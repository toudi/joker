package commands

import "github.com/toudi/joker/internal/joker"

var up = Command{
	Name: "up",
	Help: "instantiate your project",
	Run: func(jkr *joker.Joker, args []string) error {
		return jkr.Up()
	},
}

func init() {
	registerCommand(up)
}
