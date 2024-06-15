package joker

import (
	"fmt"
	"syscall"
)

func (s *Service) Down() error {
	fmt.Printf("stopping %s\n", s.definition.Name)

	if s.process != nil && s.IsAlive() {
		// this is a shell subprocess
		if s.process.SysProcAttr != nil && s.process.SysProcAttr.Setpgid {
			return syscall.Kill(-s.process.Process.Pid, syscall.SIGKILL)
		}
		// this is a regular process
		return s.process.Process.Kill()
	} else {
		fmt.Printf("[%s] does not need to be killed.\n", s.definition.Name)
	}
	return nil
}
