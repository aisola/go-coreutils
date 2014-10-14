//
// rmdir.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//
package main

import "flag"
import "fmt"
import "os"

const (
	help_text = `
    Usage: rmdir [OPTION] DIRECTORY...

    Removes directories if they are empty.

        -v, -verbose
              output a diagnostic for every directory processed
      
        -help display this help and exit
        
        -version output version information and exit
`
	version_text = `
    rmdir (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	verbose     = flag.Bool("v", false, "output a diagnostic for every directory processed.")
	verboseLong = flag.Bool("verbose", false, "see verbose")
	help        = flag.Bool("help", false, "display this help and exit")
	version     = flag.Bool("version", false, "output version information and exit")
)

// printAndExit prints a message and exits the program.
func printAndExit(message string) {
	fmt.Println(message)
	os.Exit(0)
}

// argumentIsDir returns true if the argument is a directory.
func argumentIsDir(dir *string) bool {
	file, err := os.Stat(*dir)
	if err != nil {
		prefix := "rmdir: Failed to remove"
		fmt.Printf("%s '%s': no such file or directory\n", prefix, *dir)
		return false
	} else if !file.IsDir() {
		prefix := "rmdir: Failed to remove"
		fmt.Printf("%s '%s': not a directory\n", prefix, *dir)
		return false
	} else {
		return true
	}
}

// removeDirectory attempts to remove the 'dir' directory.
func removeDirectory(dir *string) {
	err := os.Remove(*dir)
	if err != nil {
		fmt.Printf("rmdir: failed to remove '%s': %s\n", *dir, err)
	}
}

func main() {
	for index := 0; index < flag.NArg(); index++ {
		arg := flag.Arg(index)
		if *verbose || *verboseLong {
			fmt.Printf("rmdir: removing directory, '%s'\n", arg)
		}
		if argumentIsDir(&arg) {
			removeDirectory(&arg)
		}
	}
}

func init() {
	flag.Parse()
	if *help {
		printAndExit(help_text)
	}
	if *version {
		printAndExit(version_text)
	}
	if flag.NArg() == 0 {
		printAndExit("Try 'rmdir --help' for more information.")
	}
}
