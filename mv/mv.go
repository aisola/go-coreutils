//
// mv.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "bufio"
import "flag"
import "fmt"
import "os"
import "path/filepath"

const (
	help_text string = `
    Usage: mv [OPTION]... [PATH]... [PATH]
       or: mv [PATH] [PATH]
       or: mv [OPTION]
    
    move or rename files or directories
        
        --help        display this help and exit
        --version     output version information and exit

        -f, --force   remove existing destination files and never prompt the user
    ` // -v, --verbose print the name of each file before moving it
	version_text = `
    mv (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var f1 = flag.Bool("f", false, "remove existing destination files and never prompt the user")
var f2 = flag.Bool("force", false, "remove existing destination files and never prompt the user")

func input(prompt string) string {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	userinput, _ := reader.ReadString([]byte("\n")[0])

	return userinput
}

func fileExists(filep string) os.FileInfo {
	fp, err := os.Stat(filep)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return fp
}

func mover(oldp, newp string) {

	if fp := fileExists(newp); fp != nil && !*f1 && !*f2 {

		if fp.IsDir() {

			base := filepath.Base(oldp)

			if fp2 := fileExists(newp + "/" + base); fp2 != nil && !*f1 && !*f2 {

				ans := input("File '" + newp + "/" + base + "' exists. Overwrite? (y/N): ")
				if ans == "y" {
					os.Rename(oldp, newp+"/"+base)
				} else {
					os.Exit(1)
				}

			} else if fp2 != nil && (*f1 || *f2) {

				os.Rename(oldp, newp+"/"+base)

			} else if fp2 == nil {

				os.Rename(oldp, newp+"/"+base)

			}

		} else {

			ans := input("File '" + newp + "' exists. Overwrite? (y/N): ")
			if ans == "y" {
				os.Rename(oldp, newp)
			} else {
				os.Exit(1)
			}
		}

	} else if fp != nil && (*f1 || *f2) {

		os.Rename(oldp, newp)

	} else if fp == nil {

		os.Rename(oldp, newp)

	}
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	// v1 := flag.Bool("v", false, "print the name of each file before moving it")
	// v2 := flag.Bool("verbose", false, "print the name of each file before moving it")
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

	if len(files) == 2 {
		mover(files[0], files[1])
		os.Exit(0)
	} else if len(files) < 1 {
		fmt.Println("mv: destination required")
		os.Exit(1)
	} else {

		to_file, files := files[len(files)-1], files[:len(files)-1]

		if fp := fileExists(to_file); fp == nil || !fp.IsDir() {
			fmt.Println("mv: when moving multiple files, last argument must be a directory")
			os.Exit(1)
		} else {

			fmt.Println(files)
			for i := 0; i < len(files); i++ {
				mover(files[i], to_file)
			}
			os.Exit(0)

		}

	}
}
