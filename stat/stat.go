//
// stat.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "bytes"
import "flag"
import "fmt"
import "io"
import "os"
import "strings"
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
func getFileStat(index int) os.FileInfo {
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

// Opens the passwd file and returns a buffer of it's contents.
func bufferUsers() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)

	// Check to see if the file exists
	_, err := os.Stat("/etc/passwd")
	if err != nil {
		fmt.Println("Error: group file does not exist.")
		os.Exit(0)
	}

	// Cache the contents of /etc/group into a buffer
	cached, _ := os.Open("/etc/passwd")
	io.Copy(buffer, cached)
	return buffer
}

// Opens the group file and returns a buffer of it's contents.
func bufferGroups() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)

	// Check to see if the file exists
	_, err := os.Stat("/etc/group")
	if err != nil {
		fmt.Println("Error: group file does not exist.")
		os.Exit(0)
	}

	// Cache the contents of /etc/group into a buffer
	cached, _ := os.Open("/etc/group")
	io.Copy(buffer, cached)
	return buffer
}

// Converts a bytes buffer into a newline-separated string array.
func bufferToStringArray(buffer *bytes.Buffer) []string {
	return strings.Split(buffer.String(), "\n")
}

// Returns a colon separated string array for use in parsing /etc/group and /etc/user
func parseLine(line string) []string {
	return strings.Split(line, ":")
}

// Returns the username associated to a user ID
func lookupUserID(uid string) string {
	groupStringArray := bufferToStringArray(bufferUsers())
	for _, line := range groupStringArray {
		values := parseLine(line)
		if len(values) > 2 {
			if values[2] == uid {
				return values[0]
			}
		}

	}
	return uid
}

// Returns the groupname associated to a group ID
func lookupGroupID(gid string) string {
	groupStringArray := bufferToStringArray(bufferGroups())
	for _, line := range groupStringArray {
		values := parseLine(line)
		if len(values) > 2 {
			if values[2] == gid {
				return values[0]
			}
		}

	}
	return gid
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

// Resolve the symbolic link
func readLink(index int) string {
	sympath, err := os.Readlink(flag.Arg(index))
	if err == nil {
		return sympath
	} else {
		return "broken link"
	}
}

// If the file is a symbolic link, check if dereference mode is enabled.
// If dereference mode is enabled, only the path of the symbolic link is printed.
// If it is not enabled, the symlink and it's path will be printed side by side.
func dereferenceCheck(file os.FileInfo, index int) {
	if *dereference {
		fmt.Printf("  File: '%s'\n", readLink(index))
	} else {
		fmt.Printf("  File: '%s' -> '%s'\n", file.Name(), readLink(index))
	}
}

// Checks whether the file is a symbolic link and prints the file name line.
func printFileName(file os.FileInfo, index int) {
	if getType(file) == "symbolic link" {
		dereferenceCheck(file, index)
	} else {
		fmt.Printf("  File: '%s'\n", file.Name())
	}
}

// The default printing mode
func defaultMode(fi os.FileInfo, sys *syscall.Stat_t, userName, groupName string, index int) {
	printFileName(fi, index)
	fmt.Printf("  Size: %-12d Blocks: %-8d IO Block: %d %s\n", fi.Size(), sys.Blocks, sys.Blksize, getType(fi))

	// device, inode, links, permissions, uid, gid
	fmt.Printf("Device: %-12s Inode : %-8d Links: %d\n", fmt.Sprintf("%Xh/%dd", sys.Dev, sys.Dev), sys.Ino, sys.Nlink)
	fmt.Printf("Access: %s Uid: %s Gid: %s\n", fmt.Sprintf("(%#o/%s)", fi.Mode().Perm(), fi.Mode()), fmt.Sprintf("( %d/ %s)", sys.Uid, userName), fmt.Sprintf("( %d/ %s)", sys.Gid, groupName))

	// print out times
	fmt.Printf("Access: %s\n", timespecToTime(sys.Atim))
	fmt.Printf("Modify: %s\n", timespecToTime(sys.Mtim))
	fmt.Printf("Change: %s\n", timespecToTime(sys.Ctim))
}

// Loops through each argument given.
func argumentLoop() {
	for index := 0; index < flag.NArg(); index++ {
		fi := getFileStat(index)                         // Get file stats
		sys := getAdditionalFileStat(fi)                 // Get lower level file statistics.
		usr := lookupUserID(fmt.Sprintf("%d", sys.Uid))  // Get user name
		grp := lookupGroupID(fmt.Sprintf("%d", sys.Gid)) // Get group name
		defaultMode(fi, sys, usr, grp, index)                 // Send file information for printing.
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
