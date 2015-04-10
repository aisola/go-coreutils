//
// tee.go (go-coreutils) 0.1
// Copyright (C) 2015, The GO-Coreutils Developers.
//
// Written By: Haruki Tsurumoto
//
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	helpText = `
  Usage: tee [OPTION]... [FILE]...
  Copy standard input to each FILE, and also to standard output.

    -a, --append              append to the given FILEs, do not overwrite
    -i, --ignore-interrupts   ignore interrupt signals
        --help     display this help and exit
        --version  output version information and exit

  If a FILE is -, copy again to standard output.
    `
	versionText = `
    tee (go-coreutils) 0.1
    Copyright (C) 2015, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute
    it under certain conditions in LICENSE.
`
)

var (
	help                 = flag.Bool("help", false, "help")
	version              = flag.Bool("version", false, "version_text")
	appendOpt            = flag.Bool("a", false, "append to the given FILEs, do not overwrite")
	appendOptLong        = flag.Bool("append", false, "append to the given FILEs, do not overwrite")
	ignoreInterrupts     = flag.Bool("i", false, "ignore interrupt signals")
	ignoreInterruptsLong = flag.Bool("ignore-interrupts", false, "ignore interrupt signals")
)

func main() {
	overwrite := true
	ignoreSigInt := false
	flag.Parse()
	if *help {
		fmt.Println(helpText)
		os.Exit(0)
	}
	if *version {
		fmt.Println(versionText)
		os.Exit(0)
	}
	if *appendOpt || *appendOptLong {
		overwrite = false
	}
	if *ignoreInterrupts || *ignoreInterruptsLong {
		ignoreSigInt = true
	}
	// open Files
	arg := flag.Args()
	var files []*os.File
	if len(arg) >= 1 && arg[0] == "-" {
		files = append(files, os.Stdout)
	} else {
		for i := range arg {
			var f *os.File
			var err error
			if overwrite {
				f, err = os.OpenFile(arg[i], os.O_WRONLY|os.O_CREATE, 0644)
			} else {
				f, err = os.OpenFile(arg[i], os.O_WRONLY|os.O_APPEND, 0644)
			}
			if err != nil {
				f.Close()
				continue
			}
			defer f.Close()
			files = append(files, f)
		}
	}
	// signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		for {
			<-signalChan
			if !ignoreSigInt {
				os.Exit(1)
			}
		}
	}()
	// Main
	buffer := make([]byte, 1024)
	for {
		n1, err := os.Stdin.Read(buffer)
		if err != nil {
			en, ok := err.(syscall.Errno)
			if ok && int(en) == 0x4 {
				// 0x4 == EINTR: interrupted system call
				continue
			} else {
				break
			}
		}
		os.Stdout.Write(buffer[0:n1])
		for _, f := range files {
			f.Write(buffer[0:n1])
		}
		if n1 <= 0 {
			break
		}
	}
}
