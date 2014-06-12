//
// logname.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "flag"
import "fmt"
import "os"
import "os/user"

const (
	help_text = `
    Usage: logname
    
    print the name of the current user

        --help        display this help and exit
        --version     output version information and exit
    `
	version_text = `
    logname (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

func GetCurrentUser(username *string) error {

	current_user, err := user.Current()

	*username = current_user.Username

	if err != nil {
		return err
	}

	return nil
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if *help {
		fmt.Println(help)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		fmt.Println(help)
		os.Exit(0)
	}

	if flag.NArg() == 0 && flag.NFlag() == 0 {

		var username string

		err := GetCurrentUser(&username)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Println(username)

	}

}
