package main

import (
	"os"

	"github.com/carthewd/c3/command"
)

func main() {
	if _, err := command.RootCmd.ExecuteC(); err != nil {
		os.Exit(1)
	}
}
