package main

import "os"
import "fmt"
import "flag"
import "strings"

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: echo [options] [string ...]")
}

var enableEscapeChars = flag.Bool("e", false, "Enable escape characters")
var omitNewline = flag.Bool("n", false, "Don't print trailing newline")
var disableEscapeChars = flag.Bool("E", true, "Disable escape characters")

func main() {
	flag.Parse()

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
