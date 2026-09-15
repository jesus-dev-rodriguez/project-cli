package command

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
)

func TestInitCmd(t *testing.T) {
	configHome := t.TempDir()
	projectsDir := filepath.Join(t.TempDir(), "projects")
	t.Setenv("XDG_CONFIG_HOME", configHome)

	editorInput := "sh\n"
	if editors := config.AvailableEditors(); len(editors) > 0 {
		editorInput = "1\n"
	}
	var output bytes.Buffer
	cmd := InitCmd(strings.NewReader(fmt.Sprintf("%s\n%s", projectsDir, editorInput)), &output)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.ProjectsDir != projectsDir {
		t.Fatalf("ProjectsDir = %q, want %q", cfg.ProjectsDir, projectsDir)
	}
	if output.Len() == 0 {
		t.Fatal("init output is empty")
	}
}
