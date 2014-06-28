//
// expr.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//
// TODO: Implement & and | expressions.
package main

import "flag"
import "fmt"
import "math"
import "os"
import "strconv"
import "strings"

const (
	help_text = `
	Usage: expr EXPRESSION
	   or: expr OPTION
	   
	-help    display this help and exit
	
	-version output version information and exit
	   
	Print the value of EXPRESSION to standard output. EXPRESSIONS are listed below:
	
	<  - less than
	
	<= - less than or equal to
	
	=  - equal to
	
	=> - greater than or equal to
	
	>  - greater than
	
	+  - arithmetic sum of two arguments
	
	-  - arithmetic difference of two arguments
	
	*  - arithmetic product of two arguments
	
	/  - arithmetic quotient of two arguments
	
	%  - arithmetic remainder after dividing two arguments
	
	substr STRING STARTPOS LENGTH
	      substring of STRING, STARTPOS counted from 1
	      
	index STRING CHAR
	      index in STRING where any CHAR is found, or 0
	      
	length STRING
	      length of STRING
`
	version_text = `
    expr (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var (
	help    = flag.Bool("help", false, "display help information")
	version = flag.Bool("version", false, "display version information")
)

// Check initial state of flags.
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

// Prints an error indicating that the syntax is wrong and exits.
func printError() {
	fmt.Println("expr: syntax error")
	os.Exit(0)
}

// Print an error and exit the program if there are no arguments.
func checkIfNoArgumentsAreGiven() {
	if flag.NArg() == 0 {
		fmt.Println("expr: missing operand\nTry 'expr -help' for more information.")
		os.Exit(0)
	}
}

// Returns a slice of value arguments.
func getValueSlice() []float64 {
	// This will select every odd argument, which is a number.
	valueSlice := make([]float64, 0)                   // Create a slice for storing values (numbers).
	for index := 0; index <= flag.NArg(); index += 2 { // Loop through every odd argument
		value, _ := strconv.ParseFloat(flag.Arg(index), 64) // Convert the flag argument, which is a string, into float32 with strconv.
		valueSlice = append(valueSlice, value)              // Append value, which is currently a float64 value, as a float32 value to the value slice.
	}
	return valueSlice
}

// Returns a slice of modifier arguments.
func getModifierSlice() []string {
	// This will select every even argument, which is a modifier (+.-,/,*)
	modifierSlice := make([]string, 0)                // Create a slice for storing modifiers.
	for index := 1; index < flag.NArg(); index += 2 { // Loop through every even argument
		modifierSlice = append(modifierSlice, flag.Arg(index)) // Append the modifier to the modifier slice.
	}
	return modifierSlice
}

// Return true if the current modifier is an inequality symbol.
func isInequalitySymbol(test string) bool {
	var isBool bool = false
	switch test {
	case "<", "<=", "=", "=>", ">", "!=":
		isBool = true
	}
	return isBool
}

// Checks if the number can be an integer for calculating the modulus.
func floatIsInteger(number float64) bool {
	// If the float can be truncated into an int, return true.
	if math.Trunc(number) == number {
		return true
	}

	return false
}

// Checks if the numbers can be modulated, and if so returns the modulated value.
func calculateModulus(original, current float64) float64 {
	var result float64

	// Check if the numbers can be modulated
	if floatIsInteger(original) && floatIsInteger(current) {
		result = float64(int64(original) % int64(current))
	} else {
		fmt.Println("expr: non-integer argument")
		os.Exit(0)
	}

	return result
}

// Return 1 if true, 0 if false.
func booleanToFloat(boolean bool) float64 {
	if boolean {
		return 1
	} else {
		return 0
	}
}

// The initial result in an expression needs to be calculated differently than the rest.
func calculateInitialResult(firstNum, secondNum float64, modifier string) float64 {
	var result float64
	switch modifier {
	case "+":
		result = firstNum + secondNum
	case "-":
		result = firstNum - secondNum
	case "*":
		result = firstNum * secondNum
	case "/":
		result = firstNum / secondNum
	case "%":
		result = calculateModulus(firstNum, secondNum)
	default:
		printError()
	}
	return result
}

// Calculate the range of numbers to calculate between expressions.
func calculateExpressionRanges(valueSlice []float64, modifierSlice []string) ([]int, int) {
	var currentRange, equalityCount int = 0, 0
	expressionRanges := make([]int, 0)

	// Loop through each modifier and check if it is an inequality symbol, then append the range.
	for _, modifier := range modifierSlice {
		if isInequalitySymbol(modifier) {
			equalityCount++
			expressionRanges = append(expressionRanges, currentRange+1)
			currentRange = 0
		} else {
			currentRange++
		}
	}

	// The above loop will exit before appending the last range, so this will append it.
	currentRange++
	expressionRanges = append(expressionRanges, currentRange)

	return expressionRanges, equalityCount
}

// Calculate each result between inequality expressions
func calculateExpressions(valueSlice []float64, modifierSlice []string, expressionRanges []int) []float64 {
	results := make([]float64, 0)
	var result float64
	var position int = 0

	for _, currentRange := range expressionRanges {
		if currentRange == 1 {
			result = valueSlice[position]
		} else if currentRange == 2 {
			result = calculateInitialResult(valueSlice[position], valueSlice[position+1], modifierSlice[position])
		} else if currentRange > 2 {
			result = calculateInitialResult(valueSlice[position], valueSlice[position+1], modifierSlice[position])
			for index := position; index < currentRange+position-2; index++ {
				switch modifierSlice[index+1] {
				case "+":
					result += valueSlice[index+2]
				case "-":
					result -= valueSlice[index+2]
				case "*":
					result *= valueSlice[index+2]
				case "/":
					result /= valueSlice[index+2]
				case "%":
					result = calculateModulus(result, valueSlice[index+2])
				default:
					printError()
				}
			}
		}
		results = append(results, result)
		position = position + currentRange
	}

	return results
}

// Returns 1 for true, 0 for false, or any
func calculateInequalities(results []float64, expressionRanges []int, boolCount int, modifierSlice []string) float64 {
	currentValue := results[0]

	for index := 0; index < boolCount; index++ {
		currentModifier := modifierSlice[expressionRanges[index]-1]
		nextValue := results[index+1]

		switch currentModifier {
		case "<":
			currentValue = booleanToFloat(currentValue < nextValue)
		case "<=":
			currentValue = booleanToFloat(currentValue <= nextValue)
		case "=":
			currentValue = booleanToFloat(currentValue == nextValue)
		case "=>":
			currentValue = booleanToFloat(currentValue <= nextValue)
		case ">":
			currentValue = booleanToFloat(currentValue > nextValue)
		case "!=":
			currentValue = booleanToFloat(currentValue != nextValue)
		}
	}

	return currentValue
}

// Performs arithmetic calculations
func calculateArithmetic(valueSlice []float64, modifierSlice []string) {
	expressionRanges, boolCount := calculateExpressionRanges(valueSlice, modifierSlice)

	// Calculate the totals between expressions
	results := calculateExpressions(valueSlice, modifierSlice, expressionRanges)

	// Calculate inequality expressions between results or print result.
	if boolCount != 0 {
		fmt.Println(calculateInequalities(results, expressionRanges, boolCount, modifierSlice))
	} else {
		fmt.Println(results[0])
	}
}

// Returns the length of the string
func getStringLength() int {
	return len(flag.Arg(1))
}

// Returns the index value of the position of the first occurence of a character.
func getCharacterIndex() int {
	return strings.IndexByte(flag.Arg(1), flag.Arg(2)[0]) + 1
}

// Returns a substring containing only the characters in the input range.
func getSubstring() string {
	inputString := flag.Arg(1)
	start, starterr := strconv.Atoi(flag.Arg(2))
	end, enderr := strconv.Atoi(flag.Arg(3))

	// Check for errors in syntax
	if starterr != nil || enderr != nil {
		printError()
	}

	// In case the user sets a length beyond the slice size,
	// set end to the length of the input string.
	if end > len(inputString) {
		end = len(inputString)
	}

	return inputString[start-1 : end]
}

func main() {
	flag.Parse()
	processFlags()

	// If there are no arguments, print an error and exit.
	checkIfNoArgumentsAreGiven()

	// Check if length is the first argument
	switch flag.Arg(0) {
	case "match":
		//TODO
		fmt.Println("not implemented")
		os.Exit(0)
	case "substr":
		fmt.Println(getSubstring())
		os.Exit(0)
	case "index":
		fmt.Println(getCharacterIndex())
		os.Exit(0)
	case "length":
		fmt.Println(getStringLength())
		os.Exit(0)
	case "+":
		//TODO
		fmt.Println("not implemented")
	}

	// Obtain value and modifier slices
	valueSlice := getValueSlice()
	modifierSlice := getModifierSlice()

	// Determine how to process the arguments
	switch modifierSlice[0] {
	case "+", "-", "*", "/", "%", "<", "<=", "=", "=>", ">", "!=":
		calculateArithmetic(valueSlice, modifierSlice)
	default:
		printError()
	}
}
