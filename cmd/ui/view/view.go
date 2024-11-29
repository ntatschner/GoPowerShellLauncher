package view

import tea "github.com/charmbracelet/bubbletea"

type ViewChanger interface {
	ChangeView(newView tea.Model) tea.Cmd
}
