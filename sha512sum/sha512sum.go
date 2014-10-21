//
// sha512sum.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "bufio"
import "crypto/sha512"
import "flag"
import "fmt"
import "io"
import "io/ioutil"
import "os"
import "strings"

const (
	help_text string = `
    Usage: sha512sum [OPTION] [FILE]...
       or: sha512sum [OPTION] --check [FILE]...

    Print or check sha512 checksums.
    If FILE is not given or is -, read standard input.

        --help        display this help and exit
        --version     output version information and exit

        -c, --check   check sha512 sums against given list

    The sums are computed as described in RFC 1321. When checking, the input
    should be a former output of this program. The default mode is to print
    a line with a checksum, a character indicating type ('*' for binary, ' ' for
    text), and name for each FILE.
    `
	version_text = `
    sha512sum (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute
    it under certain conditions in LICENSE.
`
)

func main() {
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	check := flag.Bool("check", false, "check sha512 sums against given list")
	check1 := flag.Bool("c", false, "check sha512 sums against given list")
	flag.Parse()

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	// If you are NOT checking...
	if !*check && !*check1 {
		if flag.NArg() > 0 {
			for _, file := range flag.Args() {
				buff, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Printf("sha512sum: cannot read '%s': %s\n", file, err)
				}
				fmt.Printf("%x %s\n", sha512.Sum512(buff), file)
			}
		} else {
			buff, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Printf("sha512sum: cannot read 'STDIN': %s\n", err)
			}
			fmt.Printf("%x -\n", sha512.Sum512(buff))
		}

		// Check the files...
	} else {
		if flag.NArg() > 0 {
			NUMFAILED := 0
			for _, file := range flag.Args() {
				fp, err := os.Open(file)
				if err != nil {
					fmt.Printf("sha512sum: cannot read '%s': %s\n", file, err)
				}
				// Set up new reader
				bf := bufio.NewReader(fp)

				// loop through lines
				for {
					// get info
					line, isPrefix, err := bf.ReadLine()

					// is if EOF?
					if err == io.EOF {
						break
					}
					// another ERR?
					if err != nil {
						fmt.Printf("sha512sum: cannot read '%s': %s\n", file, err)
					}
					// is the line WAY too long?
					if isPrefix {
						fmt.Printf("sha512sum: unexpected long line: %s\n", file)
					}

					// success. check the hash
					hashfile := strings.Split(string(line), " ")
					HASH := string(hashfile[0])
					HFIL := string(hashfile[1])
					buff, err := ioutil.ReadFile(HFIL)
					if err != nil {
						fmt.Printf("sha512sum: cannot read '%s': %s\n", HFIL, err)
					}
					if HASH == fmt.Sprintf("%x", sha512.Sum512(buff)) {
						fmt.Printf("%s: OK\n", HFIL)
					} else {
						fmt.Printf("%s: FAILED\n", HFIL)
						NUMFAILED += 1
					}
				}
			}
			// Print how many TOTAL failed...
			if NUMFAILED > 0 {
				fmt.Printf("sha512sum: WARNING: %d computed checksum did NOT match\n", NUMFAILED)
			}

		} /* TODO: Implement this section...
		 else {
			buff, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Printf("sha512sum: cannot read 'STDIN': %s", err)
			}
			fmt.Printf("%x -\n", sha512.Sum(buff))
		} */
	}
}
