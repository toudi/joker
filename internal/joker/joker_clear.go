package joker

func (j *Joker) Clear(serviceName string, force bool) error {
	var err error

	if err = j.Down(); err != nil {
		return err
	}

	for _, service := range j.services {
		// allow to clear only a specific service
		if serviceName != "" && service.definition.Name != serviceName {
			continue
		}
		if err = service.Clear(j, force); err != nil {
			return err
		}
	}

	return j.state.Save()
}
