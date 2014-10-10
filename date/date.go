//
// date.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Michael Murphy
//
package main

import "flag"
import "fmt"
import "os"
import "time"

const (
	RFC3339_DATE    = "2006-01-02"
	RFC3339_SECONDS = "2006-01-02 03:04:05-07:00"
	RFC3339_NS      = "2006-01-02 03:04:05.999999999Z07:00"
	ISO8601_HOUR    = "2006-01-02T15Z0700"
	ISO8601_MINUTES = "2006-01-02T15:04Z0700"
	ISO8601_SECONDS = "2006-01-02T15:04:05Z0700"
	ISO8601_NS      = "2006-01-02T15:04:05.999999999Z0700"

	HELP_TEXT = `
	Usage: date [OPTION]... [+FORMAT]

	Display the current time in the given FORMAT.

	-I [TIMESPEC], -iso-8601=[TIMESPEC]
	        output date/time in ISO 8601 format. TIMESPEC='date' for date
	        only, 'hours', 'minutes', 'seconds', or 'ns' for date and
	        time to the indicated precision.

	-r, -reference=FILE
	      display the last modification time of FILE

	-R, -rfc-1123
	      output date and time in RFC 1123 format.
	      Example: Thu, 19 Jun 2014 03:53:45 -0500

	-rfc-3339=[TIMESPEC]
	      output date and time in RFC 3339 format.  TIMESPEC='date',
              'seconds', or 'ns' for date and time to the indicated precision.
              Date and time components are separated by a single space.
              Example: 2014-06-19 03:55:49-05:00

        -u, -utc, -universal
              print Coordinated Universal Time (UTC)

        -help display this help and exit

        -version output version information and exit
`
	VERSION_TEXT = `
	       date (go-coreutils) 0.1

         Copyright (C) 2014, The GO-Coreutils Developers.
         This program comes with ABSOLUTELY NO WARRANTY; for details see
         LICENSE. This is free software, and you are welcome to redistribute
         it under certain conditions in LICENSE.
`
)

var (
	printUTC          = flag.Bool("u", false, "print Coordinated Universal Time (UTC)")
	printUTCLong      = flag.Bool("utc", false, "print Coordinated Universal Time (UTC)")
	printUTCLonger    = flag.Bool("universal", false, "print Coordinated Universal Time (UTC)")
	referenceMode     = flag.Bool("r", false, "display the last modification time of a file")
	referenceModeLong = flag.Bool("reference", false, "display the last modification time of a file")
	printISO8601      = flag.String("I", "", "output date and time in ISO 8601 format: [date|hours|minutes|seconds]")
	printISO8601Long  = flag.String("iso-8601", "", "output date and time in ISO 8601 format: [date|hours|minutes|seconds]")
	printRFC1123      = flag.Bool("R", false, "output date and time in RFC 2822 format.")
	printRFC1123Long  = flag.Bool("rfc-1123", false, "output date and time in RFC 2822 format.")
	printRFC3339      = flag.String("rfc-3339", "", "output date and time in RFC 3339 format: [date|seconds|ns]")
	help              = flag.Bool("help", false, "display help information")
	version           = flag.Bool("version", false, "output version information")
)

// getTime returns the current time in either the default time zone or UTC.
func getTime() time.Time {
	if *printUTC {
		return time.Now().UTC()
	} else {
		return time.Now()
	}
}

// getModificationTime returns the modification time of the file.
func getModificationTime(file os.FileInfo) time.Time {
	if *printUTC {
		return file.ModTime().UTC()
	} else {
		return file.ModTime()
	}
}

// printDate prints the time based on the layout format.
func printDate(t time.Time) {
	switch {
	case *printRFC1123:
		fmt.Println(t.Format(time.RFC1123Z))
	case *printRFC3339 != "" && *printRFC3339 != "date" &&
		*printRFC3339 != "seconds" && *printRFC3339 != "ns":
		fmt.Printf("date: invalid argument '%s' for '--rfc-3339'\n"+
			"Valid arguments are:\n  - 'date'\n  - 'seconds'\n  - 'ns'\n"+
			"Try 'date --help' for more information.\n", *printRFC3339)
	case *printRFC3339 == "date":
		fmt.Println(t.Format(RFC3339_DATE))
	case *printRFC3339 == "seconds":
		fmt.Println(t.Format(RFC3339_SECONDS))
	case *printRFC3339 == "ns":
		fmt.Println(t.Format(RFC3339_NS))
	case *printISO8601 == "date":
		fmt.Println(t.Format(RFC3339_DATE))
	case *printISO8601 == "hours":
		fmt.Println(t.Format(ISO8601_HOUR))
	case *printISO8601 == "minutes":
		fmt.Println(t.Format(ISO8601_MINUTES))
	case *printISO8601 == "seconds":
		fmt.Println(t.Format(ISO8601_SECONDS))
	case *printISO8601 == "ns":
		fmt.Println(t.Format(ISO8601_NS))
	default:
		fmt.Println(t.Format(time.UnixDate))
	}
}

// getReference creates an os.FileInfo of the reference file and returns it.
func getReference() os.FileInfo {
	file, err := os.Stat(flag.Arg(0))
	if err != nil {
		fmt.Printf("date: %s - No such file or directory\n", flag.Arg(0))
		os.Exit(0)
	}
	return file
}

func main() {
	switch {
	case *referenceMode && flag.NArg() < 1:
		fmt.Println("date: option requires an argument -- 'r'")
	case *referenceMode:
		printDate(getModificationTime(getReference()))
	default:
		printDate(getTime())
	}
}

func init() {
	flag.Parse()
	if *help {
		fmt.Println(HELP_TEXT)
		os.Exit(0)
	}
	if *version {
		fmt.Println(VERSION_TEXT)
		os.Exit(0)
	}
	if *printUTCLong || *printUTCLonger {
		*printUTC = true
	}
	if *printISO8601Long != "" {
		*printISO8601 = *printISO8601Long
	}
	if *printRFC1123Long {
		*printRFC1123 = true
	}
	if *referenceModeLong {
		*referenceMode = true
	}
}
