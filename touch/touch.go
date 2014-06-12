//
// touch.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "flag"
import "fmt"
import "os"
import "time"

const (
	help_text string = `
    Usage: touch [OPTION]...
    
    set/modify file timestamps

        --help      display this help and exit
        --version   output version information and exit

        -c          do not create if file does not exist
        -t=time     set time to time
    `
	version_text = `
    touch (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

func main() {
	create := flag.Bool("c", false, "do not create if file does not exist")
	// newTime := flag.Int("t", 0, "set to time provided")
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

	files := flag.Args()

	for i := 0; i < len(files); i++ {
		now := time.Now()

		err := os.Chtimes(files[i], now, now)
		if err != nil && *create {
			fmt.Printf("touch: cannot touch '%s'\n", files[i])
			os.Exit(1)
		}

		f, err := os.OpenFile(files[i], os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("touch: cannot touch `%s': %s\n", files[i], err)
			os.Exit(1)
		}
		f.Close()
	}
}
