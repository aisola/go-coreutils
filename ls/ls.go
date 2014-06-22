//
// ls.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy, Abram C. Isola
//
package main

import "bytes"
import "fmt"
import "io"
import "io/ioutil"
import "os"
import "strings"
import "flag"
import "unsafe"
import "runtime"
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
	maxIDLength     = 0                       // Statistics for the longest id name length.
	maxSizeLength   = 0                       // Statistics for the longest file size length.
	totalCharLength = 0                       // Statistics for the total number of characters.
	maxCharLength   = 0                       // Statistics for maximum file name length.
	fileList        = make([]os.FileInfo, 0)  // A list of all files being processed
	fileLengthList  = make([]int, 0)          // A list of file character lengths
	fileModeList    = make([]string, 0)       // A list of file mode strings
	fileUserList    = make([]string, 0)       // A list of user values
	fileGroupList   = make([]string, 0)       // A list of group values
	fileModDateList = make([]string, 0)       // A list of file modication times.
	fileSizeList    = make([]int64, 0)        // A list of file sizes.
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
		if strings.HasPrefix(flag.Arg(0), ".") {
			return flag.Arg(0)
		} else {
			return flag.Arg(0) + "/"
		}

	}
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

// Scans the directory and returns a list of the contents. If the directory
// does not exist, an error is printed and the program exits.
func scanDirectory() {
	directory, err := ioutil.ReadDir(getPath())
	errorChecker(&err, "ls: "+getPath()+" - No such file or directory.\n")

	if *showHidden {
		fileList = directory
	} else {
		for _, file := range directory {
			if isHidden(file.Name()) == false {
				fileList = append(fileList, file)
			}
		}
	}
}

// Obtain file statistics
func openSymlink(file string) os.FileInfo {
	var fi os.FileInfo
	if !strings.HasPrefix(file, "/") {
		fi, _ = os.Stat(getPath() + file)
	} else {
		fi, _ = os.Stat(file)
	}
	return fi
}

// Resolve the symbolic links
func readLink(file string) string {
	sympath, err := os.Readlink(getPath() + file)
	if err == nil {
		return sympath
	} else {
		return "broken link"
	}
}

/* If the file is a symbolic link, print it in cyan; a directory, blue; an executable file,
 * green; else print the file in white. */
func colorizer(file os.FileInfo) string {
	switch {
	case file.Mode()&SYMLINK != 0:
		return CYAN_SYMLINK + file.Name()
	case file.IsDir():
		return BLUE_DIR + file.Name()
	case file.Mode()&EXECUTABLE != 0:
		return GREEN_EXECUTABLE + file.Name()
	default:
		return RESET + file.Name()
	}
}

//NOTE: Below are functions for obtaining statistics
// Checks if the date of the file is from a prior year, and if so print the year, else print
// only the hour and minute.
func dateFormatCheck(fileModTime time.Time) string {
	if fileModTime.Year() != time.Now().Year() {
		return fileModTime.Format(DATE_YEAR_FORMAT)
	} else {
		return fileModTime.Format(DATE_FORMAT)
	}
}

// Opens the passwd file and returns a buffer of it's contents.
func bufferUsers() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)

	cached, _ := os.Open("/etc/passwd")
	if err != nil {
		fmt.Println("Error: passwd file does not exist.")
		os.Exit(0)
	}
	io.Copy(buffer, cached)
	return buffer
}

// Opens the group file and returns a buffer of it's contents.
func bufferGroups() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)

	cached, _ := os.Open("/etc/group")
	if err != nil {
		fmt.Println("Error: group file does not exist.")
		os.Exit(0)
	}
	
	io.Copy(buffer, cached)
	return buffer
}

// Converts a bytes buffer into a newline-separated string array.
func bufferToStringArray(buffer *bytes.Buffer) []string {
	return strings.Split(buffer.String(), "\n")
}

// Returns a colon separated string array for use in parsing /etc/group and /etc/user
func parseLine(line string) []string {
	return strings.Split(line, ":")
}

// Returns user id
func getUID(file os.FileInfo) string {
	return fmt.Sprintf("%d", file.Sys().(*syscall.Stat_t).Uid)
}

// Returns group id
func getGID(file os.FileInfo) string {
	return fmt.Sprintf("%d", file.Sys().(*syscall.Stat_t).Gid)
}

// Obtains a list of formatted file modification dates.
func getModDateList(done chan bool) {
	for _, file := range fileList {
		fileModDateList = append(fileModDateList, dateFormatCheck(file.ModTime()))
	}
	done <- true
}

// Returns the username associated to a user ID
func lookupUserID(uid string, userStringArray []string) string {
	for _, line := range userStringArray {
		values := parseLine(line)
		if len(values) > 2 {
			if values[2] == uid {
				return values[0]
			}
		}

	}
	return uid
}

// Returns the groupname associated to a group ID
func lookupGroupID(gid string, groupStringArray []string) string {
	for _, line := range groupStringArray {
		values := parseLine(line)
		if len(values) > 2 {
			if values[2] == gid {
				return values[0]
			}
		}

	}
	return gid
}

// Obtains a list of file sizes.
func getFileSize(done chan bool) {
	for _, file := range fileList {
		fileSizeList = append(fileSizeList, file.Size())
	}
	done <- true
}

// Obtains a list of user names
func getUserList(done chan bool) {
	userBuffer := bufferToStringArray(bufferUsers())
	
	for _, file := range fileList {
		uid := lookupUserID(getUID(file), userBuffer)

		fileUserList = append(fileUserList, uid)
	}
	done <- true
}

// Obtains a list of group names
func getGroupList(done chan bool) {
	groupBuffer := bufferToStringArray(bufferGroups())
	
	for _, file := range fileList {
		gid := lookupGroupID(getGID(file), groupBuffer)
		
		fileGroupList = append(fileGroupList, gid)
	}
	done <- true
}

// Obtains a list of file character lengths.
func getFileLengthList(done chan bool) {
	for _, file := range fileList {
		fileLengthList = append(fileLengthList, len(file.Name()))
	}
	done <- true
}

// Obtains the mode type of the file in string format.
func getModeType(file os.FileInfo) string {
	return file.Mode().String()
}

// Obtains a list of mode types in string format.
func getModeTypeList(done chan bool) {
	for _, file := range fileList {
		fileModeList = append(fileModeList, file.Mode().String())
	}
	done <- true
}

// Obtains a list of colorized and spaced names for printTopToBottom.
func getColorizedList() []string {
	colorizedList := make([]string, 0)
	for index, file := range fileList { // Preprocesses the file list for printing by adding spaces.
		colorizedList = append(colorizedList, spacer(colorizer(file), fileLengthList[index]))
	}
	return colorizedList
}

// Determines the character length of the longest file name.
func getMaxCharacterLength(done chan bool) {
	for _, file := range fileList {
		if len(file.Name()) > maxCharLength {
			maxCharLength = len(file.Name())
		}
	}
	done <- true
}

// Determines the max character length of file size and user/group names/ids.
func countMaxSizeLength(done chan bool) {
	for index, file := range fileList {
		countSizeLength(file.Size())
		countIDLength(fileUserList[index], fileGroupList[index])
	}
	done <- true
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

// Determines if we can print on one line.
func printOneLineCheck(done chan bool) {
	for _, file := range fileList {
		if totalCharLength <= terminalWidth {
			totalCharLength += len(file.Name()) + 2 // The additional 2 is for spacing.
		} else {
			printOneLine = false
			done <- true
		}
	}
	printOneLine = true
	done <- true
}

// NOTE: The printing-related functions are below.
func longModePrinter() {
	// Print number of files in the directory
	fmt.Println("total:", len(fileList))

	ownershipLayout := fmt.Sprintf("%d", maxIDLength)
	sizeLayout := fmt.Sprintf("%d", maxSizeLength)
	printingLayout := "%11s %-" + ownershipLayout + "s %-" + ownershipLayout + "s %" + sizeLayout + "d %12s %s\n"
	var fileName string
	for index, file := range fileList {
		if file.Mode()&SYMLINK != 0 {
			symPath := readLink(file.Name())
			fileName = colorizer(file) + RESET + " -> " + colorizer(openSymlink(symPath))
		} else {
			fileName = colorizer(file)
		}

		fmt.Printf(printingLayout, fileModeList[index], fileUserList[index],
			fileGroupList[index], fileSizeList[index], fileModDateList[index], fileName+RESET)
	}
}

// Prints all files on one line if we can.
func oneLinePrinter() {
	for _, file := range fileList {
		fmt.Print(colorizer(file), "  ") // Print the file plus additional spacing
	}
	fmt.Println(RESET) // Print an additional line and reset the color.
}

// Prints all files on one column.
func singleColumnPrinter() {
	for _, file := range fileList {
		fmt.Println(colorizer(file))
		fmt.Print(RESET)
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
// correctly when printing in the printTopToBottom function.
func spacer(name string, charLength int) string {
	return string(name + strings.Repeat(" ", maxCharLength-charLength+SPACING))
}

// Increases index count in printTopToBottom based on current position.
// The index must take into account the fact that the last row needs files to print as well.
// After we are certain that the last row is happy, we can then start increasing index count by
// the number of rows minus one.
func indexCounter(currentIndex, column, lastRowCount, numOfRows *int) int {
	if *column >= *lastRowCount+1 {
		return *currentIndex + *numOfRows - 1
	} else {
		return *currentIndex + *numOfRows
	}
}

/* printTopToBottom takes a colorized list of files for input and gathers additional statistics required
 * for printing files from top to bottom with precision. To do that, we need to know how many files are
 * to be printed, how wide the terminal is, the maximum number of columns, how many files are on the last
 * row, and how many rows are to be printed. The first portion of this function will grab all of these
 * statistics.
 * 
 * After gaining those aforementioned statistics, it is necessary to develop an algorithm for processing
 * all of this data in a manner that allows us to print each file correctly. The for printing loop contains
 * that algorithm. */
func printTopToBottom(colorizedList []string) {
	numOfFiles := len(fileList)
	maxColumns := terminalWidth / (maxCharLength + SPACING)
	lastRowCount := numOfFiles % maxColumns
	numOfRows := countRows(&lastRowCount, &maxColumns, &numOfFiles)
	printing := true
	var currentRow, currentIndex int = 1, 0

	for printing {
		if currentRow < numOfRows { // Prints all but the last row.
			for column := 1; column < maxColumns; column++ { // Prints all but the last column.
				fmt.Print(colorizedList[currentIndex]) // Print the file.
				currentIndex = indexCounter(&currentIndex, &column, &lastRowCount, &numOfRows)
			}
			fmt.Println(colorizedList[currentIndex]) // Prints the final column in a row.
			currentRow++                             // It's time to start printing the next row.
			currentIndex = currentRow - 1            // We need to reset this for the next row.
		} else { // Prints the final row.
			for index := 1; index <= lastRowCount; index++ { // This is the final print run.
				fmt.Print(colorizedList[currentIndex]) // Print the file.
				currentIndex += numOfRows              // Switch the index to the next file.
			}
			printing = false // We are finished printing -- turn the printing press off.
		}
	}
}

// This switch will determine how we should print. The available modes for printing are long mode,
// which prints files one column at a time with statistics; single column mode, which prints all
// files on one column with any statistics; single line mode, for when no mode is set and the files
// can be printed on one line; and the default mode, which prints all files from top to bottom.
func printSwitch() {
	switch {
	case *longMode:
		longModePrinter()
	case *singleColumn:
		singleColumnPrinter()
	case printOneLine:
		oneLinePrinter()
	default:
		printTopToBottom(getColorizedList())
		fmt.Println(RESET) // Reset terminal at the end.
	}
}

// NOTE: The main function
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()+1)
	flag.Parse()
	processFlags()

	// Load the directory list
	scanDirectory()

	// Channels for the goroutines to check when they finish.
	lengthDone := make(chan bool)
	oneLineCheck := make(chan bool)
	maxCharLengthCheck := make(chan bool)
	
	// The goroutines used to grab all file statistics in parallel for a slight performance boost.
	go getFileLengthList(lengthDone)
	go getMaxCharacterLength(maxCharLengthCheck)
	go printOneLineCheck(oneLineCheck)

	// If longMode is enabled
	if *longMode {
		modeDone := make(chan bool)
		modDateDone := make(chan bool)
		sizeDone := make(chan bool)
		userDone := make(chan bool)
		groupDone := make(chan bool)
		countDone := make(chan bool)
		
		go getModeTypeList(modeDone)
		go getModDateList(modDateDone)
		go getFileSize(sizeDone)
		go getUserList(userDone)
		go getGroupList(groupDone)
		
		<-userDone
		<-groupDone
		<-sizeDone
		go countMaxSizeLength(countDone)
		<-modeDone
		<-modDateDone
		<-countDone
		fmt.Println(maxIDLength)
	}
	
	// Synchronize goroutines with main
	<-lengthDone
	<-maxCharLengthCheck
	<-oneLineCheck

	// Now that statistics have been gathered, it's time to process and print them.
	printSwitch()
}
