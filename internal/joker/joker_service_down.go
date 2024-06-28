package joker

import (
	"syscall"

	"github.com/phuslu/log"
)

func (s *Service) Down() error {
	log.Debug().Str("service", s.definition.Name).Msg("down")

	if s.process != nil && s.IsAlive() {
		// this is a shell subprocess
		if s.process.SysProcAttr != nil && s.process.SysProcAttr.Setpgid {
			return syscall.Kill(-s.process.Process.Pid, syscall.SIGKILL)
		}
		// this is a regular process
		return s.process.Process.Signal(syscall.SIGTERM)
	} else {
		log.Debug().Str("service", s.definition.Name).Msg("does not need to be killed")
	}
	return nil
}
