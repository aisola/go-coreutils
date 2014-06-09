//
// uptime.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Abram C. Isola
//
package main

import "bufio"
import "bytes"
import "fmt"
import "flag"
import "io/ioutil"
import "os"
import "strconv"
import "strings"
import "syscall"
import "time"

const (
	help_text string = `
    Usage: uptime
    
    tell how long the system has been running

        --help        display this help and exit
        --version     output version information and exit
    `
	version_text = `
    uptime (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
`
)

type Load struct {
	L1, L5, L15 float64
}

type Uptime struct {
	Time float64
}

func (self *Load) Get() error {
	line, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil
	}

	f := strings.Fields(string(line))

	self.L1, _ = strconv.ParseFloat(f[0], 64)
	self.L5, _ = strconv.ParseFloat(f[1], 64)
	self.L15, _ = strconv.ParseFloat(f[2], 64)

	return nil
}

func (self *Uptime) Get() error {
	sysinfo := syscall.Sysinfo_t{}

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return err
	}

	self.Time = float64(sysinfo.Uptime)

	return nil
}

func (self *Uptime) Format() string {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	uptime := uint64(self.Time)

	days := uptime / (60 * 60 * 24)

	if days != 0 {
		s := ""
		if days > 1 {
			s = "s"
		}
		fmt.Fprintf(w, "%d day%s, ", days, s)
	}

	minutes := uptime / 60
	hours := minutes / 60
	hours %= 24
	minutes %= 60

	fmt.Fprintf(w, "%2d:%02d", hours, minutes)

	w.Flush()
	return buf.String()
}

func Users() int { return 0 }

func main() {
	version := flag.Bool("version", false, version_text)
	flag.Parse()

	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}

	up := Uptime{}
	up.Get()
	load := Load{}
	load.Get()

	fmt.Printf(" %s up %s load average: %.2f, %.2f, %.2f\n",
		time.Now().Format("15:04:05"),
		up.Format(),
		load.L1, load.L5, load.L15)
}
