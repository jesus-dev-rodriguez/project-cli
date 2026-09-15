package command

import (
	"fmt"

	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/spf13/cobra"
)

func CreateCmd(service services.ProjectService) *cobra.Command {
	return &cobra.Command{
		Use:   "create <proyecto>",
		Short: "Crea y abre un proyecto nuevo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := service.Create(cmd.Context(), args[0]); err != nil {
				return err
			}

			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Proyecto %q creado y abierto.\n", args[0])
			return err
		},
	}
}
