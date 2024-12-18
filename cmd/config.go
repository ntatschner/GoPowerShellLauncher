package cmd

import (
	"github.com/spf13/cobra"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Create a new user directory configuration file",
	Long:  `This command creates a new user directory configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		l.Logger.Info("Creating a new configuration file")
		for _, c := range utils.ConfigStoreData {
			if c.Exists {
				l.Logger.Info("Configuration file exists", "path", c.Path)
				print("Configuration file exists at: ", c.Path)
				return
			}
		}
	},
}

func init() {
	// add the commands to the root command
	rootCmd.AddCommand(configCmd)
}
