package commands

import (
	"errors"

	"github.com/toudi/joker/internal/joker"
)

var errInvalidInvoke = errors.New("invalid invoke")

var run = Command{
	Name: "run",
	Help: "Run the specified command",
	Run: func(jkr *joker.Joker, args []string) error {
		if len(args) != 1 {
			return errInvalidInvoke
		}
		jkr.StreamHandler()
		return jkr.CallCommand(args[0])
	},
}

func init() {
	registerCommand(run)
}
