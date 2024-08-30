package joker

const rpcCmdRestartService = "restart"

func rpcCmdRestartServiceHandler(j *Joker, args string) error {
	shutdownOptions, serviceName, err := parseShutdownOptionsAndService(args)
	if err != nil {
		return err
	}

	if err := j.StopService(serviceName, shutdownOptions); err != nil {
		return err
	}

	return j.StartService(serviceName, shutdownOptions.withDependencies)
}
