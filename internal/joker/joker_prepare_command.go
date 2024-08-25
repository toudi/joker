package joker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/flosch/pongo2/v6"
	"github.com/samber/lo"
)

// this function takes the input, which can be either a list of strings or
// a string and prepares command that can be executed
func (j *Joker) prepareCommand(
	ctx context.Context,
	definition interface{},
	additionalContext *pongo2.Context,
) (*exec.Cmd, error) {
	// first, let's determine if the command is a single string or a list of
	// strings
	var commands []string
	var env []string

	if sstring, ok := definition.(string); ok {
		// it is a single string
		commands = []string{j.interpolateEnvVars(sstring, additionalContext)}
	} else if mstrings, ok := definition.([]interface{}); ok {
		commands = lo.Map(mstrings, func(sstring interface{}, index int) string {
			return j.interpolateEnvVars(sstring.(string), additionalContext)
		})
	} else if amap, ok := definition.(map[string]interface{}); ok {
		// it is a map! let's check if there are known keys we're expecting
		if envMap, exists := amap["env"]; exists {
			// if it doesn't exist then it's not a big deal
			for envName, envValue := range envMap.(map[string]interface{}) {
				envValueString := fmt.Sprintf("%s", envValue)
				// now let's treat the envValueString as a base for our template
				// and pass it trough the template rendering function.
				envValueString = j.interpolateEnvVars(envValueString, additionalContext)
				// finally, append it to env variables of the command
				env = append(env, fmt.Sprintf("%s=%s", strings.ToUpper(envName), envValueString))
			}
		}
		if commandsInterface, exists := amap["commands"]; exists {
			if commandsSlice, ok := commandsInterface.([]interface{}); ok {
				commands = lo.Map(commandsSlice, func(cstring interface{}, _ int) string {
					return j.interpolateEnvVars(cstring.(string), additionalContext)
				})
			}
		} else {
			return nil, fmt.Errorf("unexpected structure: %+v", amap)
		}
	} else {
		return nil, fmt.Errorf("unrecognized command type")
	}
	// great. if we're here then the command could be parsed
	var commandName string
	var commandArgs []string

	if len(commands) > 1 {
		// this is a series of commands therefore we have to execute it in shell
		command := strings.Join(commands, " && ")
		cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
		cmd.Env = env
		// set the group so that we can kill the subprocesses
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		return cmd, nil
	} else if len(commands) == 1 {
		// this is a single command therefore let's execute it directly
		commandParts := strings.SplitN(commands[0], " ", 2)
		if len(commandParts) >= 1 {
			commandName = commandParts[0]
			if len(commandParts) > 1 {
				commandArgs = strings.Split(commandParts[1], " ")
			}
			cmd := exec.CommandContext(ctx, commandName, commandArgs...)
			cmd.Env = env
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("bad input command")
}
