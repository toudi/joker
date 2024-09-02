package joker

import (
	"syscall"
	"time"

	"github.com/phuslu/log"
	"github.com/toudi/joker/internal/utils"
	"golang.org/x/sys/unix"
)

func (s *Service) Down(options serviceShutdownOptions) error {
	log.Debug().Str("service", s.definition.Name).Msg("down")

	if s.process != nil && s.IsAlive() {
		if err := syscall.Kill(-s.process.Process.Pid, options.signal); err != nil {
			return err
		}

		if options.wait {
			return s.IsAliveSentinel(false, 5*time.Second)
		}
	} else {
		log.Debug().Str("service", s.definition.Name).Msg("does not need to be killed")
	}
	return nil
}

func (s *Service) shutdownOptions(wait bool) serviceShutdownOptions {
	return serviceShutdownOptions{
		signal: s.getKillSignal(),
		wait:   wait,
	}
}

func (s *Service) getKillSignal() syscall.Signal {
	var signal = syscall.SIGTERM

	if s.definition.KillSignal == nil {
		log.Trace().Msg("signal not defined; using SIGTERM")
		return signal
	}

	if signalName, ok := s.definition.KillSignal.(string); ok {
		signal = utils.ParseSignalFromString(signalName)
	} else if signalNo, ok := s.definition.KillSignal.(int); ok {
		signal = utils.ParseSignalFromInt(signalNo)
	}

	log.Trace().Str("signal", unix.SignalName(signal)).Msg("returned")

	return signal
}
