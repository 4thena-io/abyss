package cli

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "abyss",
		Short: "Abyss - Developer platform",
	}

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file (overrides ABYSS_CONFIG and ABYSS_HOME)")
	cmd.AddCommand(serveCmd(&configPath))

	return cmd
}
