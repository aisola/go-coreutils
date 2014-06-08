//
// false.go (go-coreutils) 0.1
// Copyright (C) 2014, Abram C. Isola.
//
package main

import "flag"
import "log"
import "os"

const (
	help_text = `
	Usage: false [ignored command line arguments]
  	or:  false OPTION
    
	Exit with a status code indicating failure.

      --help     display this help and exit
      --version  output version information and exit
    `
	version_text = `
    false (go-coreutils) 0.1

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

	if flag.NFlag() > 1 {
		os.Exit(-1)
	}

	if *help {
		log.Fatal(help_text)
	}

	if *version {
		log.Fatal(version_text)
	}
	os.Exit(-1)
}
