//
// stat.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola, Michael Murphy
//
package main

import "flag"
import "fmt"
import "os"
import "os/user"
import "syscall"
import "time"

const (
	isExecutable = 0111              // isExcutable
	isSymlink    = os.ModeSymlink    // isSymlink
	isDevice     = os.ModeDevice     // isDevice
	isCharDevice = os.ModeCharDevice // isCharDevice

	help_text string = `
    Usage: stat [FILE]...
       or: stat [OPTION]
    
    display file or file system status
    
    THIS IS PROGRAM IS IN PROGRESS

        -L, -dereference
              follow links
          
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

var (
	dereference     = flag.Bool("L", false, "")
	dereferenceLong = flag.Bool("dereference", false, "")
)

// Process the initial flags.
func processFlags() {
	if *dereferenceLong {
		*dereference = true
	}
}

// Obtain file statistics
func getFileStat() os.FileInfo {
	fi, err := os.Lstat(flag.Arg(0))
	if err != nil {
		fmt.Printf("stat: fatal: could not open '%s': %s\n", flag.Arg(0), err)
		os.Exit(0)
	}
	return fi
}

// Obtains all file statistics information
func getAdditionalFileStat(fi os.FileInfo) *syscall.Stat_t {
	return fi.Sys().(*syscall.Stat_t)
}

// Get user information
func getUserInfo(sys *syscall.Stat_t) *user.User {
	usr, err := user.LookupId(fmt.Sprintf("%d", sys.Uid))
	if err != nil {
		fmt.Println(err)
	}
	return usr
}

// Obtain the file mode type
func getType(file os.FileInfo) string {
	switch {
	case file.IsDir():
		return "directory"
	case file.Mode()&isSymlink != 0:
		return "symbolic link"
	case file.Mode()&isDevice != 0:
		return "device file"
	case file.Mode()&isCharDevice != 0:
		return "character special file"
	case file.Mode()&isExecutable != 0:
		return "executable file"
	}
	return "regular file"
}

// Convert timespec to time
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

// If the file is a symbolic link, check if dereference mode is enabled.
// If dereference mode is enabled, only the path of the symbolic link is printed.
// If it is not enabled, the symlink and it's path will be printed side by side.
func dereferenceCheck(file os.FileInfo) {
	symPath, _ := os.Readlink(flag.Arg(0))
	if *dereference {
		fmt.Printf("  File: '%s'\n", symPath)
	} else {
		fmt.Printf("  File: '%s' -> '%s'\n", file.Name(), symPath)
	}
}

// Checks whether the file is a symbolic link and prints the file name line.
func printFileName(file os.FileInfo, index int) {
	if getType(file) == "symbolic link" {
		dereferenceCheck(file)
	} else {
		fmt.Printf("  File: '%s'\n", file.Name())
	}
}

// The default printing mode
func defaultMode(fi os.FileInfo, sys *syscall.Stat_t, usr *user.User, index int) {
	// TODO: Gid
	printFileName(fi, index)
	fmt.Printf("  Size: %-12d Blocks: %-8d IO Block: %d %s\n", fi.Size(), sys.Blocks, sys.Blksize, getType(fi))

	// device, inode, links, permissions, uid, gid
	fmt.Printf("Device: %-12s Inode : %-8d Links: %d\n", fmt.Sprintf("%Xh/%dd", sys.Dev, sys.Dev), sys.Ino, sys.Nlink)
	fmt.Printf("Access: %s Uid: %s Gid: %d\n", fmt.Sprintf("(%#o/%s)", fi.Mode().Perm(), fi.Mode()), fmt.Sprintf("( %d/ %s)", sys.Uid, usr.Username), sys.Gid)

	// print out times
	fmt.Printf("Access: %s\n", timespecToTime(sys.Atim))
	fmt.Printf("Modify: %s\n", timespecToTime(sys.Mtim))
	fmt.Printf("Change: %s\n", timespecToTime(sys.Ctim))
}

// Loops through each argument given.
func argumentLoop() {
	for index := 0; index < flag.NArg(); index++ {
		fi := getFileStat()              // Get file stats
		sys := getAdditionalFileStat(fi) // Get lower level file statistics.
		usr := getUserInfo(sys)          // Get user information
		defaultMode(fi, sys, usr, index) // Send file information for printing.
	}
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()
	processFlags()

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	argumentLoop()
}
