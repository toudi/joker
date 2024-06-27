package joker

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	fmt.Printf("Stopping processes\n")
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
				// fmt.Printf("%s was not yet checked; set recheck=true\n", serviceName)
				continue
			}
			if !recheck[serviceName] {
				// fmt.Printf("%s was already checked and the test passed\n", serviceName)
				checksLeft -= 1
				continue
			}
			if !service.HasStarted() {
				// fmt.Printf("%s did not start yet\n", serviceName)
				continue
			}
			if service.IsAlive() {
				// fmt.Printf(
				// 	"%s appears to be alive.\n",
				// 	serviceName,
				// )
				recheck[serviceName] = false
				checksLeft -= 1
				continue
			} else {
				// fmt.Printf(
				// 	"%s did not pass liveness check. stopping joker\n",
				// 	service.definition.Name,
				// )
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
