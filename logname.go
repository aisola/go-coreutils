package main

import "flag"
import "fmt"
import "os"
import "os/user"

const help_text = "Usage: logname \nPrint the name of the current user"

func GetCurrentUser(username *string) error {

	current_user, err := user.Current()

	*username = current_user.Username

	if err != nil {
		return err
	}

	return nil
}

func main() {
	help := flag.Bool("help", false, help_text)
	flag.Parse()

	if *help {
		fmt.Println(help)
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		fmt.Println(help)
		os.Exit(0)
	}

	if flag.NArg() == 0 && flag.NFlag() == 0 {

		var username string

		err := GetCurrentUser(&username)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		fmt.Println(username)

	}

}
