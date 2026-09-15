package services

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
)

func TestProjectServiceList(t *testing.T) {
	rootDir := t.TempDir()
	for _, name := range []string{"zeta", "alpha", "beta"} {
		if err := os.Mkdir(filepath.Join(rootDir, name), 0o755); err != nil {
			t.Fatalf("create project %q: %v", name, err)
		}
	}
	if err := os.WriteFile(filepath.Join(rootDir, "README.md"), []byte("file"), 0o644); err != nil {
		t.Fatalf("create file: %v", err)
	}

	got, err := NewProjectService(config.Config{
		ProjectsDir: rootDir,
		Editor:      "sh",
	}, nil).List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := []string{"alpha", "beta", "zeta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("List() = %v, want %v", got, want)
	}
}
