// yes.go (go-coreutils) 0.1
//
// yes - output a string repeatedly until killed
package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	var opts = flag.Args()
	if len(opts) == 0 {
		opts = []string{"y"}
	}

	for {
		fmt.Println(opts[0])
	}
}
