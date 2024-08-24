package joker

import (
	"github.com/samber/lo"
)

const rpcCmdRestartService = "restart"

func rpcCmdRestartServiceHandler(j *Joker, args string) error {
	shutdownOptions, serviceName, err := parseShutdownOptionsAndService(args)
	if err != nil {
		return err
	}

	service, found := lo.Find(
		j.services,
		func(s *Service) bool { return s.definition.Name == serviceName },
	)

	if !found {
		return errUnknownService
	}

	return service.Restart(j, shutdownOptions)
}
