package joker

import (
	"bufio"
	"context"
	"errors"
	"os"

	"github.com/phuslu/log"
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
		s.definition.Dir = joker.interpolateEnvVars(s.definition.Dir, nil)
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
	log.Debug().Str("service", s.definition.Name).Msg("launching")

	if err := s.bootstrap(ctx, joker); err != nil {
		return err
	}

	command, err := joker.prepareCommand(ctx, s.definition.Command, s.templateContext())
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

	return nil
}
