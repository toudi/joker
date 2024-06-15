package joker

func (j *Joker) Down() error {
	for i := len(j.services) - 1; i > -1; i -= 1 {
		if err := j.services[i].Down(); err != nil {
			return err
		}
	}
	j.runShutdownHandlers()
	return nil
}
