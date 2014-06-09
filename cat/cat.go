//
// cat.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "net"
import "os"

const version_text = `
    cat (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`

var (
	countNonBlank     = flag.Bool("b", false, "Number the non-blank output lines, starting at 1.")
	numberOutput      = flag.Bool("n", false, "Number the output lines, starting at 1.")
	squeezeEmptyLines = flag.Bool("s", false, "Squeeze multiple adjacent empty lines, causing the output to be single spaced.")
)

func openFile(s string) (io.ReadWriteCloser, error) {
	fi, err := os.Stat(s)
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeSocket != 0 {
		return net.Dial("unix", s)
	}
	return os.Open(s)
}

func dumpLines(w io.Writer, r io.Reader) (n int64, err error) {
	var lastline, line string
	br := bufio.NewReader(r)
	nr := 0
	for {
		line, err = br.ReadString('\n')
		if err != nil {
			return
		}
		if *squeezeEmptyLines && lastline == "\n" && line == "\n" {
			continue
		}
		if *countNonBlank && line == "\n" || line == "" {
			fmt.Fprint(w, line)
		} else if *countNonBlank || *numberOutput {
			nr++
			fmt.Fprintf(w, "%6d\t%s", nr, line)
		} else {
			fmt.Fprint(w, line)
		}
		lastline = line
	}
	return
}

func main() {
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	rcopy := io.Copy
	if *countNonBlank || *numberOutput || *squeezeEmptyLines {
		rcopy = dumpLines
	}

	for _, fname := range flag.Args() {
		if fname == "-" {
			rcopy(os.Stdout, os.Stdin)
		} else {
			f, err := openFile(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			rcopy(os.Stdout, f)
			f.Close()
		}
	}

}
