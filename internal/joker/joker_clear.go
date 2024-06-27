package joker

import (
	"fmt"
)

func (j *Joker) Clear(serviceName string) error {
	var err error

	if err = j.Down(); err != nil {
		return err
	}

	for _, service := range j.services {
		// allow to clear only a specific service
		if serviceName != "" && service.definition.Name != serviceName {
			continue
		}
		fmt.Printf("calling %s::clear\n", service.definition.Name)
		if err = service.Clear(j); err != nil {
			return err
		}
	}

	return j.state.Save()
}
