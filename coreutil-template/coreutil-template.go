package main

import "flag"
import "fmt"
import "os"

const (
	help_text string = `
    Usage: coreutil-template [OPTION]...
    
    A dummy template for the source of go-coreutils.

          --help     display this help and exit
          --version  output version information and exit
    `
	version_text = `
    coreutil-template (go-coreutils) 0.1

    Copyright (C) 2014 Abram C. Isola.
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

	// other code

}
