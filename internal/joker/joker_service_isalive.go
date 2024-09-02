package joker

import (
	"errors"
	"time"

	"github.com/phuslu/log"
)

var errTimeoutWaitingForService = errors.New("timeout waiting for service to be shut down")

func (s *Service) HasStarted() bool {
	return s.process.Process != nil
}
func (s *Service) IsAlive() bool {
	return s.HasStarted() && s.process.ProcessState == nil
}

// IsAliveSentinel checks whether the service's IsAlive returns expected value
// It waits `timeout` at maximum in which case the function returns timeout error.
// The function returns nil in case the `expectedValue` is met.
func (s *Service) IsAliveSentinel(expectedValue bool, timeout time.Duration) error {
	timerTimeout := time.NewTicker(timeout)
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
				Msg("waiting for process")
			if s.IsAlive() == expectedValue {
				return nil
			}
		case <-timerTimeout.C:
			log.Error().Msg("timeout reached")
			return errTimeoutWaitingForService
		}
	}

}
