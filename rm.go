package main

import "flag"
import "fmt"
import "os"

const (
	help_text string = `
    Usage: rm [OPTION]...
    
    Remove files (delete/unlink)

          --help     display this help and exit
          --version  output version information and exit
          -f         ignore if files do not exist, never prompt
          -r, -R, --recursive
              remove directories and their contents recursively
    `
	version_text = `
    rm (go-coreutils) 0.1

    Copyright (C) 2014 Abram C. Isola.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	force := flag.Bool("f", false, "Ignore if files do not exist. Never prompt.")
	r1 := flag.Bool("R", false, "Remove directories and their contents recursively.")
	r2 := flag.Bool("r", false, "Remove directories and their contents recursively.")
	r3 := flag.Bool("recursive", false, "Remove directories and their contents recursively.")
	// i1 := flag.Bool("I", false, "Remove directories and their contents recursively.")
	// i2 := flag.Bool("i", false, "Remove directories and their contents recursively.")
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

		// file exists?
		if fp, err := os.Stat(files[i]); err == nil {

			if fp.IsDir() {

				if *r1 || *r2 || *r3 {

					err = os.RemoveAll(files[i])
					if err != nil {
						fmt.Printf("rm: Cannot remove '%s': %s\n", files[i], err)
					}

				} else {

					fmt.Printf("rm: '%s'is a directory\n", files[i])

				}

			} else {

				os.Remove(files[i])

			}

		} else {

			if !*force && os.IsNotExist(err) {
				fmt.Printf("rm: Cannot remove '%s': No such file or directory\n", files[i])
			} else if !*force && !os.IsNotExist(err) {
				fmt.Printf("rm: Cannot remove '%s': %s\n", files[i], err)
			}

		}
	}
}
