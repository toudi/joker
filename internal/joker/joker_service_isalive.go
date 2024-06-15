package joker

func (s *Service) HasStarted() bool {
	return s.process.Process != nil
}
func (s *Service) IsAlive() bool {
	return s.HasStarted() && s.process.ProcessState == nil
}
