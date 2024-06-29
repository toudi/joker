package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/phuslu/log"
	"github.com/toudi/joker/cmd/commands"
	"github.com/toudi/joker/internal/joker"
)

// based on
// https://www.janekbieser.dev/posts/cli-app-with-subcommands-in-go/

type Flags struct {
	verbose   bool
	trace     bool
	jokerfile string
	statefile string
}

var subcommands []commands.Command = commands.AvailableCommands

func main() {
	var flags Flags

	flag.BoolVar(&flags.verbose, "v", false, "verbose mode")
	flag.BoolVar(&flags.trace, "vv", false, "very verbose mode (logging set to trace)")
	flag.StringVar(&flags.jokerfile, "f", "jokerfile", "jokerfile location")
	flag.StringVar(&flags.statefile, "s", ".jokerstate", "statefile location")

	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	ctx := context.Background()

	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Level:      log.InfoLevel,
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}

	if flags.verbose {
		log.DefaultLogger.Level = log.DebugLevel
	} else if flags.trace {
		log.DefaultLogger.Level = log.TraceLevel
	}

	jkr, err := joker.Joker_init(ctx, flags.jokerfile)
	if err != nil {
		log.Error().Err(err).Msg("unable to initialize joker")
		os.Exit(-1)
	}

	if err = jkr.SetStatefile(flags.statefile); err != nil {
		log.Error().Err(err).Msg("unable to set statefile")
		os.Exit(-1)
	}

	jkr.Defer(func() {
		if err := jkr.SaveState(); err != nil {
			log.Error().Err(err).Msg("unable to persist statefile")
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
