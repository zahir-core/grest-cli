package main

import (
	"os"

	"grest.dev/cmd"
)

var Version = "0.0.1"

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			cmd.PrintVersion()
			return
		} else if len(os.Args) > 2 {
			if os.Args[1] == "new" {
				if os.Args[2] == "app" {
					cmd.NewApp()
					return
				}
				if os.Args[2] == "service" {
					cmd.NewService()
					return
				}
			}
		}
	}
	cmd.UnknownCommand()
}
