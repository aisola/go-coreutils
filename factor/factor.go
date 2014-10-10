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

type factorList []int

// toString returns the factorList as a string.
func (numbers *factorList) toString() string {
	var buffer bytes.Buffer
	for _, number := range *numbers {
		buffer.WriteString(" " + strconv.Itoa(number))
	}
	return buffer.String()
}

/* getFactorList generates a factorList type containing all of the prime factors
 * of the 'number' input. The prime seive used will first check if the number is
 * divisible by the index, and if true, appends the index to the factor list. If
 * the 'index' value is greater than the square root of the 'number', the value
 * of 'number' is also a prime factor. If the square root of 'number' is equal
 * to 'index', then we can add 'index' twice and exit. If the value of 'index'
 * is two, change it to one so that we will only start checking for odd numbers
 * as prime candidates. */
func getFactorList(number int) factorList {
	var factors factorList
	for index := 2; index <= number; index += 2 {
		if number%index == 0 {
			factors = append(factors, index)
			number /= index
			index = 0
		} else if index*index > number {
			factors = append(factors, number)
			break
		} else if index*index == number {
			factors = append(factors, index)
			factors = append(factors, index)
			break
		}
		if index == 2 {
			index = 1
		}
	}
	return factors
}

/* getNumber parses the input number in string format and returns the value
 * as a number if it really is a number -- else returns an error. */
func getNumber(currentNumber string) int {
	number, err := strconv.Atoi(currentNumber)
	if err != nil {
		fmt.Printf("factor: '%s' is not a valid positive integer\n",
			flag.Arg(0))
		os.Exit(0)
	}
	return number
}

func main() {
	if flag.NArg() == 0 {
		var number int
		for {
			fmt.Scan(&number)
			factors := getFactorList(number)
			fmt.Print(number, ":", factors.toString(), "\n")
		}
	} else {
		for index := 0; index < flag.NArg(); index++ {
			number := getNumber(flag.Arg(index))
			factors := getFactorList(number)
			fmt.Print(number, ":", factors.toString(), "\n")
		}
	}
}

func init() {
	flag.Parse()
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}
}
