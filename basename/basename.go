//
// basename.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola, Michael Murphy
//
package main

import "flag"
import "fmt"
import "os"
import "path/filepath"
import "strings"

const (
	help_text string = `
    Usage: basename [SUFFIX]
       or: basename OPTION... NAME...
    
    Print NAME with any leading directory components removed.  If specified, also remove a trailing SUFFIX.

        -help     display this help and exit
        -version  output version information and exit
        
        -a, -multiple
               support multiple arguments and treat each as a NAME
               
        -s, -suffix
               remove a trailing SUFFIX
               
        -z, -zero separate output with NUL rather than newline
        
    Examples
        basename /usr/bin/sort
               -> "sort"
	       
	basename /src/basename.go .go
	       -> "basename"
	       
	basename -s .h /src/basename.go .go
	       -> "basename"
	       
	basename -a any/str1 any/str2
	       -> "str1"
	       -> "str2"
    `
	version_text = `
    basename (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
	multiple_text = "support multiple arguments and treat each as a NAME"
	suffix_text   = "remove a trailing SUFFIX"
	zero_text     = "separate output with NUL rather than newline"
)

var (
	multiple     = flag.Bool("a", false, multiple_text)
	multipleLong = flag.Bool("multiple", false, multiple_text)
	suffix       = flag.String("s", "nil", suffix_text)
	suffixLong   = flag.String("suffix", "nil", suffix_text)
	zero         = flag.Bool("z", false, zero_text)
	zeroLong     = flag.Bool("zero", false, zero_text)
)

// Display help information

func helpCheck(help *bool) {
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
}

// Display version information

func versionCheck(version *bool) {
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
	if *multipleLong {
		*multiple = true
	}
	if *suffixLong != "nil" {
		*suffix = *suffixLong
	}
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()
	helpCheck(help)
	versionCheck(version)
	var fileBase string

	switch {
	case flag.NArg() < 1: // If there are no arguments
		fmt.Println(help_text)
	case flag.NArg() == 1: // If there is only one argument
		fileBase = filepath.Base(flag.Arg(0))
		if *suffix != "nil" {
			fmt.Println(strings.TrimSuffix(fileBase, *suffix))
		} else {
			fmt.Println(fileBase)
		}
	case flag.NArg() == 2 && // If two arguments are given, and the last one is a suffix,
		strings.HasPrefix(flag.Arg(len(flag.Args())-1), "."): // trim the first argument.
		*suffix = flag.Arg(len(flag.Args()) - 1)
		fmt.Println(strings.TrimSuffix(filepath.Base(flag.Arg(0)), *suffix))
	case !*multiple: // If multiple is disabled but there is more than one argument
		fmt.Println(filepath.Base(flag.Arg(0)))
	case *multiple: // If multiple is enabled and there is more than one argument
		var arguments int
		if strings.HasPrefix(flag.Arg(len(flag.Args())-1), ".") {
			*suffix = flag.Arg(len(flag.Args()) - 1)
			arguments = len(flag.Args()) - 1
		} else {
			arguments = len(flag.Args())
		}

		for index := 0; index < arguments; index++ {
			fileBase = filepath.Base(flag.Arg(index))
			if *suffix != "nil" {
				if *zero {
					fmt.Print(strings.TrimSuffix(fileBase, *suffix))
				} else {
					fmt.Println(strings.TrimSuffix(fileBase, *suffix))
				}
			} else {
				if *zero {
					fmt.Print(fileBase)
				} else {
					fmt.Println(fileBase)
				}
			}
		}

	}
}
