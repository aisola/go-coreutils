//
// mkdir.go (go-coreutils) 0.1
// Copyright (C) 2014, The GO-Coreutils Developers.
//
// Written By: Corey Prak
//
package main

import (
  "flag"
  "fmt"
  "log"
  "os"
)

const (
  help_text string = `
    Usage: mkdir OPTION(S) DIRECTORY

        --help        display this help and exit
        --version     output version information and exit
        --parents     create parent directory/directories as 
                      needed, do nothing if already existing
  `

  version_text = `
    mkdir (go-coreutils) 0.1

    Copyright (C) 2014, The GO-Coreutils Developers.
    This program comes with ABSOLUTELY NO WARRANTY; for details see
    LICENSE. This is free software, and you are welcome to redistribute 
    it under certain conditions in LICENSE.
  `

  parents_text = `
    Create parent directory/directories as needed for a path. If any
    already exist, do nothing. 
  `

  usage_text = `usage: mkdir [-parents] directory ...`
)

var (
  help    = flag.Bool("help", false, help_text)
  version = flag.Bool("version", false, version_text)
  parents = flag.Bool("parents", false, parents_text)
)

func main() {
  flag.Parse()

  if flag.NArg() == 0 {
    fmt.Println(usage_text)
    os.Exit(0)
  }

  if *help {
    fmt.Println(help_text)
    os.Exit(0)
  }

  if *version {
    fmt.Println(version_text)
    os.Exit(0)
  }

  for i := 0; i < flag.NArg(); i++ {
    if *parents {
      mkdirAllError := os.MkdirAll(flag.Arg(i), os.ModePerm)

      if mkdirAllError != nil {
        log.Fatalln(mkdirAllError)
      }
    } else {
      mkdirError := os.Mkdir(flag.Arg(i), os.ModePerm)

      if mkdirError != nil {
        log.Fatalln(mkdirError)
      }
    }
  }
}
