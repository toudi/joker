package joker

import (
	"context"
	"fmt"
)

func (s *Service) Clear(ctx context.Context, joker *Joker) error {
	if s.definition.Cleanup != nil {
		return joker.state.ClearBootstrapped(s.definition.Name, func() error {

			if err := s.prepareDir(joker); err != nil {
				return err
			}

			command, err := joker.prepareCommand(ctx, s.definition.Cleanup)
			fmt.Printf("command: %v\n", command)
			if s.definition.Dir != "" {
				command.Dir = s.definition.Dir
			}
			if err != nil {
				return err
			}
			return command.Run()
		})
	}
	return nil
}
