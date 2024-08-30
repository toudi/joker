package joker

import (
	"flag"
	"strconv"
	"strings"
	"syscall"

	"github.com/toudi/joker/internal/utils"
)

const rpcCmdStopService = "stop"
const rpcCmdStartService = "start"

func parseShutdownOptionsAndService(input string) (serviceShutdownOptions, string, error) {
	var shutdownOptions = serviceShutdownOptions{
		signal: syscall.SIGTERM,
		wait:   true,
	}

	parser := flag.NewFlagSet("", flag.ContinueOnError)
	parser.BoolVar(
		&shutdownOptions.withDependencies,
		"with-deps",
		false,
		"process service dependencies",
	)
	parser.StringVar(&shutdownOptions.signalInput, "signal", "", "")

	if err := parser.Parse(strings.Split(input, " ")); err != nil {
		return shutdownOptions, "", err
	}

	if shutdownOptions.signalInput != "" {
		signalNo, err := strconv.Atoi(shutdownOptions.signalInput)
		if err != nil {
			// this is not an integer - let's revert to string parsing
			shutdownOptions.signal = utils.ParseSignalFromString(
				shutdownOptions.signalInput,
			)
		} else {
			shutdownOptions.signal = utils.ParseSignalFromInt(signalNo)
		}
	}

	serviceName := parser.Arg(0)

	return shutdownOptions, serviceName, nil
}

func rpcCmdStopServiceHandler(j *Joker, args string) error {
	shutdownOptions, serviceName, err := parseShutdownOptionsAndService(args)
	if err != nil {
		return err
	}

	return j.StopService(serviceName, shutdownOptions)
}

type serviceStartOptions struct {
	WithDependencies bool
}

func rpcCmdStartServiceHandler(j *Joker, args string) error {
	var startOptions serviceStartOptions
	parser := flag.NewFlagSet("", flag.ContinueOnError)
	parser.BoolVar(
		&startOptions.WithDependencies,
		"with-deps",
		false,
		"process service dependencies",
	)

	if err := parser.Parse(strings.Split(args, " ")); err != nil {
		return err
	}

	serviceName := parser.Arg(0)

	return j.StartService(serviceName, startOptions.WithDependencies)
}
