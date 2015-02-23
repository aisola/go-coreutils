package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

const (
	help_text string = `
    Usage: groups [OPTION]... [USERNAME]...

    Print the groups a user is in.

        --help     display this help and exit
        --version  output version information and exit
    `
	version_text = `
    groups (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute
    it under certain conditions in LICENSE.
`
)

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	var (
		u   *user.User
		err error
	)

	if len(os.Args) > 1 {
		username := os.Args[1]
		u, err = user.Lookup(username)
		if err != nil {
			fmt.Println("groups: " + username + ": no such user")
			os.Exit(1)
		}
	} else {
		u, err = user.Current()
	}

	if err != nil {
		log.Fatalln(err)
	}

	groups := groups(u)

	if len(os.Args) > 1 {
		fmt.Print(u.Username + " : ")
	}

	fmt.Println(strings.Join(groups, " "))
	os.Exit(0)
}

// TODO: Reading /etc/group because user.LookupGroup is not yet implemented: https://github.com/golang/go/issues/2617
func groups(u *user.User) []string {
	file, err := os.Open("/etc/group")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var groups = []string{u.Username}
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ":")
		groupname := fields[0]
		uservalue := fields[len(fields)-1]
		userlist := strings.Split(uservalue, ",")
		for _, username := range userlist {
			if username == u.Username {
				groups = append(groups, groupname)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return groups
}
