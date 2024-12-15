package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FF99FF")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#0F99EE")).Bold(true)
)

var shortcutCmd = &cobra.Command{
	Use:   "shortcut",
	Short: "Provides a way to create a new PowerShell launcher shortcut",
	Long:  `This command provides a way to create a new PowerShell launcher shortcut.`,
	Run: func(cmd *cobra.Command, args []string) {
		l.Logger.Info("Creating a new launcher shortcut")
		// Do Stuff Here
	},
}

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Loads the specified profile directly in the shell denoted by the profile",
	Long:  `This command loads the specified profile directly in the shell denoted by the profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		l.Logger.Info("Loading the specified profile")
		// Do Stuff Here
	},
}

func init() {
	// flags for the shortcut command
	shortcutCmd.Flags().StringP("name", "n", "", "The name of the shortcut")
	shortcutCmd.Flags().StringP("path", "p", "", "The path to the PowerShell script")
	// command configs
	shortcutCmd.MarkFlagRequired("name")
	shortcutCmd.MarkFlagRequired("path")

	// flags for the profiles command
	profilesCmd.Flags().StringP("path", "p", "", "The path to the profile")
	// command configs
	profilesCmd.MarkFlagRequired("path")
	// add the commands to the root command
	rootCmd.AddCommand(profilesCmd)
	rootCmd.AddCommand(shortcutCmd)
}
