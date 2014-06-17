//
// dirname.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Trey Tacon, Abram C. Isola, Michael Murphy
//
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	help_text = `
    usage: dirname [OPTION] NAME...
    
    Output each NAME with its last non-slash component and trailing slashes
    removed; if NAME contains no /'s, output '.' (meaning the current directory).

        -help     display this help and exit
        -version  output version information and exit
        
        -z, -zero
              separate output with NUL rather than newline
    `
	version_text = `
    dirname (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
	zero_text = "separate output with NUL rather than newline"
)

var (
	help     = flag.Bool("help", false, help_text)
	version  = flag.Bool("version", false, version_text)
	zero     = flag.Bool("z", false, zero_text)
	zeroLong = flag.Bool("zero", false, zero_text)
)

// Display help information

func helpCheck() {
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
}

// Display version information

func versionCheck() {
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}
}

// If zeroLong is enabled, set zero to enabled.

func processFlags() {
	if *zeroLong {
		*zero = true
	}
}

func main() {
	flag.Parse()
	processFlags()
	helpCheck()
	versionCheck()

	/* If the number of arguments given is zero, print the help text.
	 * Otherwise check if the zero flag is set and print the dirname
	 * of each file. */

	if flag.NArg() < 1 {
		fmt.Println(help_text)
		os.Exit(0)
	} else {
		/* NOTE(ttacon): we need to clean the directory first
		 * as filepath.Dir will not remove the last directory
		 * of the path the way dirname will if it ends in
		 * a trailing slash - e.g.:
		 *
		 * filepath.Dir("/Users/ttacon/") == "/Users/ttacon/"
		 * while '$ dirname /Users/ttacon/' == /Users */

		for _, file := range flag.Args() {
			if *zero {
				fmt.Print(filepath.Dir(filepath.Clean(file)))
			} else {
				fmt.Println(filepath.Dir(filepath.Clean(file)))
			}
		}
	}
}
