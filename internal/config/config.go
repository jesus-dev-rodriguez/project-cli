package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const configDirName = "project-cli"

type Config struct {
	ProjectsDir string `json:"projects_dir"`
	Editor      string `json:"editor"`
}

func Path() (string, error) {
	baseDir := os.Getenv("XDG_CONFIG_HOME")
	if baseDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolver el directorio home: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(baseDir, configDirName, "config.json"), nil
}

func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("leer la configuración %q: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decodificar la configuración %q: %w", path, err)
	}
	if err := Validate(cfg); err != nil {
		return Config{}, fmt.Errorf("validar la configuración %q: %w", path, err)
	}

	return cfg, nil
}

func Save(cfg Config) error {
	if err := Validate(cfg); err != nil {
		return fmt.Errorf("validar la configuración: %w", err)
	}

	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("crear el directorio de configuración: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("codificar la configuración: %w", err)
	}
	data = append(data, '\n')

	temp, err := os.CreateTemp(filepath.Dir(path), "config-*.tmp")
	if err != nil {
		return fmt.Errorf("crear archivo temporal de configuración: %w", err)
	}
	tempName := temp.Name()
	defer os.Remove(tempName)

	if err := temp.Chmod(0o600); err != nil {
		temp.Close()
		return fmt.Errorf("proteger el archivo de configuración: %w", err)
	}
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return fmt.Errorf("escribir la configuración: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("cerrar la configuración temporal: %w", err)
	}
	if err := os.Rename(tempName, path); err != nil {
		return fmt.Errorf("guardar la configuración %q: %w", path, err)
	}

	return nil
}

func Validate(cfg Config) error {
	if cfg.ProjectsDir == "" || !filepath.IsAbs(cfg.ProjectsDir) {
		return errors.New("projects_dir debe ser una ruta absoluta")
	}
	info, err := os.Stat(cfg.ProjectsDir)
	if err != nil {
		return fmt.Errorf("comprobar projects_dir: %w", err)
	}
	if !info.IsDir() {
		return errors.New("projects_dir no es un directorio")
	}
	if strings.TrimSpace(cfg.Editor) == "" {
		return errors.New("editor no puede estar vacío")
	}
	if _, err := exec.LookPath(cfg.Editor); err != nil {
		return fmt.Errorf("editor %q no está disponible: %w", cfg.Editor, err)
	}

	return nil
}

var knownEditors = []string{"code", "cursor", "zed", "nvim", "vim", "emacs"}

func AvailableEditors() []string {
	available := make([]string, 0, len(knownEditors))
	for _, editor := range knownEditors {
		if _, err := exec.LookPath(editor); err == nil {
			available = append(available, editor)
		}
	}
	return available
}
