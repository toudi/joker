package joker

import (
	"bufio"
	"os/exec"
)

func handleCommandStream(joker *Joker, serviceName string, command *exec.Cmd) error {
	// https://gist.github.com/mxschmitt/6c07b5b97853f05455c3fdaf48b1a8b6

	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}
	// command.Stderr = command.Stdout
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go streamHandler(stdoutScanner, serviceName, joker, false)
	go streamHandler(stderrScanner, serviceName, joker, true)

	return nil
}
