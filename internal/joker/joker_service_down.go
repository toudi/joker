package joker

import (
	"strings"
	"syscall"

	"github.com/phuslu/log"
	"golang.org/x/sys/unix"
)

func (s *Service) Down() error {
	log.Debug().Str("service", s.definition.Name).Msg("down")

	if s.process != nil && s.IsAlive() {
		// this is a shell subprocess
		if s.process.SysProcAttr != nil && s.process.SysProcAttr.Setpgid {
			return syscall.Kill(-s.process.Process.Pid, syscall.SIGKILL)
		}
		// this is a regular process
		return s.process.Process.Signal(s.getKillSignal())
	} else {
		log.Debug().Str("service", s.definition.Name).Msg("does not need to be killed")
	}
	return nil
}

func (s *Service) getKillSignal() syscall.Signal {
	var signal = syscall.SIGTERM

	if s.definition.KillSignal == nil {
		log.Trace().Msg("signal not defined; using SIGTERM")
		return signal
	}

	if signalName, ok := s.definition.KillSignal.(string); ok {
		log.Trace().Str("signal", signalName).Msg("signal defined as a string")

		if !strings.HasPrefix(signalName, "SIG") {
			signalName = "SIG" + signalName
		}

		if parsedSignal := unix.SignalNum(signalName); parsedSignal > 0 {
			signal = parsedSignal
		}
	} else if signalNo, ok := s.definition.KillSignal.(int); ok {
		log.Trace().Int("signal", signalNo).Msg("signal defined as int")

		tmpSignal := syscall.Signal(signalNo)
		if unix.SignalName(tmpSignal) != "" {
			signal = tmpSignal
		}
	}

	log.Trace().Str("signal", unix.SignalName(signal)).Msg("returned")

	return signal
}
