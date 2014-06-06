//
// stat.go (go-coreutils) 0.1
// Copyright (C) 2014, Abram C. Isola.
//
package main

import "flag"
import "fmt"
import "os"

const (
	help_text string = `
    Usage: stat [OPTION]...
    
    display file or file system status
    
    THIS IS PROGRAM IS IN PROGRESS

          --help     display this help and exit
          --version  output version information and exit
    `
	version_text = `
    stat (go-coreutils) 0.1

    Copyright (C) 2014 Abram C. Isola.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

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
    
    file := flag.Args()[0]
    
    fi, _ := os.Stat(file)
    
    // var ftype string
    // if fi.IsDir() {
    //     ftype = "directory"
    // } else {
    //     ftype = "regular file"
    // }
    
    // NEED: directory/file, Device, Blocks, IO Block, Inode, Links, Access (hex), Uid, Gid, Access Date
    fmt.Printf("  File: '%s'\n", fi.Name())
    fmt.Printf("  Size: %-12d Blocks: %-8d IO Block: %d %s\n", fi.Size(),8,4096,ftype)
    fmt.Printf("Device: %-12s Inode : %-8d Links: %d\n", "800h/1200d",1,1)
    fmt.Printf("Access: (%s) Uid: %s Gid: %s\n", fi.Mode(), "(1000/acisola)", "(1000/acisola)")
    fmt.Printf("Access: %s\n", fi.ModTime())
    fmt.Printf("Modify: %s\n", fi.ModTime())
    fmt.Printf("Change: %s\n", fi.ModTime())
    
}
