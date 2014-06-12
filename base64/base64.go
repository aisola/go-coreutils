//
// base64.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Trey Tacon, Abram C. Isola
//
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	help_text = `
    Usage: base64 [option]... [file]
       or: base64 --decode [option]... [file]

    transform data (from file or stdin) into (or from) base64 encoded form

      --help                 Display this message.
      ---version             Display version information.
      -D, --decode           Change the mode of operation, from the default of
                             encoding data, to decoding data. Input is expected
                             to be base64 encoded data, and the output will be
                             the original data.
      -w, --wrap             During encoding, wrap lines after cols characters.
                             This must be a positive number. The default of 0
                             disables wrapping.
      -i, --ignore-garbage   During decoding, ignore unrecognized bytes.
    `
	version_text = `
    base64 (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	ignoreGarbage = flag.Bool("ignore-garbage", false, "enables additional output to stderr")
	decode        = flag.Bool("decode", false, "decodes input")
	wrap          = flag.Int("wrap", 0, "wrap lines after 'wrap' columns")
	help          = flag.Bool("help", false, help_text)
	version       = flag.Bool("version", false, version_text)
)

func init() {
	flag.BoolVar(decode, "D", false, "decodes input")
	flag.IntVar(wrap, "w", 0, "wraps lines after 'wrap' columns")
	flag.BoolVar(ignoreGarbage, "i", false, "ignore unrecognized bytes")
}

func main() {
	flag.Parse()

	if *help {
		fmt.Println(help_text)
		return
	}

	if *version {
		fmt.Println(version_text)
		return
	}

	var (
		bytes []byte
		err   error
	)

	if flag.NArg() > 0 {
		bytes, err = ioutil.ReadFile(flag.Arg(0))
	} else {
		bytes, err = ioutil.ReadAll(os.Stdin)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var (
		dst     = make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
		encFunc = base64.StdEncoding.Encode
	)

	if *decode {
		dst = make([]byte, base64.StdEncoding.DecodedLen(len(bytes)))
		encFunc = func(dst, src []byte) {
			_, err := base64.StdEncoding.Decode(dst, src)
			if err != nil {
				fmt.Fprintln(os.Stderr)
				os.Exit(1)
			}
		}
	}

	encFunc(dst, bytes)

	dstString := string(dst)
	if *wrap == 0 {
		fmt.Println(dstString)
		return
	}

	var i int
	for i = 0; i+*wrap < len(dstString); i += *wrap {
		fmt.Println(dstString[i : i+*wrap])
	}
	if i < len(dstString) {
		fmt.Println(dstString[i:])
	}
}
