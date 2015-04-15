package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

var (
	skipFirst  = flag.Bool("1", false, "suppress column 1 (lines unique to FILE1)")
	skipSecond = flag.Bool("2", false, "suppress column 2 (lines unique to FILE2)")
	skipBoth   = flag.Bool("3", false, "suppress column 3 (lines that appear in both files)")
	needSort   = flag.Bool("u", false, "sort")
)

func printCol(s string, col int) {
	switch col {
	case 1:
		if *skipFirst {
			return
		}
	case 2:
		if *skipSecond {
			return
		}
		if !*skipFirst {
			fmt.Print("\t")
		}
	case 3:
		if *skipBoth {
			return
		}
		if !*skipFirst {
			fmt.Print("\t")
		}
		if !*skipSecond {
			fmt.Print("\t")
		}
	}
	fmt.Println(s)
}

func compare(fileline [][]string) {
	var i, j int
	f1 := fileline[0]
	f2 := fileline[1]
	for i < len(f1) && j < len(f2) {
		if f1[i] < f2[j] {
			printCol(f1[i], 1)
			i++
		} else if f1[i] > f2[j] {
			printCol(f2[j], 2)
			j++
		} else {
			printCol(f1[i], 3)
			i++
			j++
		}
	}
	for ; i < len(f1); i++ {
		printCol(f1[i], 1)
	}
	for ; j < len(f2); j++ {
		printCol(f2[j], 2)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.PrintDefaults()
		return
	}

	fileline := make([][]string, 2)
	var stdin bool

	for i, arg := range flag.Args() {
		var r io.Reader
		if arg == "-" && !stdin {
			stdin = true
			r = os.Stdin
		} else {
			file, err := os.Open(arg)
			if err != nil {
				flag.PrintDefaults()
				return
			}
			defer file.Close()
			r = file
		}
		s := bufio.NewScanner(r)
		for s.Scan() {
			fileline[i] = append(fileline[i], s.Text())
		}
		if *needSort {
			sort.Strings(fileline[i])
		}
	}
	compare(fileline)
}
