package commands

import (
	"flag"

	"github.com/toudi/joker/internal/joker"
)

type _clearFlags struct {
	force bool
}

var clearFlags _clearFlags

var clear = Command{
	Name: "clear",
	Help: "undo any changes performed by bootstrap for the selected service(s).",
	Run: func(jkr *joker.Joker, args []string) error {
		var err error

		flagSet := flag.NewFlagSet("clear", flag.ExitOnError)
		flagSet.BoolVar(
			&clearFlags.force,
			"force",
			false,
			"perform the action regardless of what's specified in .jokerstate file",
		)

		err = flagSet.Parse(args)

		if err != nil {
			flagSet.PrintDefaults()
			return err
		}

		// let's check if there's any service requested for the clear operation
		var serviceName = ""
		if flagSet.NArg() > 0 {
			serviceName = flagSet.Arg(0)
		}

		return jkr.Clear(serviceName, clearFlags.force)
	},
}

func init() {
	registerCommand(clear)
}
