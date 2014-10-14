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
import "os"
import "strings"

var (
	countBytes              = flag.Bool("c", false, "Print the byte counts")
	countBytesL             = flag.Bool("bytes", false, "Print the byte counts")
	countCharacters         = flag.Bool("m", false, "Print the character counts")
	countCharactersL        = flag.Bool("chars", false, "Print the character counts")
	countLines              = flag.Bool("l", false, "Print the newline counts")
	countLinesL             = flag.Bool("lines", false, "Print the newline counts")
	countSLOC               = flag.Bool("sloc", false, "Print the source lines of code")
	occurrenceRef           = flag.String("o", "", "Print the occurrences of a particular word or phrase")
	countWords              = flag.Bool("w", false, "Print the word counts")
	countWordsL             = flag.Bool("words", false, "Print the word counts")
	maxLineLength           = flag.Bool("L", false, "Print the length of the longest line")
	maxLineLengthL          = flag.Bool("max-line-length", false, "Print the length of the longest line")
	help_text        string = `
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
              
        -sloc
              print the source lines of code
              
        -o
              print the occurrences of a particular letter, word or phrase
              
        -L, -max-line-length
              print the length of the longest line
              
        -w, -words
              print the word counts
`
	version_text = `
    wc (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

// printAndExit will print a message and then exit the program.
func printAndExit(message string) {
	fmt.Println(message)
	os.Exit(0)
}

// openFile will open the file and check for errors, then return the file
func openFile(s string) *os.File {
	fi, err := os.Open(s)
	if err != nil {
		printAndExit("wc: " + s + ": No such file or directory")
	}
	return fi
}

// maxStrLength loops through each line to find the longest line.
func maxStrLength(buffer *bytes.Buffer) int {
	var maxStringLength int
	for _, line := range strings.Split(buffer.String(), "\n") {
		if len(line) > maxStringLength {
			maxStringLength = len(line)
		}
	}
	return maxStringLength
}

// bufferToStringArray returns the buffer as a newline-separated string slice.
func bufferToStringArray(buffer *bytes.Buffer) []string {
	return strings.Split(buffer.String(), "\n")
}

// isEmpty checks if the line is empty, and returns true if true
func isEmptyLine(line string) bool {
	if len(line) < 1 {
		return true
	}
	return false
}

/* removeSpacing removes tabs and spaces from the input line. This is useful
 * for properly detecting the correct SLOC. */
func removeSpacing(line *string) {
	*line = strings.Replace(*line, "\t", "", -1)
	*line = strings.Replace(*line, " ", "", -1)
}

/* isCode checks if the line is code, and returns true if true. It is
 * currently optimized for counting SLOC in Go programs. */
func isCode(line string) bool {
	removeSpacing(&line)
	prefix := strings.HasPrefix // make an alias for 'HasPrefix'
	return prefix(line, "//") || prefix(line, " *") || prefix(line, "/*") ||
		prefix(line, "*/") || prefix(line, "{") && len(line) == 1 ||
		prefix(line, "}") && len(line) == 1 || prefix(line, ")") ||
		prefix(line, "package") || prefix(line, "import") ||
		prefix(line, "var (") || prefix(line, "const (")
}

// slocCounter counts the source lines of code
func slocCounter(buffer *bytes.Buffer, count int) int {
	for _, line := range bufferToStringArray(buffer) {
		if !isEmptyLine(line) && !isCode(line) {
			count++
		}
	}
	return count
}

// occurrenceCounter counts the number of occurrences of occurrenceRef.
func occurrenceCounter(buffer *bytes.Buffer) int {
	return strings.Count(buffer.String(), *occurrenceRef)
}

// lineCount returns the number of lines by splitting the buffer's newlines.
func lineCount(buffer *bytes.Buffer) int {
	return strings.Count(buffer.String(), "\n")
}

// wordcount returns the number of words by splitting the buffer's spaces/fields.
func wordCount(buffer *bytes.Buffer) int {
	return len(strings.Fields(buffer.String()))
}

// bufferProcessor will print information relating to the input flag.
func bufferProcessor(fileName string, buffer *bytes.Buffer) {
	switch {
	case *countBytes || *countBytesL:
		fmt.Println(buffer.Len(), fileName)
	case *countCharacters || *countCharactersL:
		fmt.Println(buffer.Len(), fileName)
	case *countLines || *countLinesL:
		fmt.Println(lineCount(buffer), fileName)
	case *maxLineLength || *maxLineLengthL:
		fmt.Println(maxStrLength(buffer), fileName)
	case *countWords || *countWordsL:
		fmt.Println(wordCount(buffer), fileName)
	case *countSLOC:
		fmt.Println(slocCounter(buffer, 0), fileName)
	case len(*occurrenceRef) != 0: // Count occurences if not empty.
		fmt.Println(occurrenceCounter(buffer), fileName)
	default: // Print all if no argument is given.
		fmt.Println(lineCount(buffer), wordCount(buffer), buffer.Len(),
			fileName)
	}
}

func main() {
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		buffer := bytes.NewBuffer(nil) // create a buffer
		io.Copy(buffer, os.Stdin)      // copy stdin into the buffer
		bufferProcessor("", buffer)    // print the results
	} else {
		for file := 0; file < flag.NArg(); file++ {
			buffer := bytes.NewBuffer(nil)
			io.Copy(buffer, openFile(flag.Arg(file)))
			bufferProcessor(flag.Arg(file), buffer)
		}
	}
}

func init() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()
	if *help {
		printAndExit(help_text)
	}
	if *version {
		printAndExit(version_text)
	}
}
