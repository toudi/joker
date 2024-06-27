package commands

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/toudi/joker/internal/joker"
)

type Command struct {
	Name string
	Help string
	Run  func(jkr *joker.Joker, args []string) error
}

var AvailableCommands []Command

func registerCommand(command Command) {
	AvailableCommands = append(AvailableCommands, command)
}

func Run(name string, jkr *joker.Joker, args []string) {
	cmdIdx := slices.IndexFunc(AvailableCommands, func(cmd Command) bool {
		return cmd.Name == name
	})

	if cmdIdx < 0 {
		fmt.Fprintf(os.Stderr, "command \"%s\" not found\n\n", name)
		flag.Usage()
		os.Exit(1)
	}

	if err := AvailableCommands[cmdIdx].Run(jkr, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
}
