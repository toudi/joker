package joker

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/phuslu/log"
)

var signalChannel chan os.Signal

func (j *Joker) Up() error {
	// after a lot of googling I was able to understand the pattern.
	// first of all, we're setting up the channel for listening for
	// system events:
	signalChannel = make(chan os.Signal, 1)
	// next, we're setting up a "done" channel, just in order to block
	// this function (Up)
	done := make(chan bool, 1)
	// we're setting up two signals that we want to listen to:
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	j.StreamHandler()

	// the key point is to listen for these signals *before* any of the
	// processes are spawned up.
	// Second piece of the puzzle is to set the gid attribute of the process:
	// https://stackoverflow.com/questions/35433741/in-golang-prevent-child-processes-to-receive-signals-from-calling-process
	// Otherwise the processes will be direct children of joker process and as a
	// result all of the signals that joker receives will be propagated to the
	// subprocesses as well, which is not what we want.
	go func() {
		<-signalChannel
		signal.Stop(signalChannel)

		log.Debug().Msg("stopping processes")

		_ = j.Down()
		// notify the main program that it can gracefuly finish
		// as we're done with processing signals.
		done <- true
	}()

	// now we can start the processes.
	for _, service := range j.services {
		if err := service.Up(j.ctx, j, serviceStartOptions{}); err != nil {
			log.Error().Err(err).Msg("unable to instantiate project. stopping joker")
			return j.Down()
		}
		// this will no-op if the state doesn't need to be saved but it does
		// introduce the benefit that if a service was successfully bootstrapped
		// then the state will reflect that ASAP.
		_ = j.state.Save()
	}

	// only start these once we know the project is up
	go j.livenessCheck(signalChannel)
	if err := j.startRPCListener(); err != nil {
		signalChannel <- syscall.SIGTERM
	}

	<-done

	return nil
}

func (j *Joker) Shutdown() {
	signalChannel <- syscall.SIGTERM
}
