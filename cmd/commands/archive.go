package commands

import (
	"errors"
	"flag"
	"os/exec"
	"strings"

	"github.com/phuslu/log"
	"github.com/toudi/joker/internal/joker"
)

type _archiveFlags struct {
	destinationFile string
}

var archiveFlags _archiveFlags
var (
	errSpecifyDestinationFile = errors.New("specify destination file")
)

var archive = Command{
	Name: "archive",
	Help: "create archive from your data_dir so that you can move it to another machine",
	Run: func(jkr *joker.Joker, args []string) error {
		flagSet := flag.NewFlagSet("archive", flag.ContinueOnError)
		flagSet.StringVar(&archiveFlags.destinationFile, "o", "", "output file")

		if err := flagSet.Parse(args); err != nil {
			return err
		}

		if archiveFlags.destinationFile == "" {
			return errSpecifyDestinationFile
		}

		dataDir, err := jkr.GetDataDir()
		if err != nil {
			return err
		}

		command := exec.Command(
			"tar",
			"-czf",
			archiveFlags.destinationFile,
			"-C",
			strings.TrimRight(dataDir, "/"),
			".",
		)
		log.Debug().Str("cmd", command.String()).Msg("execute")
		return command.Run()
	},
}

func init() {
	registerCommand(archive)
}
