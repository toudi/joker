package joker

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/phuslu/log"
)

func (j *Joker) Up() error {
	j.StreamHandler()

	for _, service := range j.services {
		if err := service.Up(j.ctx, j); err != nil {
			return err
		}
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go j.livenessCheck(signalChannel)

	<-signalChannel

	j.ctx.Done()

	log.Debug().Msg("stopping processes")

	return j.Down()
}
