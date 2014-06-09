//
// whoami.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
// 
// Written By: Abram C. Isola
//
package main

import "os"
import "fmt"
import "log"
import "os/user"
import "flag"

const (
	help_text string = `
    Usage: whoami [OPTION]
    
    Print the user name associated with the current effective user ID.

        --help     display this help and exit
        --version  output version information and exit
    `
	version_text = `
    whoami (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	current_user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(current_user.Username)
	os.Exit(0)
}
