package main

import "os"
import "log"

// Get PID of Parent
var process = os.Getppid()

func main() {
	pproc, err := os.FindProcess(process)

	if err != nil {
		log.Fatalln(err)
	} else {
		pproc.Kill()
	}
}
