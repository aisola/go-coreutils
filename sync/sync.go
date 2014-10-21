//
// sync.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy, Abram C. Isola
//
package main

import "fmt"
import "flag"
import "syscall"

const (
	helpText = `
    Usage: sync [OPTION]
    
    Force changed blocks to disk; update the super block.
    
        -help display this help and exit
        
        -version
              output version information and exit
`
	versionText = `
    sync (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	help    = flag.Bool("help", false, "display help information")
	version = flag.Bool("version", false, "display version information")
)

func main() {
	syscall.Sync()
}

func init() {
	flag.Parse()
	if *help {
		fmt.Println(helpText)
	}
	if *version {
		fmt.Println(versionText)
	}
}
