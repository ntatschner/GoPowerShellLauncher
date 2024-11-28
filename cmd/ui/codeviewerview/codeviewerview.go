package codeviewerview

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type model struct {
	codeviewer  viewport.Model
	profilePath string
	viewChanger view.ViewChanger
}

func New(path string, windowSize tea.WindowSizeMsg, viewChanger view.ViewChanger) model {
	l.Logger.Info("Initializing code viewer")
	ws := windowSize
	vp := viewport.New(ws.Height, ws.Width)
	content, err := utils.LoadProfileContent(path)
	if err != nil {
		l.Logger.Error("Failed to load profile content", "error", err)
	}
	vp.SetContent(content)
	return model{
		codeviewer: viewport.Model{
			Width:             ws.Width - 10,
			Height:            ws.Height - 10,
			MouseWheelEnabled: true,
			Style: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true).
				Padding(1, 1).Align(lipgloss.Center).
				Background(lipgloss.Color("#8396a6")),
		},
		profilePath: path,
		viewChanger: viewChanger,
	}

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.codeviewer.Height = msg.Height - 10
		m.codeviewer.Width = msg.Width - 10
	case tea.KeyMsg:
		switch msg.String() {

		}
	}
	var cmd tea.Cmd
	m.codeviewer, cmd = m.codeviewer.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.codeviewer.View()
}
