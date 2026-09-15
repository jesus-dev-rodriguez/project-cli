package services

import (
	"context"

	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
	"github.com/jesus-dev-rodriguez/project-cli/internal/process"
)

type ProjectService interface {
	List() ([]string, error)
	Open(ctx context.Context, name string) error
	Create(ctx context.Context, name string) error
	Path(name string) (string, error)
}

type projectService struct {
	rootDir       string
	runner        process.Runner
	editor        string
}

func NewProjectService(cfg config.Config, runner process.Runner) ProjectService {
	return &projectService{
		rootDir:       cfg.ProjectsDir,
		runner:        runner,
		editor:        cfg.Editor,
	}
}
