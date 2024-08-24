package joker

func (j *Joker) Down() error {
	for i := len(j.services) - 1; i > -1; i -= 1 {
		var service = j.services[i]
		if err := j.services[i].Down(service.shutdownOptions(true)); err != nil {
			return err
		}
	}
	j.runShutdownHandlers()
	return nil
}
