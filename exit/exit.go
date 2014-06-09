//
// exit.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
// 
// Written By: Abram C. Isola
//
package main

import "os"
import "log"
import "fmt"
import "flag"

const (
	help_text string = `
    Usage: exit [OPTION]
    
    exit from a program, shell or log out of a unix network
        
        --help        display this help and exit
        --version     output version information and exit
    `
	version_text = `
    exit (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

// Get PID of Parent
var process = os.Getppid()

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

	pproc, err := os.FindProcess(process)

	if err != nil {
		log.Fatalln(err)
	} else {
		pproc.Kill()
	}
}
