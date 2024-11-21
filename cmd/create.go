package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FF99FF")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#0F99EE")).Bold(true)
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new launcher shortcut",
	Long:  `Create a new PowerShell launcher shortcut with the given name and options selected herein.`,
	Run: func(cmd *cobra.Command, args []string) {
		GoPowerShellLauncher.Logger.Info("Creating a new launcher shortcut")
		// Do Stuff Here
	},
}
