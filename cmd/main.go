package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/toudi/joker/cmd/commands"
	"github.com/toudi/joker/internal/joker"
)

// based on
// https://www.janekbieser.dev/posts/cli-app-with-subcommands-in-go/

type Flags struct {
	verbose   bool
	jokerfile string
	statefile string
}

var subcommands []commands.Command = commands.AvailableCommands

func main() {
	var flags Flags

	flag.BoolVar(&flags.verbose, "v", false, "verbose mode")
	flag.StringVar(&flags.jokerfile, "f", "jokerfile", "jokerfile location")
	flag.StringVar(&flags.statefile, "s", ".jokerstate", "statefile location")

	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	ctx := context.Background()

	jkr, err := joker.Joker_init(ctx, flags.jokerfile)
	if err != nil {
		fmt.Printf("[error] unable to initialize joker: %v\n", err)
		os.Exit(-1)
	}

	if err = jkr.SetStatefile(flags.statefile); err != nil {
		fmt.Printf("[error] unable to set statefile: %v\n", err)
		os.Exit(-1)
	}

	jkr.Defer(func() {
		if err := jkr.SaveState(); err != nil {
			fmt.Printf("unable to persist statefile: %v\n", err)
		}
	})

	commands.Run(flag.Arg(0), jkr, flag.Args()[1:])
}

func usage() {
	intro := `joker is a bare-matel alternative to docker-compose.
	
Usage:
	joker [flags] <command> [command flags]`
	fmt.Fprintln(os.Stderr, intro)

	fmt.Fprintln(os.Stderr, "\nCommands:")
	for _, cmd := range subcommands {
		fmt.Fprintf(os.Stderr, "  %-8s %s\n", cmd.Name, cmd.Help)
	}
	fmt.Fprintln(os.Stderr, "\nFlags:")
	// Prints a help string for each flag we defined earlier using
	// flag.BoolVar (and related functions)
	flag.PrintDefaults()

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(
		os.Stderr,
		"Run `joker <command> -h` to get help for a specific command\n\n",
	)
}
