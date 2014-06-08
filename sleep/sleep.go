//
// sleep.go (go-coreutils) 0.1
// Copyright (C) 2014, Abram C. Isola.
//
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	help_text = `
    Usage: sleep NUMBER[SUFFIX]...
    or:  sleep OPTION
    
    Pause for NUMBER seconds.  SUFFIX may be 's' for seconds (the default),
    'm' for minutes, 'h' for hours or 'd' for days.  Unlike most implementations
    that require NUMBER be an integer, here NUMBER may be an arbitrary floating
    point number.  Given two or more arguments, pause for the amount of time
    specified by the sum of their values.

          --help     display this help and exit
          --version  output version information and exit
    `
	version_text = `
    sleep (go-coreutils) 0.1

    Copyright (C) 2014 Abram C. Isola.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

func usage() {
	fmt.Printf("sleep: missing operand\nTry 'sleep --help' for more information.\n")
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if flag.NArg() == 0 && flag.NFlag() == 0 {
		usage()
		os.Exit(1)
	}

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	var total time.Duration

	// coreutil's sleep says: "Given two or more arguments, pause for the amount
	// of time specified by the sum of their value"
	for i := 0; i < flag.NArg(); i++ {
		d, err := time.ParseDuration(flag.Arg(i))
		if err != nil {
			fmt.Printf("sleep: invalid time interval '%s'\n", flag.Arg(i))
			os.Exit(1)
		}

		total = total + d
	}

	// sleep for a total time of passed times
	time.Sleep(total)
}
