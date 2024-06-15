package main

import (
	"context"
	"fmt"
	"os"

	"github.com/toudi/joker/internal/joker"
	"github.com/toudi/joker/internal/jokerfile"
)

func main() {
	var exitcode int = 0

	defer func() {
		os.Exit(exitcode)
	}()

	if len(os.Args) < 2 {
		fmt.Printf("usage: joker [up|down|clear]\n")
		exitcode = 1
		return
	}

	cmd := os.Args[1]

	const srcfile = "jokerfile"
	const statefile = ".jokerstate"

	config, err := jokerfile.Parse(srcfile)
	if err != nil {
		fmt.Printf("error opening jockerfile: %v\n", err)
		exitcode = -1
		return
	}

	jkr, _ := joker.Joker_init(config)
	if err := jkr.SetStatefile(statefile); err != nil {
		fmt.Printf("unable to set statefile: %v\n", err)
		exitcode = -1
		return
	}

	jkr.Defer(func() {
		if err := jkr.SaveState(); err != nil {
			fmt.Printf("unable to persist statefile: %v\n", err)
		}
	})

	ctx := context.Background()

	if cmd == "up" {
		if err := jkr.Up(ctx); err != nil {
			fmt.Printf("error running up: %v\n", err)
			// try to stop the services though this may not be successful.
			// but it doesn't hurt to try
			_ = jkr.Down()
			exitcode = 1
			return
		}
	} else if cmd == "down" {
		if err := jkr.Down(); err != nil {
			fmt.Printf("error running down: %v\n", err)
			exitcode = 1
			return
		}
	} else if cmd == "clear" {
		var serviceName = ""
		if len(os.Args) == 3 {
			serviceName = os.Args[2]
		}
		if err := jkr.Clear(ctx, serviceName); err != nil {
			fmt.Printf("error running clear: %v\n", err)
			exitcode = 1
			return
		}
	} else {
		fmt.Printf("unsupported command.\n")
		exitcode = 1
		return
	}
}
