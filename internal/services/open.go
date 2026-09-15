package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func (s *projectService) Open(ctx context.Context, name string) error {
	projectPath, err := s.resolveProject(name)
	if err != nil {
		return err
	}

	if err := s.runner.Run(ctx, s.editor, projectPath); err != nil {
		return fmt.Errorf("abrir el proyecto %q con %q: %w", name, s.editor, err)
	}

	return nil
}

func (s *projectService) resolveProject(name string) (string, error) {
	if err := validateProjectName(name); err != nil {
		return "", err
	}

	projectPath := filepath.Join(s.rootDir, name)
	info, err := os.Stat(projectPath)
	if err != nil {
		return "", fmt.Errorf("buscar el proyecto %q: %w", name, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("el proyecto %q no es un directorio", name)
	}

	return projectPath, nil
}
