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
	       
	basename -s .go /src/basename.go .go
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
	help         = flag.Bool("help", false, help_text)
	version      = flag.Bool("version", false, version_text)
)

// If zeroLong is enabled, set zero to enabled.
func processFlags() {
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}
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

// A switch to check arguments and process them accordingly.
func argumentCheck() {
	switch {
	case flag.NArg() < 1: // If there are no arguments
		fmt.Println(help_text)
	case flag.NArg() == 1: // If there is only one  argument
		checkSuffix(getBaseName())
	case flag.NArg() == 2 && suffixExists(): // If there is an argument and a suffix
		fmt.Println(strings.TrimSuffix(getBaseName(), flag.Arg(len(flag.Args()) - 1)))
	case !*multiple: // If multiple is disabled but there is more than one argument
		fmt.Println(getBaseName())
	case *multiple: // If multiple is enabled and there is more than one argument
		multiFilePrinter()
	}
}

// Obtain the basename.
func getBaseName() string {
	return filepath.Base(flag.Arg(0))
}

// Checks if a suffix is set and prints the basename accordingly.
func checkSuffix(baseName string) {
	if *suffix != "nil" {
		fmt.Println(strings.TrimSuffix(baseName, *suffix))
	} else {
		fmt.Println(baseName)
	}
}

// Trim suffix from the basename of the file.
func trimSuffix(baseName string) string {
	return strings.TrimSuffix(baseName, *suffix)
}

// Check if the last argument is a suffix
func suffixExists() bool {
	if strings.HasPrefix(flag.Arg(len(flag.Args())-1), ".") {
		return true
	} else {
		return false
	}
}

// Used in multiFilePrinter for checking if zeroMode is enabled.
func checkZero(baseName string) {
	switch {
	case *suffix != "nil" && *zero:
		fmt.Print(strings.TrimSuffix(baseName, *suffix))
	case *suffix != "nil":
		fmt.Println(strings.TrimSuffix(baseName, *suffix))
	case *zero:
		fmt.Print(baseName)
	default:
		fmt.Println(baseName)
	}
}

// Prints all basenames
func multiFilePrinter() {
	var arguments int
	if suffixExists() {
		*suffix = flag.Arg(len(flag.Args())-1)
		arguments = len(flag.Args())-1
	} else {
		arguments = len(flag.Args())
	}

	for index := 0; index < arguments; index++ {
		checkZero(filepath.Base(flag.Arg(index)))
	}
}

func main() {
	flag.Parse()
	processFlags()
	argumentCheck()
}
