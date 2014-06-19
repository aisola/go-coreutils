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

// If zeroLong is enabled, set zero to enabled.
func processFlags() {
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}
	if *zeroLong {
		*zero = true
	}
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
}

// Return the dirname
func getDirName() string {
	return filepath.Dir(filepath.Clean(file))
}

/* If the number of arguments given is zero, print the help text.
 * Otherwise check if the zero flag is set and print the dirname
 * of each file. */
func argumentCheck() {
	if flag.NArg() < 1 {
		fmt.Println(help_text)
		os.Exit(0)
	} else {
		for _, file := range flag.Args() {
			if *zero {
				fmt.Print(getDirName())
			} else {
				fmt.Println(getDirName())
			}
		}
	}
}

func main() {
	flag.Parse()
	processFlags()
	argumentCheck()
}
