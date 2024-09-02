package joker

import "syscall"

type serviceShutdownOptions struct {
	wait bool
	// should joker wait for command to complete
	signalInput string
	// this one wil be passed by the user
	signal syscall.Signal
	// this one will be passed to the service after parsing from the string one
	// defaults to SIGTERM
	withDependencies bool
	// WARNING: not implemented for now, but a potential improvement.
	// the idea would be that if you pass -deps flag, then joker would
	// stop any services that depend on the service, then run the restart
	// and then reverse the process.
}

func (s *Service) Restart(joker *Joker, options serviceShutdownOptions) error {
	if err := s.Down(options); err != nil {
		return err
	}

	if err := s.Up(joker.ctx, joker, serviceStartOptions{
		WithDependencies: options.withDependencies,
		Wait:             options.wait,
	}); err != nil {
		return err
	}

	return nil
}
