//
// true.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Trey Tacon, Abram C. Isola
//
package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	help_text = `
    usage: true [ignored command line arguments]
       or: true OPTION
    
    Exit with a status code indicating success

        --help     display this help and exit
        --version  output version information and exit
    `
	version_text = `
    true (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	help    = flag.Bool("help", false, help_text)
	version = flag.Bool("version", false, version_text)
)

func main() {
	flag.Parse()

	if *help {
		fmt.Println(help_text)
		return
	}

	if *version {
		fmt.Println(version_text)
		return
	}

	os.Exit(0)
}
