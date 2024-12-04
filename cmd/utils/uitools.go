package utils

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusBar Message Update

type StatusBarUpdate list.Model

func StatusBarMessage(m list.Model, msg string, style lipgloss.Style) tea.Cmd {
	msg = style.Render(msg)
	m.NewStatusMessage(msg)
	return func() tea.Msg {
		return StatusBarUpdate(m)
	}
}
