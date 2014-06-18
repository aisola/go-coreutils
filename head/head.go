//
// head.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//

package main

import "bytes"
import "flag"
import "fmt"
import "io"
import "os"
import "strings"

const (
	help_text string = `
    Usage: printHead [OPTION]... [FILE]...
       
    Print the first 10 lines of each FILE to standard output. With more than one FILE, precede
    each with a printHeader giving the file name. With no FILE, or when FILE is -, read standard input.
    
    
    Mandatory arguments to long options are mandatory for short options too.

       -help        display this help and exit
       -version     output version information and exit
       
       -c, -bytes=K
             print the first K bytes of each file
       
       -n, -lines=K
              output the first K lines

       -q, -quiet, -silent
              never output printHeaders giving file names
`
	version_text = `
    cat (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for deprintHeads see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
	bytes_text  = "print the first K bytes of each file"
	lines_text  = "output the last K lines"
	silent_text = "never output printHeaders giving file names"
)

var (
	help          = flag.Bool("help", false, help_text)
	version       = flag.Bool("version", false, version_text)
	bytesMode     = flag.Int("c", 0, bytes_text)
	bytesModeLong = flag.Int("bytes", 0, bytes_text)
	lines         = flag.Int("n", 10, lines_text)
	linesLong     = flag.Int("lines", 10, lines_text)
	silent        = flag.Bool("q", false, silent_text)
	silentLong    = flag.Bool("quiet", false, silent_text)
	silentLong2   = flag.Bool("silent", false, silent_text)
)

/* The processFlags function will check initial values of flags and act
 * accordingly. */

func processFlags() {
	if *bytesModeLong != 0 {
		*bytesMode = *bytesModeLong
	}
	if *linesLong != 10 {
		*lines = *linesLong
	}
	if *silentLong || *silentLong2 {
		*silent = true
	}
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}
}

func errorChecker(input error, message string) {
	if input != nil {
		fmt.Println(message)
		os.Exit(0)
	}
}

/* The bufferLines function will open the file and copy it's contents into a
 * bytes.Buffer, which is then sent to the printHead function for printing. */

func bufferLines(s string) {
	buffer := bytes.NewBuffer(nil)
	file, err := os.Open(s)
	errorChecker(err, "wc: "+s+": No such file or directory")
	io.Copy(buffer, file)
	file.Close()
	printHead(buffer)
}

/* The bufferBytes function will create a byte array as long as the specified
 * bytesMode parameter is set for, and then reads the file into that buffer and
 * prints the results. */

func bufferBytes(s string) {
	bytesBuffer := make([]byte, *bytesMode)
	file, err := os.Open(s)
	errorChecker(err, "wc: "+s+": No such file or directory")
	file.Read(bytesBuffer)
	file.Close()
	fmt.Println(string(bytesBuffer))
}

// If silent mode is enabled, do not print the file name.

func silentCheck(filename string) {
	if !*silent {
		fmt.Printf("==> %s <==\n", filename)
	}
}

/* If there is more than one file to process, this function will loop through
 * each file argument, check if silent is enabled, then check if
 * we are to be printing a specified number of bytes or lines, and sends the
 * current file to either bufferBytes or bufferLines for printing.
 *
 * In addition to printing lines, this function will add a newline before every
 * new file. */

func multiFileProcessor() {
	for index, currentFile := range flag.Args() {
		silentCheck(currentFile)
		if *bytesMode != 0 {
			bufferBytes(currentFile)
		} else {
			bufferLines(currentFile)
			if index+1 != flag.NArg() && !*silent {
				fmt.Println()
			}
		}
	}
}

/* The bufferToStringArray function will convert the bytes.Buffer into a string
 * array. */

func bufferToStringArray(buffer *bytes.Buffer) []string {
	bufferString := buffer.String()
	return strings.Split(bufferString, "\n")
}

/* The printHead function will take the buffered bytes input and send it to
 * bufferToStringArray. After obtaining the string array of the buffer, printHead will
 * loop through each line, ending at the Kth line. */

func printHead(buffer *bytes.Buffer) {
	stringArray := bufferToStringArray(buffer)
	for index := 0; index < *lines; index++ {
		fmt.Println(stringArray[index])
	}
}

func main() {
	flag.Parse()
	processFlags()

	/* If no file is given, or the file is -, read standard input
	 * and output to standard output. Otherwise, open the file and
	 * begin reading it. */

	switch { // If there are no files or the file is "-", copy stdin to stdout.
	case flag.NArg() < 1 || flag.Arg(0) == "-":
		io.Copy(os.Stdout, os.Stdin)
	case flag.NArg() == 1: // If there is only one file
		if *bytesMode != 0 {
			bufferBytes(flag.Arg(0))
		} else {
			bufferLines(flag.Arg(0))
		}
	default: // If there is more than one file
		multiFileProcessor()
	}
}
