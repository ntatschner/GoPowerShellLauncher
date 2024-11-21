package mainview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

var (
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF")).Italic(true)
)

type sessionState int

const (
	mainView sessionState = iota
	profileView
	confirmationView
)

// MainView is the main view of the application
type model struct {
	state    sessionState
	profiles tea.Model
	shells   tea.Model
}

func (m model) Init() tea.Cmd {
	l.Logger.Info("Initializing MainView")
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("GoPowerShellLauncher"),
		subtleStyle.Render("Press q to quit"),
	)
}
