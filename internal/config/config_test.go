package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPathUsesXDGConfigHome(t *testing.T) {
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)

	got, err := Path()
	if err != nil {
		t.Fatalf("Path() error = %v", err)
	}

	want := filepath.Join(configHome, "project-cli", "config.json")
	if got != want {
		t.Fatalf("Path() = %q, want %q", got, want)
	}
}

func TestSaveAndLoad(t *testing.T) {
	configHome := t.TempDir()
	projectsDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)

	want := Config{ProjectsDir: projectsDir, Editor: "sh"}
	if err := Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != want {
		t.Fatalf("Load() = %+v, want %+v", got, want)
	}

	info, err := os.Stat(filepath.Join(configHome, "project-cli", "config.json"))
	if err != nil {
		t.Fatalf("stat config: %v", err)
	}
	if permissions := info.Mode().Perm(); permissions != 0o600 {
		t.Fatalf("config permissions = %o, want 600", permissions)
	}
}

func TestLoadRejectsInvalidJSON(t *testing.T) {
	configHome := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configHome)
	if err := os.MkdirAll(filepath.Join(configHome, "project-cli"), 0o700); err != nil {
		t.Fatalf("create config directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configHome, "project-cli", "config.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load() error = %v, want JSON error", err)
	}
}
