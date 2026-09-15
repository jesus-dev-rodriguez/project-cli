package services

func (s *projectService) Path(name string) (string, error) {
	return s.resolveProject(name)
}
