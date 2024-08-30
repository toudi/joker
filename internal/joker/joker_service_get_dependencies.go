package joker

import (
	"slices"

	"github.com/samber/lo"
)

func (j *Joker) getServiceByName(serviceName string) (*Service, bool) {
	return lo.Find(
		j.services,
		func(item *Service) bool { return item.definition.Name == serviceName },
	)
}

func (j *Joker) getServiceWithDependencies(serviceName string) []*Service {
	// is this a valid service ?
	var result []*Service

	for _, service := range j.services {
		if service.definition.Name == serviceName ||
			slices.Contains(service.definition.Depends, serviceName) {
			result = append(result, service)
		}
	}

	// now we have to sort the found services so that all the dependencies
	// will be before the actual service.
	slices.SortFunc(result, func(a, b *Service) int {
		// we want to bubble the service itself to the end of the list
		if a.definition.Name == serviceName {
			return 1
		}
		if b.definition.Name == serviceName {
			return -1
		}
		// otherwise we don't care
		return 0
	})

	return result
}
