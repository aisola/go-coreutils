//
// echo.go (go-coreutils) 0.1
// Copyright (C) 2014, Abram C. Isola.
//
package main

import "os"
import "fmt"
import "flag"
import "strings"

const version_text = `
    echo (go-coreutils) 0.1

    Copyright (C) 2014 Abram C. Isola.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: echo [options] [string ...]")
}

func main() {
	enableEscapeChars := flag.Bool("e", false, "Enable escape characters")
	omitNewline := flag.Bool("n", false, "Don't print trailing newline")
	disableEscapeChars := flag.Bool("E", true, "Disable escape characters")
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	concatenated := strings.Join(flag.Args(), " ")

	a := []rune(concatenated)

	length := len(a)

	ai := 0

	if length != 0 {
		for i := 0; i < length; {
			c := a[i]
			i++
			if (*enableEscapeChars == true || *disableEscapeChars == false) && c == '\\' && i < length {
				c = a[i]
				i++
				switch c {
				case 'a':
					c = '\a'
				case 'b':
					c = '\b'
				case 'c':
					os.Exit(0)
				case 'e':
					c = '\x1B'
				case 'f':
					c = '\f'
				case 'n':
					c = '\n'
				case 'r':
					c = '\r'
				case 't':
					c = '\t'
				case 'v':
					c = '\v'
				case '\\':
					c = '\\'
				case 'x':
					c = a[i]
					i++
					if '9' >= c && c >= '0' && i < length {
						hex := (c - '0')
						c = a[i]
						i++
						if '9' >= c && c >= '0' && i <= length {
							c = 16*(c-'0') + hex
						}
					}
				}
			}
			a[ai] = c
			ai++
		}
	}

	os.Stdout.WriteString(string(a[:ai]))
	if *omitNewline == false {
		fmt.Print("\n")
	}
}
