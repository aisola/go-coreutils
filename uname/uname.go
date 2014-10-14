//
// uname.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy & Abram Isola
//

// +build linux

package main

import "flag"
import "fmt"
import "io/ioutil"
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

// sysinfo stores all information regarding the system in strings.
type sysinfo struct {
	name      string
	node      string
	release   string
	version   string
	machine   string
	domain    string
	os        string
	processor string
}

// utsnameToString converts the utsname to a string and returns it.
func utsnameToString(unameArray [65]int8) string {
	var byteString [65]byte
	var indexLength int
	for ; unameArray[indexLength] != 0; indexLength++ {
		byteString[indexLength] = uint8(unameArray[indexLength])
	}
	return string(byteString[:indexLength])
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

// bufferCPUInfo returns a newline-delimited string slice of '/proc/cpuinfo'.
func bufferCPUInfo() []string {
	cached, _ := ioutil.ReadFile("/proc/cpuinfo")
	return strings.Split(string(cached), "\n")
}

// getProcessorName returns the processor name.
func getProcessorName() string {
	if runtime.GOOS == "linux" {
		return bufferCPUInfo()[4][13:]
	} else {
		return "unknown"
	}
}

// getSystemInfo returns a sysinfo struct containing system information.
func getSystemInfo() *sysinfo {
	var utsname syscall.Utsname
	_ = syscall.Uname(&utsname)
	sys := sysinfo{
		name:      utsnameToString(utsname.Sysname),
		node:      utsnameToString(utsname.Nodename),
		release:   utsnameToString(utsname.Release),
		version:   utsnameToString(utsname.Version),
		machine:   utsnameToString(utsname.Machine),
		domain:    utsnameToString(utsname.Domainname),
		os:        getOS(),
		processor: getProcessorName(),
	}
	return &sys
}

/* unameString generates a string for printing based on input arguments and
 * system information gathered by 'sys'. */
func (sys *sysinfo) unameString() string {
	if flag.NFlag() == 0 {
		return sys.name
	}
	printArray := make([]string, 0)
	if *printAll {
		printArray = append(printArray,
			fmt.Sprintf("%s %s %s %s %s %s %s", sys.name, sys.node,
				sys.release, sys.version, sys.machine,
				sys.processor, sys.os))
	}
	if *printKernelname || *printKernelnameLong {
		printArray = append(printArray, sys.name)
	}
	if *printNodename || *printNodenameLong {
		printArray = append(printArray, sys.node)
	}
	if *printRelease || *printReleaseLong {
		printArray = append(printArray, sys.release)
	}
	if *printVersion || *printVersionLong {
		printArray = append(printArray, sys.version)
	}
	if *printMachine || *printMachineLong {
		printArray = append(printArray, sys.machine)
	}
	if *printDomain || *printDomainLong {
		printArray = append(printArray, sys.domain)
	}
	if *printOS || *printOSLong {
		printArray = append(printArray, sys.os)
	}
	if *printProcessor || *printProcessorLong {
		printArray = append(printArray, sys.processor)
	}
	return strings.Join(printArray, " ")
}

func main() {
	sys := getSystemInfo()
	fmt.Println(sys.unameString())
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
