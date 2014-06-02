package main

import "fmt"
import "log"
import "os"

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println(pwd)
	}

}
