package joker

import (
	"errors"
	"strings"

	"github.com/phuslu/log"
)

var (
	errUnknownCommand     = errors.New("unknown command")
	errNoInstructions     = errors.New("command does not contain any instructions")
	errInvalidInstruction = errors.New("invalid instruction")
	errUnknownService     = errors.New("unknown service")
)

const rpcCmdCall = "call"

func rpcCmdCallHandler(j *Joker, command string) error {
	input := strings.Split(command, " ")
	if len(input) != 1 {
		return errInvalidInstruction
	}

	// first of all, check if the command exists
	cmdDefinition, exists := j.config.Commands[command]
	if !exists {
		return errUnknownCommand
	}

	// great. we know the command exists, now let's try to interpret it.
	tmpSlice, ok := cmdDefinition.([]interface{})
	if !ok {
		return errNoInstructions
	}

	log.Debug().Str("command", command).Msg("executing")

	// even better. there are instructions. Let's proceed
	for _, instruction := range tmpSlice {
		instructionString, ok := instruction.(string)
		if !ok {
			// right now we can bail. moving forward, we could maybe support
			// stuff like conditional execution
			return errInvalidInstruction
		}

		log.Debug().Msgf("command: %s", instructionString)

		var isRpc bool = false

		for prefix, handler := range availableRPCs {
			if strings.HasPrefix(instructionString, prefix) {
				isRpc = true
				if err := handler(j, strings.TrimLeft(strings.TrimPrefix(instructionString, prefix), " ")); err != nil {
					return err
				}
			}
		}

		if !isRpc {
			// it's not a known RPC therefore the only thing left is to
			// try to interpret it as a command.

			interpolated := j.interpolateEnvVars(instructionString, nil)
			log.Debug().Msgf("command after interpolation: %s\n", interpolated)
			process, err := j.prepareCommand(j.ctx, interpolated, nil)
			if err != nil {
				return err
			}
			if err = handleCommandStream(j, command, process); err != nil {
				return err
			}
			if err = process.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (j *Joker) CallCommand(commandName string) error {
	return rpcCmdCallHandler(j, commandName)
}
