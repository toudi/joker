package joker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
)

func streamHandler(scanner *bufio.Scanner, service string, joker *Joker, stderr bool) {
	for scanner.Scan() {
		joker.streamChan <- StreamLine{
			Service: service,
			Line:    scanner.Text(),
			Stderr:  stderr,
		}
	}
}

func (s *Service) prepareDir(joker *Joker) error {
	if s.definition.Dir != "" {
		s.definition.Dir = joker.interpolateEnvVars(s.definition.Dir)
		// check if the directory exists
		_, err := os.Stat(s.definition.Dir)
		if os.IsNotExist(err) {
			if err := os.MkdirAll(s.definition.Dir, 0755); err != nil {
				return errors.Join(err, errors.New("cannot create directory"))
			}
		}
	}

	return nil
}

func (s *Service) Up(ctx context.Context, joker *Joker) error {
	fmt.Printf("launching %s\n", s.definition.Name)

	if err := s.prepareDir(joker); err != nil {
		return err
	}

	if err := s.bootstrap(ctx, joker); err != nil {
		return err
	}

	command, err := joker.prepareCommand(ctx, s.definition.Command)
	if err != nil {
		return err
	}
	if s.definition.Dir != "" {
		command.Dir = s.definition.Dir
	}
	if err = handleCommandStream(joker, s.definition.Name, command); err != nil {
		return err
	}

	s.process = command

	if err = s.process.Start(); err != nil {
		return err
	}

	go func() {
		_ = s.process.Wait()
	}()

	// has the hot reload been asked for ?
	if s.definition.HotReload != nil {
		if err = s.HotReloadHandler(joker); err != nil {
			return err
		}
	}

	return nil
}
