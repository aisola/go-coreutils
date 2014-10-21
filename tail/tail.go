//
// tail.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//

package main

import "bytes"
import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "os"
import "strings"

const (
	help_text string = `
    Usage: tail [OPTION]... [FILE]...
       
    Print the last 10 lines of each FILE to standard output. With more than one FILE, precede each with a header giving the file name. With no FILE, or when FILE is -, read standard input.
    
    
    Mandatory arguments to long options are mandatory for short options too.

       -help        display this help and exit
       -version     output version information and exit

       -c, --bytes=K
              output the last K bytes; or use -n +K to output starting with the Kth byte.

       -n, -lines=K
              output the last K lines; or use -n +K to output starting with 
the Kth

       -q, -quiet, -silent
              never output headers giving file names
`
	version_text = `
    tail (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for deprintTails see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
	bytes_text  = "output the last K bytes"
	lines_text  = "output the last K lines"
	silent_text = "never output headers giving file names"
)

var (
	help        = flag.Bool("help", false, help_text)
	version     = flag.Bool("version", false, version_text)
	lines       = flag.Int("n", 10, lines_text)
	linesLong   = flag.Int("lines", 10, lines_text)
	bytesF      = flag.Int("c", 0, bytes_text)
	bytesFLong  = flag.Int("bytes", 0, bytes_text)
	silent      = flag.Bool("q", false, silent_text)
	silentLong  = flag.Bool("quiet", false, silent_text)
	silentLong2 = flag.Bool("silent", false, silent_text)
)

// bufferFile returns a byte slice of the file contents.
func bufferFile(s string) []byte {
	buffer, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Println(err, "wc: "+s+": No such file or directory")
		os.Exit(0)
	}
	return buffer
}

// silentCheck prints the file name if silent mode is enabled.
func silentCheck(filename string) {
	if !*silent {
		fmt.Printf("==> %s <==\n", filename)
	}
}

// multiFileLineProcessor prints that last K lines of every file.
func multiFileLineProcessor() {
	for index, currentFile := range flag.Args() {
		silentCheck(currentFile)
		printTailingLines(string(bufferFile(currentFile)))
		if index+1 != flag.NArg() && !*silent {
			fmt.Println()
		}
	}
}

// multiFileByteProcessor prints the last K bytes of every file.
func multiFileByteProcessor() {
	for index, currentFile := range flag.Args() {
		silentCheck(currentFile)
		printTailingBytes(bufferFile(currentFile))
		if index+1 != flag.NArg() && !*silent {
			fmt.Println()
		}
	}
}

/* splitAndCountLines splits the buffered string into a newline-deliminted
 * string slice and returns the slice along with the line count. */
func splitAndCountLines(buffer string) ([]string, int) {
	return strings.Split(buffer, "\n"), strings.Count(buffer, "\n")
}

// printTailingLines prints the last N lines from the input buffer.
func printTailingLines(buffer string) {
	lineSlice, totalLines := splitAndCountLines(buffer)
	lineCount := totalLines
	if *lines < lineCount {
		lineCount = *lines
	}
	for index := totalLines - lineCount; index < totalLines; index++ {
		fmt.Println(lineSlice[index])
	}
}

// printTailingBytes prints the last N bytes from the input buffer.
func printTailingBytes(buffer []byte) {
	totalBytes := len(buffer)
	byteCount := totalBytes
	if *bytesF < byteCount {
		byteCount = *bytesF
	}
	for index := totalBytes - byteCount; index < totalBytes; index++ {
		fmt.Print(string(buffer[index]))
	}
}

// getStdin will get input from stdin if there are no file arguments for tail.
func getStdin() {
	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, os.Stdin)
	if *bytesF == 0 {
		printTailingLines(buffer.String())
	} else {
		printTailingBytes(buffer.Bytes())
	}
}

// oneFile will use the first file argument as an argument for tail.
func oneFile() {
	if *bytesF == 0 {
		printTailingLines(string(bufferFile(flag.Arg(0))))
	} else {
		printTailingBytes(bufferFile(flag.Arg(0)))
	}
}

// multipleFiles will launch the proper function for looping through all files.
func multipleFiles() {
	if *bytesF == 0 {
		multiFileLineProcessor()
	} else {
		multiFileByteProcessor()
	}
}

func main() {
	switch {
	case flag.NArg() == 0 || flag.Arg(0) == "-":
		getStdin()
	case flag.NArg() == 1:
		oneFile()
	default:
		multipleFiles()
	}
}

func init() {
	flag.Parse()
	if *linesLong != 10 {
		*lines = *linesLong
	}
	if *silentLong || *silentLong2 {
		*silent = true
	}
	if *bytesFLong != 0 {
		*bytesF = *bytesFLong
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
