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
	bytes_text      = "Print the byte counts"
	characters_text = "Print the character counts"
	lines_text      = "Print the newline counts"
	words_text      = "Print the word counts"
	linelength_text = "Print the length of the longest line"
)

var (
	countBytes          = flag.Bool("c", false, bytes_text)
	countBytesLong      = flag.Bool("bytes", false, bytes_text)
	countCharacters     = flag.Bool("m", false, characters_text)
	countCharactersLong = flag.Bool("chars", false, characters_text)
	countLines          = flag.Bool("l", false, lines_text)
	countLinesLong      = flag.Bool("lines", false, lines_text)
	countWords          = flag.Bool("w", false, words_text)
	countWordsLong      = flag.Bool("words", false, words_text)
	maxLineLength       = flag.Bool("L", false, linelength_text)
	maxLineLengthLong   = flag.Bool("max-line-length", false, linelength_text)
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

/* errorChecker performs error handling via printing an error message and then
 * cleanly exiting the program. */

func errorChecker(input error, message string) {
	if input != nil {
		fmt.Println(message)
		os.Exit(0)
	}
}

// The openFile function will open the file and check for errors.

func openFile(s string) (io.ReadWriteCloser, error) {
	fi, err := os.Stat(s)
	errorChecker(err, "wc: "+s+": No such file or directory")

	if fi.Mode()&os.ModeSocket != 0 {
		return net.Dial("unix", s)
	}
	return os.Open(s)
}

/* To find the maximum string length, we must loop through a newline-delimited
 * string while counting the length of each line with len(). If the line is
 * longer than the longest recorded line before it, maxStringLength will be
 * updated to reflect it. */

func countMaxStringLength(input []string) int {
	var maxStringLength int
	for _, line := range input {
		if len(line) > maxStringLength {
			maxStringLength = len(line)
		}
	}
	return maxStringLength
}

/* The bufferProcessor function will take the buffered input and process it
 * uniquely based on which flag was given to the program. */

func bufferProcessor(fileName *string, buffer *bytes.Buffer) {
	switch {
	case *countLines: // Print the number of lines
		fmt.Println(strings.Count(buffer.String(), "\n"), *fileName)
	case *countWords: // Print the number of words
		fmt.Println(len(strings.Fields(buffer.String())), *fileName)
	case *countCharacters: // Print the number of characters (same as bytes?)
		fmt.Println(buffer.Len(), *fileName)
	case *countBytes: // Print the number of bytes
		fmt.Println(buffer.Len(), *fileName)
	case *maxLineLength: // Print the length of the longest line
		fmt.Println(countMaxStringLength(strings.Split(buffer.String(), "\n")), *fileName)
	default:
		fmt.Println(strings.Count(buffer.String(), "\n"), len(strings.Fields(buffer.String())), buffer.Len(), *fileName)
	}
}

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()
	processFlags() // Process long-form flags to their short-form variants.
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
		errorChecker(err, "wc: "+args[0]+": No such file or directory")
		io.Copy(buffer, file)
		file.Close()
	}

	bufferProcessor(&args[0], buffer) // Send the buffer for processing.
}
