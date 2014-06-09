//
// pwd.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
// 
// Written By: Abram C. Isola
//
package main

import "fmt"
import "log"
import "os"
import "flag"

const version_text = `
    pwd (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`

func main() {
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(pwd)
	}

}
