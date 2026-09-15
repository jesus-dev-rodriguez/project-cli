package command

import "github.com/spf13/cobra"

func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pjt",
		Short: "Gestiona tus proyectos locales",
	}
}
