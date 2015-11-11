//
// rm.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "errors"
import "flag"
import "fmt"
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
func RemoveAll(path string) (err error) {
	// Is this a directory we need to recurse into?
	dir, serr := os.Lstat(path)
	if serr != nil {
		if serr, ok := serr.(*os.PathError); ok && (os.IsNotExist(serr.Err) || serr.Err == syscall.ENOTDIR) {
			return nil
		}
		return serr
	}

	// remove complete dir
	if dir.IsDir() {
		return removeDir(path)
	}

	// remove file
	if *interactivei && !yesOrNo("rm: do you want to remove '%s'? (y/N) " + path) {
		return nil
	}

	return os.Remove(path)
}

func removeDir(path string) (err error) {
	if !*recursiveR && !*recursiver && !*recursive {
		return errors.New(fmt.Sprintf("rm: '%s' is a directory", path))
	}

	if *interactivei && !yesOrNo("rm: descend into directory '%s'? (y/N) " + path) {
		return nil
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}

	// Remove contents.
	names, err := fd.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		if *interactivei && !yesOrNo("rm: remove '%s'? (y/N)" + path) {
			continue
		}

		err := RemoveAll(path + string(os.PathSeparator) + name)
		if err != nil {
			return err
		}
	}

	// Close directory, because windows won't remove opened directory.
	// and by the way it's always a good idea to close stuff you're done with
	fd.Close()

	// Remove directory.
	return os.Remove(path)
}

// return true on "y" and "yes"; otherwise false
func yesOrNo(question string) bool {
	var answer string

	fmt.Println(question)
	_, err := fmt.Scanln(&answer)
	if err != nil {
		goto out
	}

	if answer == "y" || answer == "yes" {
		return true
	}

out:
	return false
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()
	exitCode := 0

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
		err := RemoveAll(files[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error() + "\n")
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
