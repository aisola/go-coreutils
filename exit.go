package main

import "os"
import "log"
import "fmt"
import "flag"

const version_text = `
    exit (go-coreutils) 0.1

    Copyright (C) 2014 Abram C. Isola.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`

// Get PID of Parent
var process = os.Getppid()

func main() {
	version := flag.Bool("version", false, version_text)
	flag.Parse()

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
