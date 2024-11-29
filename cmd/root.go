/*
Copyright © 2024 Nigel Tatschner <ntatschner@gmail.com>
*/
package cmd

import (
	"os"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/mainview"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// Get size of the terminal window

func GetWindowSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	l.Logger.Info("Getting terminal size", "Width", width, "Height", height)
	if err != nil {
		l.Logger.Error("Error getting terminal size", "Error", err)
		return 0, 0
	}
	return width, height
}

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
		w, h := GetWindowSize()
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.GoPowerShellLauncher.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
