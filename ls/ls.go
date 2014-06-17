//
// ls.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy, Abram C. Isola
//

package main

import "fmt"
import "io/ioutil"
import "os"
import "strings"
import "flag"
import "unsafe"
import "syscall"
import "time"

const ( // Constant variables used throughout the program.
	TERMINAL_INFO    = 0x5413         // Used in the getTerminalWidth function
	EXECUTABLE       = 0111           // File executable bit
	SYMLINK          = os.ModeSymlink // Symlink bit
	CYAN_SYMLINK     = "\x1b[36;1m"   // Cyan terminal color
	BLUE_DIR         = "\x1b[34;1m"   // Blue terminal color
	GREEN_EXECUTABLE = "\x1b[32;1m"   // Green terminal color
	RESET            = "\x1b[0m"      // Reset terminal color
	SPACING          = 1              // Spacing between columns

	help_text string = `
    Usage: ls [OPTION]...
    
    list files and directories in working directory

        --help        display this help and exit
        --version     output version information and exit
        
        -a  include hidden files and directories
        -l  use a long listing format
        -1  list in a single column
    `
	version_text = `
    ls (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

var ( // Default flags and variables.
	showHidden      = flag.Bool("a", false, "list hidden files and directories")
	oneColumn       = flag.Bool("1", false, "list files by one column")
	longMode        = flag.Bool("l", false, "use a long listing format")
	isHidden        = false                   // Sets whether or not a file should be hidden.
	printOneLine    = true                    // Sets whether or not to print on one row.
	terminalWidth   = int(getTerminalWidth()) // Grabs information on the current terminal width.
	maxCharLength   = 0                       // Statistics for the largest file name length.
	totalCharLength = 0                       // Statistics for the total number of characters.
)

func main() {
	fileNameList := make([]string, 0) // Stores a list of files with their unaltered name.
	fileLengthList := make([]int, 0)  // Stores a list of character lengths per file.
	//fileModeList := make([]string, 0)
	//fileUserList := make([]string, 0)
	//fileGroupList := make([]string, 0)
	fileModDateList := make([]time.Time, 0)
	fileSizeList := make([]int64, 0)
	
	help := flag.Bool("help", false, help_text)
	version := flag.Bool("version", false, version_text)
	flag.Parse()      // Parse flags
	path := getPath() // Obtains the path to search

	// Display help information

	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	// Display version information

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	/* Scans the directory and returns a list of the contents. If the directory
	 * does not exist, an error is printed and the program exits. */

	directory, err := ioutil.ReadDir(path)
	errorChecker(&err, "ls: " + path + " - No such file or directory.\n")

	/* This for loop will loop through each file in the directory in ascending order
	 * and call the filterHidden function to determine whether or not it should
	 * display the file. */

	var colorized string      // Stores the colorized name
	var currentFileLength int // Stores the current file length
	for _, file := range directory {
		filterHidden(file.Name()) // Check to see if the file is to be hidden or not.
		if isHidden == false {
			/* If it should be displayed, it will call the colorizer function to
			 * determine what kind of file the file is (Directory, Symlink, Executable,
			 * or File) and colorize the file based on that. */

			colorized = colorizer(file.Name(), // Send current file to the colorizer function
				file.IsDir(),                         // Check if the file is a directory
				file.Mode()&SYMLINK == SYMLINK,       // Check if it is a symbolic link
				file.Mode()&EXECUTABLE == EXECUTABLE, // Check if it is an executable
			)

			/* Rather than print directly to the terminal, we store name and length of
			 * the name in slices for later use. */

			currentFileLength = len(file.Name())                       // Obtains character length of the current file
			fileNameList = append(fileNameList, colorized)             // append the colorized file to the list
			fileLengthList = append(fileLengthList, currentFileLength) // append the name

			/* Finally, this is a good spot for getting some statistics on the totalCharLength
			 * and maxCharLength of the files. */

			if totalCharLength <= terminalWidth { // Determines if we can print on one line.
				totalCharLength += currentFileLength + 2 // The additional 2 is for spacing.
			} else {
				printOneLine = false
			}

			if currentFileLength > maxCharLength { // Determines the longest file length.
				maxCharLength = currentFileLength
			}

			// If longMode is enabled, we need to grab extra statistics.
			if *longMode {
				//fileUserList
				//fileGroupList
				//fileModeList = append(fileModeList, file.Mode())
				fileModDateList = append(fileModDateList, file.ModTime())
				fileSizeList = append(fileSizeList, file.Size())
			}
		}
	}

	// This switch will determine how we should print.

	switch {
	case *longMode: // If longMode is enabled
		for index, _ := range fileNameList {
			longModePrinter(fileNameList[index],
				//fileModeList[index],
				//fileUserList[index],
				//fileGroupList[index],
				fileSizeList[index],
				fileModDateList[index],
			)
		}
	case printOneLine: // If we can print everything on one row, do it.
		for _, file := range fileNameList {
			fmt.Print(file, "  ") // Print the file plus additional spacing
		}
		fmt.Println(RESET) // Print an additional line and reset the color.
	case *oneColumn: // If the -1 flag is set, print everything in one column.
		for _, file := range fileNameList {
			fmt.Println(file) // Print the file and create a new line
			fmt.Print(RESET)  // Same as above, except we do not need a new line here.
		}
	default: // Properly spaces our files and then prints from top to bottom.
		spaced := make([]string, 0)
		for index, file := range fileNameList { // Preprocesses the file list for printing by adding spaces.
			spaced = append(spaced, spacer(&file, fileLengthList[index]))
		}
		printTopToBottom(spaced) // Let's get printing!
		fmt.Println(RESET)       // Reset terminal at the end.
	}
}

type termsize struct { // Stores information regarding the terminal size.
	Row, Col, Xpixel, Ypixel uint16
}

func getTerminalWidth() uint { // Obtains the current width of the terminal.
	ws := &termsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(TERMINAL_INFO),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col)
}

func errorChecker(err *error, message string) {
	if *err != nil {
		fmt.Print(message)
		os.Exit(0)
	}
}

// If there is no argument, set the directory path to the current working directory

func getPath() string {
	args := flag.Args()
	
	if len(args) < 1 {
		path, err := os.Getwd()
		errorChecker(&err, "ls: Could not obtain the current working directory.\n")
		return path
	} else {
		return args[0]
	}
}

/* The filterHidden function will check if the hidden flag is set, and if not, whether
 * the file is a hidden file or not. If the flag is not set and the file is a hidden
 * file (starts with a period), then it will not be printed. */

func filterHidden(file string) {
	switch {
	case *showHidden: // If the hidden flag is enabled
		isHidden = false
	case strings.HasPrefix(file, "."): // If it is disabled and the file is hidden
		isHidden = true
	default: // If it is disabled and the file is not hidden
		isHidden = false
	}
}

/* If the file is a symbolic link, print it in cyan; a directory, blue; an executable file,
 * green; else print the file in white. */

func colorizer(file string, isDir, isSymlink, isExecutable bool) string {
	switch {
	case isSymlink:
		return CYAN_SYMLINK + file
	case isDir:
		return BLUE_DIR + file
	case isExecutable:
		return GREEN_EXECUTABLE + file
	default:
		return RESET + file
	}
}

/* The spacer function will add spaces to the end of each file name so that they line up
 * correctly when printing. */

func spacer(name *string, charLength int) string {
	return string(*name + strings.Repeat(" ", maxCharLength-charLength+SPACING))
}

/* The countRows function counts the number of rows that will be printed. The number is
 * determined by dividing the number of files by the maximum number of columns. However,
 * if there is a remaindera of files left over for an incomplete row, we add an additional
 * row. */

func countRows(lastRowCount, maxColumns, numOfFiles *int) int {
	if *lastRowCount == 0 {
		return *numOfFiles / *maxColumns
	}
	return *numOfFiles / *maxColumns + 1 // Add additional row if the last row is incomplete.
}

/* Finally, the printTopToBottom function will take the spaced slice and print all the
 * files from top to bottom. */

func printTopToBottom(spaced []string) {
	/* The first thing that we need to know is how many columns we will be printing, how
	 * many files are left over to be printed on the last row, and the number of rows that
	 * we will be printing. */

	numOfFiles := len(spaced)                                       // Total number of files
	maxColumns := terminalWidth / (maxCharLength + SPACING)         // Number of columns to print
	lastRowCount := numOfFiles % maxColumns                         // Number of files on the last row
	numOfRows := countRows(&lastRowCount, &maxColumns, &numOfFiles) // Number of rows to print
	printing := true                                                // Turn the printer on by default
	var currentRow, currentIndex int = 1, 0

	/* This our magnificent printing press. It will first start on the default case, which
	 * will print the majority of the files from the first row all the way to the next to
	 * last row. The tricky party is trying to figure out an algorithm that allows us to
	 * print from top to bottom. To do this, we need to know how many files are left over
	 * on the last row.
	 *
	 * For example, say you have three rows, eight columns, and the last row has four files.
	 * To get everything to play nicely together, we start by looping through the file list
	 * in intervals based on the number of rows (three in our example). However, at some point,
	 * the last row is going to kill our algorithm. To counter-act that, we need to switch it
	 * to counting in intervals of two after the fifth file has been printed on each row.
	 *
	 * Once we get to the last row, we can simply print everyting in an interval based on the
	 * number of rows until we have printed the last file. */

	for printing { // Keep printing as long as the printer is turned on.
		if currentRow < numOfRows { // Prints all rows except the last row.
			for column := 1; column < maxColumns; column++ { // Prints all but the last column.
				fmt.Print(spaced[currentIndex]) // Print the file.
				if column >= lastRowCount+1 {   // This is where we see if it's time to switch our index
					currentIndex += numOfRows - 1 // interval to one less than then initial interval value.
				} else {
					currentIndex += numOfRows // Select the next file in the index based on the row interval.
				}
			}
			fmt.Println(spaced[currentIndex]) // Prints the final column in a row.
			currentRow++                      // It's time to start printing the next row.
			currentIndex = currentRow - 1     // We need to reset this for the next row.
		} else { // Prints the final row.
			for index := 1; index <= lastRowCount; index++ { // This is the final print run.
				fmt.Print(spaced[currentIndex]) // Print the file.
				currentIndex += numOfRows       // Switch the index to the next file.
			}
			printing = false // We are finished printing -- turn the printing press off.
		}
	}
}

// The longModePrinter function prints files and their statistics one line at a time.

func longModePrinter(fileName string, fileSize int64, fileDate time.Time) {
	fmt.Printf("%10d %s %s\n", fileSize, fileDate, fileName+RESET)
}
