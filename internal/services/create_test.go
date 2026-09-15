package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
)

func TestProjectServiceCreateCreatesAndOpensProject(t *testing.T) {
	rootDir := t.TempDir()
	runner := new(fakeRunner)
	service := NewProjectService(config.Config{
		ProjectsDir: rootDir,
		Editor:      "code",
	}, runner)

	if err := service.Create(context.Background(), "new-app"); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	projectPath := filepath.Join(rootDir, "new-app")
	info, err := os.Stat(projectPath)
	if err != nil {
		t.Fatalf("stat created project: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("created project is not a directory")
	}
	if runner.name != "code" || len(runner.args) != 1 || runner.args[0] != projectPath {
		t.Fatalf("runner = (%q, %v), want (%q, [%q])", runner.name, runner.args, "code", projectPath)
	}
}

func TestProjectServiceCreateRejectsExistingProject(t *testing.T) {
	rootDir := t.TempDir()
	projectPath := filepath.Join(rootDir, "existing")
	if err := os.Mkdir(projectPath, 0o755); err != nil {
		t.Fatalf("create existing project: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectPath, "keep.txt"), []byte("keep"), 0o644); err != nil {
		t.Fatalf("write existing content: %v", err)
	}

	runner := new(fakeRunner)
	service := NewProjectService(config.Config{
		ProjectsDir: rootDir,
		Editor:      "code",
	}, runner)
	if err := service.Create(context.Background(), "existing"); err == nil {
		t.Fatal("Create() error = nil, want error")
	}
	if runner.name != "" {
		t.Fatal("editor was run for an existing project")
	}
}

func TestProjectServiceCreateRejectsInvalidName(t *testing.T) {
	runner := new(fakeRunner)
	service := NewProjectService(config.Config{
		ProjectsDir: t.TempDir(),
		Editor:      "code",
	}, runner)

	if err := service.Create(context.Background(), "../outside"); err == nil {
		t.Fatal("Create() error = nil, want error")
	}
	if runner.name != "" {
		t.Fatal("editor was run for an invalid project name")
	}
}
