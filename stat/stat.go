//
// stat.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "flag"
import "fmt"
import "os"
import "os/user"
import "syscall"
import "time"

const (
	help_text string = `
    Usage: stat [FILE]...
       or: stat [OPTION]
    
    display file or file system status
    
    THIS IS PROGRAM IS IN PROGRESS

        --help     display this help and exit
        --version  output version information and exit
    `
	version_text = `
    stat (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
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
		fi, err := os.Stat(files[i])
		if err != nil {
			fmt.Printf("stat: fatal: could not open '%s': %s", files[i], err)
		}

		// get type
		var ftype string
		if fi.IsDir() {
			ftype = "directory"
		} else {
			ftype = "regular file"
		}

		// get all file information
		sys := fi.Sys().(*syscall.Stat_t)

		// get user information
		usr, err := user.LookupId(fmt.Sprintf("%d", sys.Uid))
		if err != nil {
			fmt.Println(err)
		}

		// TODO: Gid
		fmt.Printf("  File: '%s'\n", fi.Name())
		fmt.Printf("  Size: %-12d Blocks: %-8d IO Block: %d %s\n", fi.Size(), sys.Blocks, sys.Blksize, ftype)

		// device, inode, links, permissions, uid, gid
		fmt.Printf("Device: %-12s Inode : %-8d Links: %d\n", fmt.Sprintf("%Xh/%dd", sys.Dev, sys.Dev), sys.Ino, sys.Nlink)
		fmt.Printf("Access: %s Uid: %s Gid: %d\n", fmt.Sprintf("(%#o/%s)", fi.Mode().Perm(), fi.Mode()), fmt.Sprintf("( %d/ %s)", sys.Uid, usr.Username), sys.Gid)

		// print out times
		fmt.Printf("Access: %s\n", timespecToTime(sys.Atim))
		fmt.Printf("Modify: %s\n", timespecToTime(sys.Mtim))
		fmt.Printf("Change: %s\n", timespecToTime(sys.Ctim))
	}
}
