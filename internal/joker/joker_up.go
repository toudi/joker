package joker

import (
	"os"
	"os/signal"
	"syscall"
	"time"

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

func (j *Joker) livenessCheck(interrupt chan os.Signal) {
	// all that we have to do is to just iterate trough
	// all the services, check that they have been
	// instantiated and make sure that their processes did
	// not exit with a non-zero exitcode.
	var finished bool = false

	var recheck = make(map[string]bool)

	for !finished {
		var checksLeft int = len(j.services)

		for _, service := range j.services {
			var serviceName = service.definition.Name
			if _, exists := recheck[serviceName]; !exists {
				recheck[serviceName] = true
				log.Debug().Str("service", serviceName).Msg("was not yet checked for liveness")
				continue
			}
			if !recheck[serviceName] {
				log.Debug().Str("service", serviceName).Msg("was already checked")
				checksLeft -= 1
				continue
			}
			if !service.HasStarted() {
				log.Debug().Str("service", serviceName).Msg("did not start yet")
				continue
			}
			if service.IsAlive() {
				log.Debug().Str("service", serviceName).Msg("appears to be alive")
				recheck[serviceName] = false
				checksLeft -= 1
				continue
			} else {
				log.Debug().Str("service", serviceName).Msg("did not pass liveness check. stopping joker")
				defer func() {
					interrupt <- syscall.SIGTERM
				}()
				checksLeft = 0
				break
			}
		}

		finished = checksLeft == 0

		time.Sleep(time.Second)
	}
}
