//
// ls.go (go-coreutils) 0.1
// Copyright (C) 2014, Abram C. Isola.
//
// MUCH OF THIS CODE WAS PROVIDED BY:
//     Michael Murphy (mmstickman@gmail.com)
//         https://gist.github.com/mmstick/2046b2a8ec2c5c381ad9
//
package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "strings"
  "flag"
  "unsafe"
  "syscall"
)

const (
  TERMINAL_INFO = 0x5413 //
  EXECUTABLE = 0111 // File executable bit
  CYAN_SYMLINK = "\x1b[36;1m"
  BLUE_DIR = "\x1b[34;1m"
  GREEN_EXECUTABLE = "\x1b[32;1m"
  RESET = "\x1b[0m"
  SPACING = 1
)

type winsize struct {
  Row, Col, Xpixel, Ypixel uint16
}

func getTerminalWidth() uint {
  ws := &winsize{}
  retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
    uintptr(syscall.Stdin),
    uintptr(TERMINAL_INFO),
    uintptr(unsafe.Pointer(ws)))
  if int(retCode) == -1 {
    panic(errno)
  }
  return uint(ws.Col)
}

func main() {
  
  // A list of variables used in this function
  
  var colorized string
  var isHidden, printOneLine bool = false, true
  var totalCharLength, maxCharLength int
  fileNameList := make([]string, 0)
  fileLengthList := make([]int, 0)
  terminalWidth := int(getTerminalWidth())
  
  
  // A list of available flags
  
  showHidden := flag.Bool("a", false, "list hidden files and directories")
  oneColumn := flag.Bool("1", false, "list files by one column")
  flag.Parse()
  
  // If there is no argument, set the directory path to "./"
  
  var path string
  args := flag.Args()
  if len(args) < 1 { path = "./"} else { path = args[0] }
  
  /* Scans the directory and returns a list of the contents. If the directory
   * does not exist, an error is printed and the program exits. */
  
  directory, err := ioutil.ReadDir(path)
  if err != nil {
    fmt.Print("gols: ", path, " - No such file or directory\n")
    return
  }
  
  /* This for loop will loop through each file in the directory in ascending order
   * and call the filterHidden function to determine whether or not it should
   * display the file. */

  for _, file := range directory {
    filterHidden(file.Name(), showHidden, &isHidden)
      if isHidden == false {
	
	/* If it should be displayed, it will call the setModeType function to
	 * determine what kind of file the file is (Directory, Symlink, Executable,
	 * or File).
	 *
	 * After determining the modeType, we send the modeType and name of the file
	 * to the colorizer function to colorize the name based on the modeType. */
	
	colorized = colorizer(file.Name(), setModeType(file.IsDir(),
	  file.Mode() & os.ModeSymlink == os.ModeSymlink,
	  file.Mode() & EXECUTABLE == EXECUTABLE))
	
	/* Rather than print directly to the terminal, we store name and length of
	 * the name in slices for later use. */
	
	fileNameList = append(fileNameList, colorized)
	fileLengthList = append(fileLengthList, len(file.Name()))
	
	/* Finally, this is a good spot for getting some statistics on the totalCharLength
	 * and maxCharLength of the files. */
	
	if totalCharLength <= terminalWidth {
	  totalCharLength += len(file.Name()) + SPACING
	} else {
	  printOneLine = false
	}
	
	if len(file.Name()) > maxCharLength {
	  maxCharLength = len(file.Name())
	}
      }
  }
  
  /* This switch will determine whether we should print everything in one line, in one
   * column, or send it to the printTopToBottom function for printing. */
  
  switch {
    case printOneLine: // If we can print everything on one line, do it.
      for _, file := range fileNameList {
	fmt.Print(file, "  ")
      }
      fmt.Println(RESET) // Print an additional line and reset the color.
    case *oneColumn: // If the -1 flag is set, print everything in one column.
      for _, file := range fileNameList {
	fmt.Println(file)
	fmt.Print(RESET) // Same as above, except we do not need a new line here.
      }
    default: // Properly spaces our files and then prints from top to bottom.
      spaced := make([]string, 0)
      for index, file := range fileNameList {
	spaced = append(spaced, spacer(&file, &fileLengthList[index], &maxCharLength))
      }
      printTopToBottom(spaced, &maxCharLength, &terminalWidth) // Let's get printing!
      fmt.Println(RESET) // Once again.
  }
}

/* The filterHidden function will check if the hidden flag is set, and if not, whether
 * the file is a hidden file or not. If the flag is not set and the file is a hidden
 * file (starts with a period), then it will not be printed. */

func filterHidden (file string, showHidden, isHidden *bool) {
  switch {
    case *showHidden:
      *isHidden = false
    case strings.HasPrefix(file, "."):
      *isHidden = true
    default:
      *isHidden = false
  }
}

/* The setModeType function will determine whether the file is a directory, symlink,
 * executable, or file and then set the modeType accordingly. */

func setModeType (isDir, isSymlink, isExecutable bool) string {
  var modeType string
  switch {
    case isSymlink:
      modeType = "Symlink"
    case isDir:
      modeType = "Dir"
    case isExecutable:
      modeType = "Exe"
  }
  return modeType
}

/* If the file is a symbolic link, print it in cyan; a directory, blue; an executable file,
 * green; else print the file in white. */

func colorizer(file, modeType string) string {
  switch modeType {
    case "Symlink":
      return CYAN_SYMLINK + file
    case "Dir":
      return BLUE_DIR + file
    case "Exe":
      return GREEN_EXECUTABLE + file
    default:
      return RESET + file
  }
  return RESET + file
}

/* The spacer function will add spaces to the end of each file name so that they line up
 * correctly when printing. */

func spacer(name *string, charLength, maxCharLength *int) string {
  return string(*name + strings.Repeat(" ", *maxCharLength - *charLength + SPACING))
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

func printTopToBottom(spaced []string, maxCharLength, terminalWidth *int) {
  
  /* The first thing that we need to know is how many columns we will be printing, how
   * many files are left over to be printed on the last row, and the number of rows that
   * we will be printing. */
  
  numOfFiles := len(spaced)
  maxColumns := *terminalWidth / (*maxCharLength + SPACING)
  lastRowCount := numOfFiles % maxColumns
  numOfRows := countRows(&lastRowCount, &maxColumns, &numOfFiles)
  var currentRow, currentIndex int = 1, 0
  printing := true // Set the initial value 
  
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
  
  for printing {
    switch {
      case currentRow < numOfRows: // Prints all rows except the last row.
	for column := 1; column < maxColumns; column++ { // Prints all but the last column.
	  fmt.Print(spaced[currentIndex])
	  if column >= lastRowCount + 1 { // This is where we see if it's time to switch our
	    currentIndex += numOfRows - 1 // interval to the number of rows minus one.
	  } else {
	  currentIndex += numOfRows
	  }
	}
	fmt.Println(spaced[currentIndex]) // Prints the final column in a row.
	currentRow++ // It's time to start printing the next row.
	currentIndex = currentRow - 1 // We need to reset this for the next row.
      default: // Prints the final row.
	for index := 1; index <= lastRowCount; index++ {
	  fmt.Print(spaced[currentIndex])
	  currentIndex += numOfRows
	}
	printing = false // We are finished printing.
    }
  }
}