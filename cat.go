package main

import "bufio"
import "flag"
import "fmt"
import "io"
import "net"
import "os"

var (
	countNonBlank     = flag.Bool("b", false, "Number the non-blank output lines, starting at 1.")
	numberOutput      = flag.Bool("n", false, "Number the output lines, starting at 1.")
	squeezeEmptyLines = flag.Bool("s", false,
		"Squeeze multiple adjacent empty lines, causing the output to be single spaced.")
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
	flag.Parse()
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
