//
// yes.go (go-coreutils) 0.1
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
	help_text string = `
    Usage: yes STRING
       or: yes OPTION
    
    output a string repeatedly until killed

        --help        display this help and exit
        --version     output version information and exit
    `
	version_text = `
    yes (go-coreutils) 0.1

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

	var opts = flag.Args()
	if len(opts) == 0 {
		opts = []string{"y"}
	}

	for {
		fmt.Println(opts[0])
	}
}
