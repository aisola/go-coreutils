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
import "os/user"
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
	DATE_FORMAT      = "Jan _2 15:04" // Format date
	DATE_YEAR_FORMAT = "Jan _2  2006" // If the file is from a previous year

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
	help            = flag.Bool("help", false, "display help information")
	version         = flag.Bool("version", false, "display version information")
	showHidden      = flag.Bool("a", false, "list hidden files and directories")
	singleColumn    = flag.Bool("1", false, "list files by one column")
	longMode        = flag.Bool("l", false, "use a long listing format")
	printOneLine    = true                    // Sets whether or not to print on one row.
	terminalWidth   = int(getTerminalWidth()) // Grabs information on the current terminal width.
	maxCharLength   = 0                       // Statistics for the largest file name length.
	maxIDLength     = 0                       // Statistics for the longest id name length.
	maxSizeLength   = 0                       // Statistics for the longest file size length.
	totalCharLength = 0                       // Statistics for the total number of characters.
)

// Stores information regarding the terminal size.
type termsize struct {
	Row, Col, Xpixel, Ypixel uint16
}

// Obtains the current width of the terminal.
func getTerminalWidth() uint {
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

// Displays error messages
func errorChecker(err *error, message string) {
	if *err != nil {
		fmt.Print(message)
		os.Exit(0)
	}
}

// If there is no argument, set the directory path to the current working directory
func getPath() string {
	if flag.NArg() < 1 {
		path, err := os.Getwd()
		errorChecker(&err, "ls: Could not obtain the current working directory.\n")
		return path
	} else {
		return flag.Arg(0)
	}
}

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

// Scans the directory and returns a list of the contents. If the directory
// does not exist, an error is printed and the program exits.
func scanDirectory() []os.FileInfo {
	directory, err := ioutil.ReadDir(getPath())
	errorChecker(&err, "ls: "+getPath()+" - No such file or directory.\n")
	return directory
}

// Checks if the file can be shown
func isHidden(file string) bool {
	switch {
	case *showHidden: // If the hidden flag is enabled
		return false
	case strings.HasPrefix(file, "."): // If it is disabled and the file is hidden
		return true
	default: // If it is disabled and the file is not hidden
		return false
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

// Determines the longest file length.
func countMaxCharacterLength(currentFileLength int) {
	if currentFileLength > maxCharLength {
		maxCharLength = currentFileLength
	}
}

// Checks if the date of the file is from a prior year, and if so print the year, else print
// only the hour and minute.
func dateFormatCheck(fileModTime time.Time) string {
	if fileModTime.Year() != time.Now().Year() {
		return fileModTime.Format(DATE_YEAR_FORMAT)
	} else {
		return fileModTime.Format(DATE_FORMAT)
	}
}

// Returns the group ID that owns the file.
func getGID(file os.FileInfo) string {
	return fmt.Sprintf("%d", file.Sys().(*syscall.Stat_t).Gid)
}

// Returns the user ID that owns the file.
func getUID(file os.FileInfo) string {
	return fmt.Sprintf("%d", file.Sys().(*syscall.Stat_t).Uid)
}

// Converts uid/gid into names.
func idToName(id string) string {
	name, err := user.LookupId(id)
	if err == nil {
		return name.Username // Returns the name represented as an actual name.
	} else {
		return id // Returns the name represented in numerical format.
	}
}

// Determines the maximum id name length for printing with long mode.
func countIDLength(uid, gid string) {
	if len(uid) > maxIDLength {
		maxIDLength = len(uid)
	}
	if len(gid) > maxIDLength {
		maxIDLength = len(gid)
	}
}

// Determines the maximum size name length for printing with long mode.
func countSizeLength(fileSize int64) {
	length := len(fmt.Sprintf("%d", fileSize))
	if length > maxSizeLength {
		maxSizeLength = length
	}
}

// If longMode is enabled, we need to grab extra statistics.
func longModeCheck(fileMode, fileModTime string, fileUID, fileGID string, fileSize int64,
	fileModeList *[]string, fileSizeList *[]int64, fileModDateList *[]string,
	fileUserList, fileGroupList *[]string) {
	if *longMode {
		*fileUserList = append(*fileUserList, fileUID)
		*fileGroupList = append(*fileGroupList, fileGID)
		*fileModeList = append(*fileModeList, fileMode)
		*fileModDateList = append(*fileModDateList, fileModTime)
		*fileSizeList = append(*fileSizeList, fileSize)
		countIDLength(idToName(fileUID), idToName(fileGID))
		countSizeLength(fileSize)
	}
}

/* The longModePrinter function prints files and their statistics one line at a time.
 * This printer will create an initial printf layout based on the maxIDLength and maxSizeLength,
 * then print the total number of files, and finally loop through each file in the list and print
 * them based on the layout. */
func longModePrinter(fileNameList, fileModeList []string, fileSizeList []int64, fileDate []string,
	fileUserList, fileGroupList []string) {
	printingLayout := "%11s %-" + fmt.Sprintf("%d", maxIDLength) + "s %-" + fmt.Sprintf("%d", maxIDLength) + "s %" + fmt.Sprintf("%d", maxSizeLength) + "d %12s %s\n"
	fmt.Println("total:", len(fileNameList))
	for index, _ := range fileNameList {
		fmt.Printf(printingLayout, fileModeList[index], idToName(fileUserList[index]),
			idToName(fileGroupList[index]), fileSizeList[index], fileDate[index],
			fileNameList[index]+RESET)
	}
}

// The countRows function counts the number of rows that will be printed. The number is
// determined by dividing the number of files by the maximum number of columns. However,
// if there is a remaindera of files left over for an incomplete row, we add an additional
// row.
func countRows(lastRowCount, maxColumns, numOfFiles *int) int {
	if *lastRowCount == 0 {
		return *numOfFiles / *maxColumns
	}
	return *numOfFiles / *maxColumns + 1 // Add additional row if the last row is incomplete.
}

// The spacer function will add spaces to the end of each file name so that they line up
// correctly when printing.
func spacer(name *string, charLength int) string {
	return string(*name + strings.Repeat(" ", maxCharLength-charLength+SPACING))
}

// Finally, the printTopToBottom function will take the fileList slice and print all the
// files from top to bottom.
func printTopToBottom(fileList []string) {
	numOfFiles := len(fileList)                                     // Total number of files
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
	 * the last row is going to kill our algorithm. To counteract that, we need to switch it
	 * to counting in intervals of two after the fifth file has been printed on each row.
	 *
	 * Once we get to the last row, we can simply print everyting in an interval based on the
	 * number of rows until we have printed the last file. */
	for printing { // Keep printing as long as the printer is turned on.
		if currentRow < numOfRows { // Prints all rows except the last row.
			for column := 1; column < maxColumns; column++ { // Prints all but the last column.
				fmt.Print(fileList[currentIndex]) // Print the file.
				if column >= lastRowCount+1 { // This is where we see if it's time to switch our index
					currentIndex += numOfRows - 1
				} else {
					currentIndex += numOfRows
				}
			}
			fmt.Println(fileList[currentIndex]) // Prints the final column in a row.
			currentRow++                                            // It's time to start printing the next row.
			currentIndex = currentRow - 1                           // We need to reset this for the next row.
		} else { // Prints the final row.
			for index := 1; index <= lastRowCount; index++ { // This is the final print run.
				fmt.Print(fileList[currentIndex]) // Print the file.
				currentIndex += numOfRows                          // Switch the index to the next file.
			}
			printing = false // We are finished printing -- turn the printing press off.
		}
	}
}

// Determines if we can print on one line.
func printOneLineCheck(currentFileLength int) {
	if totalCharLength <= terminalWidth {
		totalCharLength += currentFileLength + 2 // The additional 2 is for spacing.
	} else {
		printOneLine = false
	}
}

// Prints all files on one line if we can.
func oneLinePrinter(fileNameList *[]string) {
	for _, file := range *fileNameList {
		fmt.Print(file, "  ") // Print the file plus additional spacing
	}
	fmt.Println(RESET) // Print an additional line and reset the color.
}

// Prints all files on one column.
func singleColumnPrinter(fileNameList *[]string) {
	for _, file := range *fileNameList {
		fmt.Println(file)
		fmt.Print(RESET)
	}
}

// This switch will determine how we should print.
func printSwitch(fileNameList, fileModeList, fileModDateList []string, fileSizeList []int64,
	fileUserList, fileGroupList []string, fileLengthList []int) {
	switch {
	case *longMode: // If longMode is enabled
		longModePrinter(fileNameList, fileModeList, fileSizeList, fileModDateList,
			fileUserList, fileGroupList)
	case printOneLine: // If we can print everything on one row, do it.
		oneLinePrinter(&fileNameList)
	case *singleColumn: // If the -1 flag is set, print everything in a single column.
		singleColumnPrinter(&fileNameList)
	default: // Properly spaces our files and then prints from top to bottom.
		spaced := make([]string, 0)
		for index, file := range fileNameList { // Preprocesses the file list for printing by adding spaces.
			spaced = append(spaced, spacer(&file, fileLengthList[index]))
		}
		printTopToBottom(spaced) // Let's get printing!
		fmt.Println(RESET)             // Reset terminal at the end.
	}
}

func getModeType(file os.FileInfo) string {
	return file.Mode().String()
}

func main() {
	flag.Parse()
	processFlags()

	fileNameList := make([]string, 0)     // Stores a list of files with their name.
	fileModeTypeList := make([]string, 0) // Stores a list of mode types
	fileModeList := make([]string, 0)     // Stores a list of file mode bits.
	fileUserList := make([]string, 0)     // Stores a list of user ownerships.
	fileGroupList := make([]string, 0)    // Stores a list of group ownerships.
	fileModDateList := make([]string, 0)  // Stores a list of file modication times.
	fileSizeList := make([]int64, 0)      // Stores a list of file sizes.
	fileLengthList := make([]int, 0)  // Stores a list of character lengths per file.

	/* This for loop will loop through each file in the directory in ascending order
	 * and check if the file is to be displayed. If it is to be displayed, it will colorize
	 * the file name and obtain a lot of additional statistics for use in different printing
	 * modes in this program. */
	for _, file := range scanDirectory() {
		if isHidden(file.Name()) == false {
			fileModeTypeList = append(fileModeTypeList, getModeType(file))
			colorized := colorizer(file.Name(),
				file.IsDir(),
				file.Mode()&SYMLINK == SYMLINK,
				file.Mode()&EXECUTABLE == EXECUTABLE)
			fileNameList = append(fileNameList, colorized)
			currentFileLength := len(file.Name())
			fileLengthList = append(fileLengthList, currentFileLength)
			countMaxCharacterLength(currentFileLength)
			printOneLineCheck(currentFileLength)
			longModeCheck(file.Mode().String(), dateFormatCheck(file.ModTime()), getUID(file), getGID(file), file.Size(),
				&fileModeList, &fileSizeList, &fileModDateList, &fileUserList, &fileGroupList)
		}
	}

	// Send the recently obtained file statistics to the printSwitch function to determine how to print the files.
	printSwitch(fileNameList, fileModeList, fileModDateList, fileSizeList, fileUserList, fileGroupList, fileLengthList)
}
