//
// uname.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy & Abram Isola
//
package main

import "bytes"
import "flag"
import "fmt"
import "io"
import "os"
import "runtime"
import "strings"
import "syscall"

const (
	help_text = `
    Usage: uname [OPTION]...

    Print certain system information.  With no OPTION, same as -s.

        -help        display this help and exit
        -version     output version information and exit

        -a, all
              print all information, in the following order.

        -s, -kernel-name
              print the kernel name

        -n, -nodename
              print the network node hostname

        -r, -kernel-release
              print the kernel release

        -v, -kernel-version
              print the kernel version

        -m, -machine
              print the machine hardware name

        -o, -operating-system
              print the operating system

        -p, -processor-name
              print the processor name
    `
	version_text = `
    uname (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute
    it under certain conditions in LICENSE.
`
)

var (
	help                = flag.Bool("help", false, "print help")
	version             = flag.Bool("version", false, "print version")
	printAll            = flag.Bool("a", false, "print all information")
	printAllLong        = flag.Bool("all", false, "print all information")
	printKernelname     = flag.Bool("s", false, "print the kernel name")
	printKernelnameLong = flag.Bool("kernel-name", false, "print the kernel name")
	printNodename       = flag.Bool("n", false, "print the network node hostname")
	printNodenameLong   = flag.Bool("nodename", false, "print the network node hostname")
	printRelease        = flag.Bool("r", false, "print the kernel release")
	printReleaseLong    = flag.Bool("kernel-release", false, "print the kernel release")
	printVersion        = flag.Bool("v", false, "print the kernel version")
	printVersionLong    = flag.Bool("kernel-version", false, "print the kernel version")
	printMachine        = flag.Bool("m", false, "print the machine architecture")
	printMachineLong    = flag.Bool("machine", false, "print the machine architecture")
	printDomain         = flag.Bool("d", false, "print the domain name the machine belongs to")
	printDomainLong     = flag.Bool("domain", false, "print the domain name the machine belongs to")
	printOS             = flag.Bool("o", false, "print the operating system")
	printOSLong         = flag.Bool("operating-system", false, "print the operating system")
	printProcessor      = flag.Bool("p", false, "print the processor name")
	printProcessorLong  = flag.Bool("processor-name", false, "print the processor name")
)

// Each utsname is a 65-width int array. Therefore, to convert it into something readable,
// We must convert the int8's into a string
func utsnameToString(unameArray [65]int8) string {
	var byteString [65]byte

	var indexLength int
	for ; unameArray[indexLength] != 0; indexLength++ {
		byteString[indexLength] = uint8(unameArray[indexLength])
	}
	return string(byteString[0:indexLength])
}

// Returns the operating system name.
//TODO: Add additional operating systems.
func getOS() string {
	var osname string

	if runtime.GOOS == "linux" {
		osname = "GNU/Linux"
	}

	return osname
}

// Buffers /proc/cpuinfo into memory for parsing.
func bufferCPUInfo() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	cached, _ := os.Open("/proc/cpuinfo")
	io.Copy(buffer, cached)
	return buffer
}

// Converts a bytes buffer into a newline-separated string array.
func bufferToStringArray(buffer *bytes.Buffer) []string {
	return strings.Split(buffer.String(), "\n")
}

// Parses the cpuinfo file for CPU information
func parseCPUInfo() string {
	infoArray := bufferToStringArray(bufferCPUInfo())
	modelLine := infoArray[4] // model name is on the fifth line.
	return modelLine[13:]     // The name information is stored after the 18th char.
}

// Returns the processor name
func getProcessorName() string {
	if runtime.GOOS == "linux" {
		return parseCPUInfo()
	} else {
		return "unknown"
	}
}

func main() {
	flag.Parse()
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	// Obtain information about the system.
	var utsname syscall.Utsname
	_ = syscall.Uname(&utsname)
	sysname := utsnameToString(utsname.Sysname)
	nodename := utsnameToString(utsname.Nodename)
	release := utsnameToString(utsname.Release)
	version := utsnameToString(utsname.Version)
	machine := utsnameToString(utsname.Machine)
	domain := utsnameToString(utsname.Domainname)
	osname := getOS()
	processorname := getProcessorName()

	if flag.NFlag() == 0 {
		fmt.Println(sysname)
		os.Exit(0)
	}

	// Store printing information in an array.
	printArray := make([]string, 0)
	if *printAll {
		printArray = append(printArray,
			fmt.Sprintf("%s %s %s %s %s %s %s", sysname, nodename,
				release, version, machine, processorname, osname))
	}
	if *printKernelname || *printKernelnameLong {
		printArray = append(printArray, sysname)
	}
	if *printNodename || *printNodenameLong {
		printArray = append(printArray, nodename)
	}
	if *printRelease || *printReleaseLong {
		printArray = append(printArray, release)
	}
	if *printVersion || *printVersionLong {
		printArray = append(printArray, version)
	}
	if *printMachine || *printMachineLong {
		printArray = append(printArray, machine)
	}
	if *printDomain || *printDomainLong {
		printArray = append(printArray, domain)
	}
	if *printOS || *printOSLong {
		printArray = append(printArray, osname)
	}
	if *printProcessor || *printProcessorLong {
		printArray = append(printArray, processorname)
	}

	// Print the information.
	fmt.Println(strings.Join(printArray, " "))
}
