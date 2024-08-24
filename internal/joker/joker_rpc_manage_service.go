package joker

import (
	"flag"
	"strconv"
	"strings"
	"syscall"

	"github.com/samber/lo"
	"github.com/toudi/joker/internal"
)

const rpcCmdStopService = "stop"
const rpcCmdStartService = "start"

func parseShutdownOptionsAndService(input string) (serviceShutdownOptions, string, error) {
	var shutdownOptions = serviceShutdownOptions{
		signal: syscall.SIGTERM,
		wait:   true,
	}

	parser := flag.NewFlagSet("", flag.ContinueOnError)
	parser.BoolVar(&shutdownOptions.withDependencies, "deps", false, "")
	parser.StringVar(&shutdownOptions.signalInput, "signal", "", "")

	if err := parser.Parse(strings.Split(input, " ")); err != nil {
		return shutdownOptions, "", err
	}

	if shutdownOptions.signalInput != "" {
		signalNo, err := strconv.Atoi(shutdownOptions.signalInput)
		if err != nil {
			// this is not an integer - let's revert to string parsing
			shutdownOptions.signal = internal.ParseSignalFromString(
				shutdownOptions.signalInput,
			)
		} else {
			shutdownOptions.signal = internal.ParseSignalFromInt(signalNo)
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

	service, found := lo.Find(
		j.services,
		func(s *Service) bool { return s.definition.Name == serviceName },
	)

	if !found {
		return errUnknownService
	}

	return service.Down(shutdownOptions)
}

func rpcCmdStartServiceHandler(j *Joker, serviceName string) error {
	service, found := lo.Find(
		j.services,
		func(s *Service) bool { return s.definition.Name == serviceName },
	)

	if !found {
		return errUnknownService
	}

	return service.Up(j.ctx, j)
}
