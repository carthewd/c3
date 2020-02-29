package main

import (
	"fmt"
	"os"

	"github.com/carthewd/c3/command"
)

func main() {
	if cmd, err := command.RootCmd.ExecuteC(); err != nil {
		fmt.Println(cmd.UsageString())
		os.Exit(1)
	}
}
