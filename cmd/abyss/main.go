package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}

func rootCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "abyss",
		Short: "Abyss - Developer platform",
	}

	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file (default: /etc/abyss/config.yaml)")
	cmd.AddCommand(serveCmd(&configPath))

	return cmd
}
