package parentview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/mainview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/profileselector"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/shellview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type ParentModel struct {
	currentView tea.Model
	screens     map[string]tea.Model
}

func New() ParentModel {
	return ParentModel{
		currentView: mainview.New(),
		screens: map[string]tea.Model{
			"mainView":    mainview.New(),
			"profileView": profileselector.New(),
			"shellView":   shellview.New(),
		},
	}
}

func (m ParentModel) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m ParentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case utils.SwitchViewMsg:
		if nextView, ok := m.screens[string(msg)]; ok {
			m.currentView = nextView
		}
	default:
		var cmd tea.Cmd
		m.currentView, cmd = m.currentView.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m ParentModel) View() string {
	return m.currentView.View()
}
