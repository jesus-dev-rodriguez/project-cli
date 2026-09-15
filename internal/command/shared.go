package command

import (
	"github.com/jesus-dev-rodriguez/project-cli/internal/services"
	"github.com/spf13/cobra"
)

func completeProjects(service services.ProjectService) cobra.CompletionFunc {
	return func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		projects, err := service.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return projects, cobra.ShellCompDirectiveNoFileComp
	}
}
