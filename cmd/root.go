/*
Copyright © 2024 Nigel Tatschner <ntatschner@gmail.com>
*/
package cmd

import (
	"os"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/mainview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GoPowerShellLauncher",
	Short: "Launches PowerShell \"Profile\" scripts with ease.",
	Long: `GoPowerShellLauncher is a CLI tool that allows you to easily launch PowerShell scripts,
	called "Profiles", with a simple command. It is designed to be used in conjunction with
	"PowerShell Profile" scripts that are designed to be run in a specific environment.
	You can create shortcuts to your favorite PowerShell Profile scripts, and launch them.`,
	Run: func(cmd *cobra.Command, args []string) {
		l.Logger.Info("Launching PowerShell Launcher UI")
		w, h := utils.GetWindowSize()
		tprogram := tea.NewProgram(mainview.NewMainModel(tea.WindowSizeMsg{Width: w, Height: h}), tea.WithAltScreen(), tea.WithFPS(120))
		if _, err := tprogram.Run(); err != nil {
			l.Logger.Error("Error starting the program", "Error", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
