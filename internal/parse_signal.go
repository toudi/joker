package internal

import (
	"strings"
	"syscall"

	"github.com/phuslu/log"
	"golang.org/x/sys/unix"
)

func ParseSignalFromString(signalName string) syscall.Signal {
	var signal = syscall.SIGTERM

	log.Trace().Str("signal", signalName).Msg("signal defined as a string")

	if !strings.HasPrefix(signalName, "SIG") {
		signalName = "SIG" + signalName
	}

	if parsedSignal := unix.SignalNum(signalName); parsedSignal > 0 {
		signal = parsedSignal
	}

	return signal
}

func ParseSignalFromInt(signalNo int) syscall.Signal {
	var signal = syscall.SIGTERM

	log.Trace().Int("signal", signalNo).Msg("signal defined as int")

	tmpSignal := syscall.Signal(signalNo)
	if unix.SignalName(tmpSignal) != "" {
		signal = tmpSignal
	}

	return signal
}
