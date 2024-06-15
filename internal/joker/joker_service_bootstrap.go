package joker

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) bootstrap(ctx context.Context, joker *Joker) error {
	// does the service have any bootstrapping code?
	// if it does, let's check if it needs to be executed. maybe it was
	// already bootstrapped?
	if s.definition.Bootstrap != nil {
		if err := joker.state.SetBootstrapped(s.definition.Name, func() error {
			fmt.Printf("bootstrapping %s\n", s.definition.Name)
			command, err := joker.prepareCommand(ctx, s.definition.Bootstrap)
			if err != nil {
				return err
			}
			if s.definition.Dir != "" {
				command.Dir = s.definition.Dir
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
