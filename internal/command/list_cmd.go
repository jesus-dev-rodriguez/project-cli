package command

import (
	"fmt"

	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/spf13/cobra"
)

func ListCmd(service services.ProjectService) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lista los proyectos disponibles",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			projects, err := service.List()
			if err != nil {
				return err
			}

			for _, name := range projects {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), name); err != nil {
					return fmt.Errorf("escribir la lista de proyectos: %w", err)
				}
			}

			return nil
		},
	}
}
