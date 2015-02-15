//
// mkdir.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Corey Prak
//
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	help_text string = `
    Usage: mkdir OPTION(S) DIRECTORY

        --help        display this help and exit
        --version     output version information and exit
        --parents     create parent directory/directories as
                      needed, do nothing if already existing
        --verbose     print a message for each created directory
  `

	version_text = `
    mkdir (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute
    it under certain conditions in LICENSE.
  `

	parents_text = `
    Create parent directory/directories as needed for a path. If any
    already exist, do nothing.
  `

	usage_text = `usage: mkdir [-parents, -verbose] directory ...`

	verbose_text = `Print a message for each created directory.`
)

var (
	help    = flag.Bool("help", false, help_text)
	version = flag.Bool("version", false, version_text)
	parents = flag.Bool("parents", false, parents_text)
	verbose = flag.Bool("verbose", false, verbose_text)
)

func extend(slice []string, element string) []string {
	if len(slice) == cap(slice) {
		newSlice := make([]string, len(slice), 2*cap(slice))
		copy(newSlice, slice)
		slice = newSlice
	}
	n := len(slice)
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func getAllPaths(dir string) []string {
	slice := make([]string, 0, 10)
	i := 0
	_, err := os.Stat(dir)
	for err != nil {
		if dir == "." || dir == "/" {
			break
		}
		slice = extend(slice, dir)
		i += 1
		dir = filepath.Dir(dir)
		_, err = os.Stat(dir)
	}
	return slice
}

func printAllPaths(slice []string) {
	i := len(slice)
	for i > 0 {
		i -= 1
		fmt.Printf("mkdir: created directory `%s'\n", slice[i])

	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println(usage_text)
		os.Exit(0)
	}

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	for i := 0; i < flag.NArg(); i++ {

		if *parents {
			paths := getAllPaths(flag.Arg(i))
			mkdirAllError := os.MkdirAll(flag.Arg(i), os.ModePerm)

			if mkdirAllError != nil {
				fmt.Println(mkdirAllError)
			} else if *verbose {
				printAllPaths(paths)
			}
		} else {
			mkdirError := os.Mkdir(flag.Arg(i), os.ModePerm)

			if mkdirError != nil {
				fmt.Println(mkdirError)
			} else if *verbose {
				fmt.Printf("mkdir: created directory '%s'\n", flag.Arg(i))
			}
		}

	}
}
