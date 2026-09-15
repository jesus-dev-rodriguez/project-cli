package command

import (
	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/spf13/cobra"
)

func OpenCmd(service services.ProjectService) *cobra.Command {
	return &cobra.Command{
		Use:               "open <proyecto>",
		Short:             "Abre un proyecto con el editor configurado",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeProjects(service),
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.Open(cmd.Context(), args[0])
		},
	}
}
