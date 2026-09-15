package services

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jesus-dev-rodriguez/project-cli/internal/process"
	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
)

type fakeRunner struct {
	name string
	args []string
}

func (r *fakeRunner) Run(_ context.Context, name string, args ...string) error {
	r.name = name
	r.args = args
	return nil
}

var _ process.Runner = (*fakeRunner)(nil)

func TestOpenService(t *testing.T) {
	rootDir := t.TempDir()
	projectPath := filepath.Join(rootDir, "my-app")
	if err := os.Mkdir(projectPath, 0o755); err != nil {
		t.Fatalf("create project: %v", err)
	}

	runner := new(fakeRunner)
	service := NewProjectService(config.Config{ProjectsDir: rootDir, Editor: "code"}, runner)
	if err := service.Open(context.Background(), "my-app"); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	if runner.name != "code" || !reflect.DeepEqual(runner.args, []string{projectPath}) {
		t.Fatalf("runner = (%q, %v), want (%q, [%q])", runner.name, runner.args, "code", projectPath)
	}
}

func TestOpenServiceRejectsInvalidName(t *testing.T) {
	service := NewProjectService(config.Config{ProjectsDir: t.TempDir(), Editor: "code"}, new(fakeRunner))

	err := service.Open(context.Background(), "../outside")
	if err == nil {
		t.Fatal("Open() error = nil, want error")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Open() error = %v, want invalid-name error", err)
	}
}
