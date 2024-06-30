package joker

import (
	"github.com/phuslu/log"
)

func (s *Service) Clear(joker *Joker, force bool) error {
	if s.definition.Cleanup != nil {
		return joker.state.ClearBootstrapped(s.definition.Name, force, func() error {
			log.Debug().Str("service", s.definition.Name).Msg("clear")

			command, err := joker.prepareCommand(joker.ctx, s.definition.Cleanup)

			log.Debug().Str("command", command.String())

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
