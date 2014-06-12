//
// rm.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "flag"
import "fmt"
import "io"
import "os"
import "syscall"

const (
	help_text string = `
    Usage: rm [OPTION]...
    
    remove files (delete/unlink)

        --help     display this help and exit
        --version  output version information and exit
        -f         ignore if files do not exist, never prompt
        -i         prompt before each removal
        -r, -R, --recursive
            remove directories and their contents recursively
    `
	version_text = `
    rm (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	force        = flag.Bool("f", false, "ignore if files do not exist, never prompt")
	recursiveR   = flag.Bool("R", false, "remove directories and their contents recursively")
	recursiver   = flag.Bool("r", false, "remove directories and their contents recursively")
	recursive    = flag.Bool("recursive", false, "remove directories and their contents recursively")
	interactivei = flag.Bool("i", false, "prompt before each removal")
	// interactiveI = flag.Bool("I", false, "")
)

// MODIFIED FROM THE os.RemoveAll() implimentation
// RemoveAll removes all files/directories below
// and prompts if the option is set.
func RemoveAll(path string) error {
	var answer string
	var err error

	// Is this a directory we need to recurse into?
	dir, serr := os.Lstat(path)
	if serr != nil {
		if serr, ok := serr.(*os.PathError); ok && (os.IsNotExist(serr.Err) || serr.Err == syscall.ENOTDIR) {
			return nil
		}
		return serr
	}

	if !dir.IsDir() {
		// Not a directory;
		if *interactivei {
			fmt.Printf("rm: do you want to remove '%s'? (y/N) ", path)
			_, err = fmt.Scanln(&answer)
			if err != nil {
				return err
			}

			if answer == "y" || answer == "yes" {
				os.Remove(path)
			}
		} else {
			os.Remove(path)
		}

		return nil
	}

	// Turns out, it's a directory...
	if !*recursiveR && !*recursiver && !*recursive {
		fmt.Printf("rm: '%s' is a directory\n", path)
		return nil
	}

	if *interactivei {
		fmt.Printf("rm: descend into directory '%s'? (y/N) ", path)
		_, err = fmt.Scanln(&answer)
		if err != nil {
			return err
		}

		if answer == "y" || answer == "yes" {

			fd, err := os.Open(path)
			if err != nil {
				return err
			}

			// Remove contents & return first error.
			err = nil
			for {
				names, err1 := fd.Readdirnames(100)
				for _, name := range names {
					err1 := RemoveAll(path + string(os.PathSeparator) + name)
					if err == nil {
						err = err1
					}
				}
				if err1 == io.EOF {
					break
				}
				// If Readdirnames returned an error, use it.
				if err == nil {
					err = err1
				}
				if len(names) == 0 {
					break
				}
			}

			// Close directory, because windows won't remove opened directory.
			fd.Close()

		} else {
			return nil
		}

	} else {
		fd, err := os.Open(path)
		if err != nil {
			return err
		}

		// Remove contents & return first error.
		err = nil
		for {
			names, err1 := fd.Readdirnames(100)
			for _, name := range names {
				err1 := RemoveAll(path + string(os.PathSeparator) + name)
				if err == nil {
					err = err1
				}
			}
			if err1 == io.EOF {
				break
			}
			// If Readdirnames returned an error, use it.
			if err == nil {
				err = err1
			}
			if len(names) == 0 {
				break
			}
		}

		// Close directory, because windows won't remove opened directory.
		fd.Close()
	}

	if *interactivei {

		fmt.Printf("rm: remove '%s'? (y/N) ", path)
		_, err := fmt.Scanln(&answer)
		if err != nil {
			return err
		}

		if answer == "y" || answer == "yes" {
			// Remove directory.
			err1 := os.Remove(path)
			if err == nil {
				err = err1
			}
		}

	} else {
		// Remove directory.
		err1 := os.Remove(path)
		if err == nil {
			err = err1
		}
	}

	return err
}

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

	files := flag.Args()
	for i := 0; i < len(files); i++ {
		RemoveAll(files[i])
	}
}
