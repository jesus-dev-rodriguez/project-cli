package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (s *projectService) Create(ctx context.Context, name string) error {
	if err := validateProjectName(name); err != nil {
		return err
	}

	projectPath := filepath.Join(s.rootDir, name)
	if err := os.Mkdir(projectPath, 0o755); err != nil {
		return fmt.Errorf("crear el proyecto %q: %w", name, err)
	}

	if err := s.runner.Run(ctx, s.editor, projectPath); err != nil {
		return fmt.Errorf("abrir el proyecto recién creado %q con %q: %w", name, s.editor, err)
	}

	return nil
}

func validateProjectName(name string) error {
	if name == "" || name == "." || name == ".." ||
		filepath.Base(name) != name ||
		strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("nombre de proyecto inválido: %q", name)
	}

	return nil
}
