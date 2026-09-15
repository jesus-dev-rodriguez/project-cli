package command

import (
	"fmt"

	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/spf13/cobra"
)

func PathCmd(service services.ProjectService) *cobra.Command {
	return &cobra.Command{
		Use:   "path <proyecto>",
		Short: "Muestra la ruta absoluta de un proyecto",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := service.Path(args[0])
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), path)
			return err
		},
	}
}
