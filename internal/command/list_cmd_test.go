package command

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
)

func TestListCmd(t *testing.T) {
	rootDir := t.TempDir()
	for _, name := range []string{"second", "first"} {
		if err := os.Mkdir(filepath.Join(rootDir, name), 0o755); err != nil {
			t.Fatalf("create project %q: %v", name, err)
		}
	}
	service := services.NewProjectService(config.Config{ProjectsDir: rootDir, Editor: "sh"}, nil)

	cmd := ListCmd(service)
	var output bytes.Buffer
	cmd.SetOut(&output)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := output.String(), "first\nsecond\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
