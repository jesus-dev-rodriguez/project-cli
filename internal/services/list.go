package services

import (
	"fmt"
	"os"
	"sort"
)

func (s *projectService) List() ([]string, error) {
	entries, err := os.ReadDir(s.rootDir)
	if err != nil {
		return nil, fmt.Errorf("leer el directorio de proyectos %q: %w", s.rootDir, err)
	}

	projects := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			projects = append(projects, entry.Name())
		}
	}

	sort.Strings(projects)
	return projects, nil
}
