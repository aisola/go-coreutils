//
// wc.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//

package main

import "bytes"
import "flag"
import "fmt"
import "io"
import "net"
import "os"
import "strings"

const (
	help_text string = `
    Usage: wc [OPTION]... [FILE]...
       
    Print newline, word, and byte counts for each FILE, and a total line if
    more than one FILE is specifileed. With no FILE, or when FILE is -,
    read standard input. A word is a non-zero-length sequence of characters
    delimited by white spaces.
    
    
    The options below may be used to select which counts are printed, always in
    the following order: newline, word, character, byte, maximum line length.

        -help        display this help and exit
        -version     output version information and exit
        
        -c, -bytes
              print the byte counts
        
        -m, -chars
              print the character counts
        
        -l, -lines
              print the newline counts
              
        -L, -max-line-length
              print the length of the longest line
              
        -w, -words
              print the word counts
`
	version_text = `
    cat (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	countBytes          = flag.Bool("c", false, "Print the byte counts")
	countBytesLong      = flag.Bool("bytes", false, "Print the byte counts")
	countCharacters     = flag.Bool("m", false, "Print the character counts")
	countCharactersLong = flag.Bool("chars", false, "Print the character counts")
	countLines          = flag.Bool("l", false, "Print the newline counts")
	countLinesLong      = flag.Bool("lines", false, "Print the newline counts")
	countWords          = flag.Bool("w", false, "Print the word counts")
	countWordsLong      = flag.Bool("words", false, "Print the word counts")
	maxLineLength       = flag.Bool("L", false, "Print the length of the longest line")
	maxLineLengthLong   = flag.Bool("max-line-length", false, "Print the length of the longest line")
)

// The processFlags function will process the long forms of flags.

func processFlags() {
	switch {
	case *countBytesLong:
		*countBytes = true
	case *countLinesLong:
		*countLines = true
	case *countCharactersLong:
		*countCharacters = true
	case *countWordsLong:
		*countWords = true
	case *maxLineLengthLong:
		*maxLineLength = true
	}
}

// The getFile function will get the file from flag arguments.

func errorChecker(input error) {
	if input != nil {
		panic(input)
	}
}

// The openFile function will open the file.

func openFile(s string) (io.ReadWriteCloser, error) {
	fi, err := os.Stat(s)
	errorChecker(err)
	
	if fi.Mode()&os.ModeSocket != 0 {
		return net.Dial("unix", s)
	}
	return os.Open(s)
}

// Count the maximum string length for the maxLineLength case.

func countMaxStringLength(input []string) int {
	var maxStringLength int
	for _, line := range input {
		if len(line) > maxStringLength {
			maxStringLength = len(line)
		}
	}
	return maxStringLength
}

// The outputPrinter function will take the buffered input and process it.

func outputPrinter(fileName *string, buffer *bytes.Buffer) {
	switch {
	case *countBytes: // Print the number of bytes
		fmt.Println(buffer.Len(), *fileName)
	case *countLines: // Print the number of lines
		fmt.Println(strings.Count(buffer.String(), "\n"))
	case *countCharacters: // Print the number of characters (same as bytes?)
		fmt.Println(buffer.Len(), *fileName)
	case *countWords: // Print the number of words
		fmt.Println(len(strings.Fields(buffer.String())), *fileName)
	case *maxLineLength: // Print the length of the longest line
		fmt.Println(countMaxStringLength(strings.Split(buffer.String(), "\n")), *fileName)
	}
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()
	args := flag.Args()
	buffer := bytes.NewBuffer(nil) // Used to buffer the input
	
	// Display help information

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	
	// Display version information
	
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	/* If no file is given, or the file is -, read standard input
	 * and output to standard output. Otherwise, open the file and
	 * begin reading it. */
	
	if len(args) < 1 || args[0] == "-" {
		io.Copy(os.Stdout, os.Stdin)
	} else {
		file, err := openFile(args[0])
		errorChecker(err)
		io.Copy(buffer, file)
		file.Close()
	}
	
	outputPrinter(&args[0], buffer) // Send the buffer for processing.
}
