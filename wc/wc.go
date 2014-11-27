//
// wc.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//
package main

import "bufio"
import "bytes"
import "flag"
import "fmt"
import "os"
import "strings"
import "unicode/utf8"

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

// slocCounter counts the source lines of code
func slocCounter(buffer []byte, count int) int {
	// Returns true if the input line is not an empty.
	isNotEmpty := func(buffer []byte) bool {
		return !(len(buffer) < 1)
	}
	// Returns true if the input line is a line of code
	isCode := func(line string) bool {
		line = strings.Replace(line, "\t", "", -1)
		line = strings.Replace(line, " ", "", -1)
		prefix := strings.HasPrefix // make an alias for 'HasPrefix'
		return prefix(line, "//") || prefix(line, " *") ||
			prefix(line, "/*") || prefix(line, "*/") ||
			prefix(line, "{") && len(line) == 1 ||
			prefix(line, "}") && len(line) == 1 ||
			prefix(line, ")") || prefix(line, "package") ||
			prefix(line, "import") || prefix(line, "var (") ||
			prefix(line, "const (")
	}
	if isNotEmpty(buffer) && !isCode(string(buffer)) {
		count++
	}
	return count
}

// occurrenceCounter counts the number of occurrences of occurrenceRef.
func occurrenceCounter(buffer []byte) int {
	return bytes.Count(buffer, []byte(*occurrenceRef))
}

// wordcount returns the number of words by splitting the buffer's fields.
func wordCount(buffer []byte) int {
	return len(strings.Fields(string(buffer)))
}

/* characterCount counts the number of characters in the bytes buffer by
 * sending a bytes slice of the buffer to the RuneCount function in utf8. This
 * gives support for non-ASCII characters that don't fit inside of a byte. */
func characterCount(buffer []byte) int {
	return utf8.RuneCount(buffer)
}

// wcstat stores statistics for each file processsed by wc
type wcstat struct {
	bytes      int
	characters int
	lines      int
	maxLength  int
	words      int
	sloc       int
	occurences int
	fileName   string
}

// maxLineLength determines the maximum line length for the entire file.
// NOTE: GNU wc counts by bytes rather than by runes, which is not accurate. It
// also does not correctly detect tabs, giving a false line length size.
func (wc *wcstat) maxLineLength(buffer []byte) {
	length := utf8.RuneCount(buffer)
	if wc.maxLength < length {
		wc.maxLength = length
	}
}

// processStats will process the input buffer and append the obtained stats
// to the wcstat struct.
func (wc *wcstat) getStats(buffer []byte) {
	switch {
	case *countBytes || *countBytesL:
		wc.bytes += len(buffer) + 1
	case *countCharacters || *countCharactersL:
		wc.characters += characterCount(buffer) + 1
	case *countLines || *countLinesL:
		wc.lines++
	case *maxLineLength || *maxLineLengthL:
		wc.maxLineLength(buffer)
	case *countWords || *countWordsL:
		wc.words += wordCount(buffer)
	case *countSLOC:
		wc.sloc += slocCounter(buffer, 0)
	case len(*occurrenceRef) != 0: // Count occurences if not empty.
		wc.occurences += occurrenceCounter(buffer)
	default: // Print all if no argument is given.
		wc.lines++
		wc.words += wordCount(buffer)
		wc.bytes += len(buffer)
	}
}

// printStats prints the statistics for the current file.
func (wc *wcstat) printStats() {
	switch {
	case *countBytes || *countBytesL:
		fmt.Println(wc.bytes, wc.fileName)
	case *countCharacters || *countCharactersL:
		fmt.Println(wc.characters, wc.fileName)
	case *countLines || *countLinesL:
		fmt.Println(wc.lines, wc.fileName)
	case *maxLineLength || *maxLineLengthL:
		fmt.Println(wc.maxLength, wc.fileName)
	case *countWords || *countWordsL:
		fmt.Println(wc.words, wc.fileName)
	case *countSLOC:
		fmt.Println(wc.sloc, wc.fileName)
	case len(*occurrenceRef) != 0: // Count occurences if not empty.
		fmt.Println(wc.occurences, wc.fileName)
	default: // Print all if no argument is given.
		fmt.Println(wc.lines, wc.words, wc.bytes+wc.lines, wc.fileName)
	}
}

// scanFile scans each file, line by line, gathering statistics using a scanner.
func (wc *wcstat) scanFile(scanner *bufio.Scanner) {
	for scanner.Scan() {
		wc.getStats(scanner.Bytes())
	}
}

func main() {
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		var wc wcstat
		wc.scanFile(bufio.NewScanner(os.Stdin))
		wc.printStats()
	} else {
		for file := 0; file < flag.NArg(); file++ {
			wc := wcstat{fileName: flag.Arg(file)}
			wc.scanFile(bufio.NewScanner(openFile(flag.Arg(file))))
			wc.printStats()
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
