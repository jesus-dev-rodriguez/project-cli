package services

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
)

func TestProjectServicePath(t *testing.T) {
	rootDir := t.TempDir()
	projectPath := filepath.Join(rootDir, "my-app")
	if err := os.Mkdir(projectPath, 0o755); err != nil {
		t.Fatalf("create project: %v", err)
	}

	service := NewProjectService(config.Config{
		ProjectsDir: rootDir,
		Editor:      "code",
	}, new(fakeRunner))
	got, err := service.Path("my-app")
	if err != nil {
		t.Fatalf("Path() error = %v", err)
	}
	if got != projectPath {
		t.Fatalf("Path() = %q, want %q", got, projectPath)
	}
}

func TestProjectServicePathRejectsMissingProject(t *testing.T) {
	service := NewProjectService(config.Config{
		ProjectsDir: t.TempDir(),
		Editor:      "code",
	}, new(fakeRunner))

	if _, err := service.Path("missing"); err == nil {
		t.Fatal("Path() error = nil, want error")
	}
}
