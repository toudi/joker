package joker

const rpcCmdShutdown = "shutdown"

func rpcCmdShutdownHandler(j *Joker, _ string) error {
	j.Shutdown()
	return nil
}
