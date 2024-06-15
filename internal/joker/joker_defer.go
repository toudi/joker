package joker

func (j *Joker) Defer(f func()) {
	j.shutdownFunctions = append(j.shutdownFunctions, f)
}

func (j *Joker) runShutdownHandlers() {
	for _, function := range j.shutdownFunctions {
		function()
	}
}
