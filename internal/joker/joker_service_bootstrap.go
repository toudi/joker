package joker

import (
	"context"
	"errors"
	"fmt"

	"github.com/phuslu/log"
)

func (s *Service) bootstrap(ctx context.Context, joker *Joker) error {
	// does the service have any bootstrapping code?
	// if it does, let's check if it needs to be executed. maybe it was
	// already bootstrapped?
	if s.definition.Bootstrap != nil {
		if err := joker.state.SetBootstrapped(s.definition.Name, func() error {
			log.Info().Str("service", s.definition.Name).Msg("bootstrap")

			command, err := joker.prepareCommand(ctx, s.definition.Bootstrap, s.templateContext())

			if err != nil {
				return err
			}

			if s.definition.Dir != "" {
				command.Dir = s.definition.Dir
				log.Debug().Str("dir", s.definition.Name).Msg("set working dir")
			}

			handleCommandStream(joker, s.definition.Name, command)

			if err = command.Run(); err != nil {
				return errors.Join(fmt.Errorf("error running bootstrap command"), err)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
