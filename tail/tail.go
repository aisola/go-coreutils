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
import "os"
import "strings"

const (
	help_text string = `
    Usage: tail [OPTION]... [FILE]...
       
    Print the last 10 lines of each FILE to standard output. With more than one FILE, precede
    each with a header giving the file name. With no FILE, or when FILE is -, read standard input.
    
    
    Mandatory arguments to long options are mandatory for short options too.

       -help        display this help and exit
       -version     output version information and exit


       -n, -lines=K
              output the last K lines, instead of the last 10; or use -n +K to output starting with the Kth

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
	lines_text  = "output the last K lines"
	silent_text = "never output headers giving file names"
)

var (
	help        = flag.Bool("help", false, help_text)
	version     = flag.Bool("version", false, version_text)
	lines       = flag.Int("n", 10, lines_text)
	linesLong   = flag.Int("lines", 10, lines_text)
	silent      = flag.Bool("q", false, silent_text)
	silentLong  = flag.Bool("quiet", false, silent_text)
	silentLong2 = flag.Bool("silent", false, silent_text)
)

/* The processFlags function will check initial values of flags and act
 * accordingly. */

func processFlags() {
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
 * bytes.Buffer, which is then sent to the printTail function for printing. */

func bufferLines(s string) {
	buffer := bytes.NewBuffer(nil)
	file, err := os.Open(s)
	errorChecker(err, "wc: "+s+": No such file or directory")
	io.Copy(buffer, file)
	file.Close()
	printTail(buffer)
}

// If silent mode is enabled, do not print the file name.

func silentCheck(filename string) {
	if !*silent {
		fmt.Printf("==> %s <==\n", filename)
	}
}

/* If there is more than one file to process, this function will loop through
 * each file argument, check if silent is enabled, then finally send the file
 * to be buffered and printed by bufferLines.
 *
 * In addition to printing lines, this function will add a newline before every
 * new file. */

func multiFileProcessor() {
	for index, currentFile := range flag.Args() {
		silentCheck(currentFile)
		bufferLines(currentFile)
		if index+1 != flag.NArg() && !*silent {
			fmt.Println()
		}
	}
}

/* The splitAndCount function will return the number of lines and a newline delimited
 * string array. To do this, we first need to convert the buffer into a string via
 * buffer.String(), then we can find the number of strings by using strings.Count()
 * and split the string with strings.Split, both accepting a "\n" for delimiting. */

func splitAndCount(buffer *bytes.Buffer) (int, []string) {
	bufferString := buffer.String()
	return strings.Count(bufferString, "\n"), strings.Split(bufferString, "\n")
}

/* The printTail function will take the buffered bytes input and send it to splitAndCount.
 * After obtaining the number of lines and a string array of the buffer, printTail will
 * loop through each line, starting at the Kth file from the end. */

func printTail(buffer *bytes.Buffer) {
	numOfLines, stringArray := splitAndCount(buffer)
	for index := numOfLines - *lines; index < numOfLines; index++ {
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
		bufferLines(flag.Arg(0))
	default: // If there is more than one file
		multiFileProcessor()
	}
}
