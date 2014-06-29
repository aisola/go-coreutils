//
// factor.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//
package main

import "bytes"
import "flag"
import "fmt"
import "strconv"
import "os"

const (
	help_text = `
    Usage: factor [NUMBER]...
    
    Print the prime factors of each specified integer number. If none are
    specified on the command line, read them from standard output.
    
    -help display this help and exit
    
    -version output version information and exit
`
	version_text = `
    factor (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for deprintTails see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	help    = flag.Bool("help", false, "display help information")
	version = flag.Bool("version", false, "display version information")
)

func processFlags() {
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}
}

// Returns the integer of the current argument.
func getNumber(currentNumber string) int {
	number, err := strconv.Atoi(currentNumber)

	// Check if the number to factor is actually a number.
	if err != nil {
		fmt.Printf("factor: '%s' is not a valid positive integer\n", flag.Arg(0))
		os.Exit(0)
	}

	return number
}

// Returns a slice of all the prime factors associated with a number.
func getFactors(number int) []int {
	factors := make([]int, 0)

	// If the number is
	for index := 2; index < number; index++ {
		if number%index == 0 {
			factors = append(factors, index)
			number = number / index
			index = 1
		}
	}

	// Append the final prime number to the factors slice.
	factors = append(factors, number)

	return factors
}

// Returns a string of the factor slice.
func factorsToString(numbers []int) string {
	var buffer bytes.Buffer

	for _, number := range numbers {
		buffer.WriteString(" " + strconv.Itoa(number))
	}

	return buffer.String()
}

// Prints factors for each argument given.
func printFactors() {
	for index := 0; index < flag.NArg(); index++ {
		number := getNumber(flag.Arg(index))
		factors := getFactors(number)
		fmt.Print(number, ":", factorsToString(factors), "\n")
	}
}

func main() {
	flag.Parse()
	processFlags()
	printFactors()
}
