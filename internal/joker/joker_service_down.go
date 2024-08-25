package joker

import (
	"errors"
	"syscall"
	"time"

	"github.com/phuslu/log"
	"github.com/toudi/joker/internal/utils"
	"golang.org/x/sys/unix"
)

var errTimeoutShuttingDownService = errors.New("timeout waiting for service to be shut down")

func (s *Service) Down(options serviceShutdownOptions) error {
	log.Debug().Str("service", s.definition.Name).Msg("down")

	if s.process != nil && s.IsAlive() {
		// this is a shell subprocess
		if s.process.SysProcAttr != nil && s.process.SysProcAttr.Setpgid {
			return syscall.Kill(-s.process.Process.Pid, syscall.SIGKILL)
		}

		// this is a regular process
		if err := s.process.Process.Signal(options.signal); err != nil {
			return err
		}

		if options.wait {
			timerTimeout := time.NewTicker(5 * time.Second)
			timerPoll := time.NewTicker(200 * time.Millisecond)

			defer func() {
				timerTimeout.Stop()
				timerPoll.Stop()
			}()

			for {
				select {
				case <-timerPoll.C:
					log.Trace().
						Str("service", s.definition.Name).
						Msg("waiting for process to finish")
					if !s.IsAlive() {
						return nil
					}
				case <-timerTimeout.C:
					log.Error().Msg("timeout reached")
					return errTimeoutShuttingDownService
				}
			}
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
