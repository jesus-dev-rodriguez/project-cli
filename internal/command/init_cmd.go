package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
	"github.com/spf13/cobra"
)

func InitCmd(in io.Reader, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Configura la carpeta de proyectos y el editor",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := promptConfig(bufio.NewReader(in), out)
			if err != nil {
				return err
			}
			if err := config.Save(cfg); err != nil {
				return err
			}
			_, err = fmt.Fprintln(out, "Configuración guardada correctamente.")
			return err
		},
	}
}

func promptConfig(reader *bufio.Reader, out io.Writer) (config.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config.Config{}, fmt.Errorf("resolver el directorio home: %w", err)
	}

	defaultProjectsDir := filepath.Join(homeDir, "Proyectos")
	if _, err := fmt.Fprintf(out, "Carpeta de proyectos [%s]: ", defaultProjectsDir); err != nil {
		return config.Config{}, err
	}
	projectsDir, err := readLine(reader)
	if err != nil {
		return config.Config{}, fmt.Errorf("leer la carpeta de proyectos: %w", err)
	}
	if projectsDir == "" {
		projectsDir = defaultProjectsDir
	}
	if strings.HasPrefix(projectsDir, "~/") {
		projectsDir = filepath.Join(homeDir, strings.TrimPrefix(projectsDir, "~/"))
	}
	projectsDir, err = filepath.Abs(filepath.Clean(projectsDir))
	if err != nil {
		return config.Config{}, fmt.Errorf("resolver la carpeta de proyectos: %w", err)
	}
	if err := os.MkdirAll(projectsDir, 0o755); err != nil {
		return config.Config{}, fmt.Errorf("crear la carpeta de proyectos: %w", err)
	}

	editors := config.AvailableEditors()
	editor, err := selectEditor(reader, out, editors)
	if err != nil {
		return config.Config{}, err
	}

	return config.Config{ProjectsDir: projectsDir, Editor: editor}, nil
}

func selectEditor(reader *bufio.Reader, out io.Writer, editors []string) (string, error) {
	if len(editors) == 0 {
		if _, err := fmt.Fprintln(out, "No se encontró un editor conocido instalado."); err != nil {
			return "", err
		}
		if _, err := fmt.Fprint(out, "Escribe el ejecutable de un editor disponible: "); err != nil {
			return "", err
		}
		editor, err := readLine(reader)
		if err != nil {
			return "", fmt.Errorf("leer el editor: %w", err)
		}
		if editor == "" {
			return "", fmt.Errorf("debes indicar un editor")
		}
		if _, err := execLookPath(editor); err != nil {
			return "", fmt.Errorf("el editor %q no está disponible: %w", editor, err)
		}
		return editor, nil
	}

	if _, err := fmt.Fprintln(out, "Selecciona tu editor:"); err != nil {
		return "", err
	}
	for index, editor := range editors {
		if _, err := fmt.Fprintf(out, "%d) %s\n", index+1, editor); err != nil {
			return "", err
		}
	}
	if _, err := fmt.Fprint(out, "Opción: "); err != nil {
		return "", err
	}
	value, err := readLine(reader)
	if err != nil {
		return "", fmt.Errorf("leer la opción del editor: %w", err)
	}
	index, err := strconv.Atoi(value)
	if err != nil || index < 1 || index > len(editors) {
		return "", fmt.Errorf("opción de editor inválida: %q", value)
	}

	return editors[index-1], nil
}

var execLookPath = func(file string) (string, error) {
	return exec.LookPath(file)
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
