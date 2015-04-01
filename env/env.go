//
// env.go (go-coreutils) 0.1
// Copyright (C) 2015, The GO-Coreutils Developers.
//
// Written By: Haruki Tsurumoto
//
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	help_text = `
    Usage: env [OPTION]... [-] [NAME=VALUE]... [COMMAND [ARG]...]

    Set each NAME to VALUE in the environment and run COMMAND.

      -i, --ignore-environment  start with an empty environment
      -0, --null           end each output line with 0 byte rather than newline
      -u, --unset=NAME     remove variable from the environment
        --help     display this help and exit
        --version  output version information and exit
    `
	version_text = `
    env (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute
    it under certain conditions in LICENSE.
`
)

var (
	help          = flag.Bool("help", false, "help")
	version       = flag.Bool("version", false, "version_text")
	ignoreEnv     = flag.Bool("i", false, "start with an empty environment")
	ignoreEnvLong = flag.Bool("ignore-environment", false, "start with an empty environment")
	nullOpt       = flag.Bool("0", false, "end each output line with 0 byte rather than newline")
	nullOptLong   = flag.Bool("null", false, "end each output line with 0 byte rather than newline")
	unset         = flag.String("u", "", "remove variable from the environment")
	unsetLong     = flag.String("unset", "", "remove variable from the environment")
	environ       = os.Environ()
)

func setenv(name, value string) {
	for i := 0; i < len(environ); i++ {
		e := strings.SplitN(environ[i], "=", 2)
		if e[0] == name {
			environ[i] = name + "=" + value
			return
		}
	}
	environ = append(environ, name+"="+value)
}
func unsetenv(name string) {
	for i := 0; i < len(environ); {
		e := strings.SplitN(environ[i], "=", 2)
		if e[0] == name && len(e) == 2 {
			environ = append(environ[:i], environ[i+1:]...) // delete
		} else {
			i++
		}
	}
}

func main() {
	optNullTerminateOutput := false
	flag.Parse()
	if *help {
		fmt.Println(help_text)
		os.Exit(0)
	}
	if *version {
		fmt.Println(version_text)
		os.Exit(0)
	}
	if *ignoreEnv || *ignoreEnvLong {
		environ = make([]string, 0)
	}
	if *nullOpt || *nullOptLong {
		optNullTerminateOutput = true
	}
	if *unset != "" {
		unsetenv(*unset)
	}
	if *unsetLong != "" {
		unsetenv(*unsetLong)
	}
	arg := flag.Args()
	if len(arg) >= 1 && arg[0] == "-" {
		environ = make([]string, 0)
		arg = arg[1:]
	}
	if len(arg) >= 1 {
		for i, _ := range arg {
			if strings.Index(arg[i], "=") > 0 {
				e := strings.SplitN(arg[i], "=", 2)
				setenv(e[0], e[1])
			} else {
				// run COMMAND
				cmd := exec.Command(arg[i], arg[i+1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Env = environ
				err := cmd.Run()
				if err != nil {
					_, lookPathErr := exec.LookPath(arg[i])
					if lookPathErr != nil {
						fmt.Printf("env: %s: No such file or directory\n", arg[i])
					} else {
						fmt.Println(err)
					}
					os.Exit(1)
				}
				os.Exit(0)
			}
		}
	}
	// print all Environment
	for _, s := range environ {
		if optNullTerminateOutput {
			fmt.Print(s + string('\000'))
		} else {
			fmt.Println(s)
		}
	}
}
