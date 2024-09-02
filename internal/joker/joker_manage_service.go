package joker

import "slices"

func (j *Joker) getServiceDependencyQueue(
	serviceName string,
	withDependencies bool,
) ([]*Service, error) {
	var queue []*Service

	if service, exists := j.getServiceByName(serviceName); !exists {
		return nil, errUnknownService
	} else {
		queue = append(queue, service)
	}

	if withDependencies {
		queue = j.getServiceWithDependencies(serviceName)
	}

	return queue, nil
}

func (j *Joker) StopService(serviceName string, options serviceShutdownOptions) error {
	queue, err := j.getServiceDependencyQueue(serviceName, options.withDependencies)
	if err != nil {
		return err
	}

	for _, service := range queue {
		if err := service.Down(options); err != nil {
			return err
		}
	}

	return nil
}

func (j *Joker) StartService(serviceName string, options serviceStartOptions) error {
	queue, err := j.getServiceDependencyQueue(serviceName, options.WithDependencies)
	if err != nil {
		return err
	}

	// when starting the service back, we have to reverse the queue as
	// the target service will be last on the list.

	slices.Reverse(queue)

	for _, service := range queue {
		if err := service.Up(j.ctx, j, options); err != nil {
			return err
		}
	}

	return nil

}
