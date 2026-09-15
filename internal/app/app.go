package app

import (
	"io"

	"github.com/jesus-dev-rodriguez/project-cli/internal/command"
	"github.com/jesus-dev-rodriguez/project-cli/internal/config"
	"github.com/jesus-dev-rodriguez/project-cli/internal/process"
	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/spf13/cobra"
)

// Run composes the application and executes the requested command.
func Run(args []string, in io.Reader, out, errOut io.Writer) error {
	rootCmd, err := NewCommand(args, in, out, errOut)
	if err != nil {
		return err
	}
	return rootCmd.Execute()
}

// NewCommand builds the command tree and its application dependencies.
func NewCommand(args []string, in io.Reader, out, errOut io.Writer) (*cobra.Command, error) {
	cfg, err := config.Load()
	if err != nil && !isSetupCommand(args) && !isHelpCommand(args) {
		if initErr := command.InitCmd(in, out).Execute(); initErr != nil {
			return nil, initErr
		}
		cfg, err = config.Load()
	}

	rootCmd := command.RootCmd()
	rootCmd.SetArgs(args)
	rootCmd.SetIn(in)
	rootCmd.SetOut(out)
	rootCmd.SetErr(errOut)
	rootCmd.AddCommand(command.InitCmd(in, out))
	if err != nil {
		return rootCmd, nil
	}

	projectService := services.NewProjectService(
		cfg,
		process.NewRunner(in, out, errOut),
	)
	rootCmd.AddCommand(command.ListCmd(projectService))
	rootCmd.AddCommand(command.OpenCmd(projectService))
	rootCmd.AddCommand(command.CreateCmd(projectService))
	rootCmd.AddCommand(command.PathCmd(projectService))

	return rootCmd, nil
}

func isSetupCommand(args []string) bool {
	return len(args) > 0 && args[0] == "init"
}

func isHelpCommand(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" || arg == "help" {
			return true
		}
	}
	return false
}
