package parentview

import (
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/common"
)

type sessionState int

const (
	menuView sessionState = iota
	profilesView
	conformationView
)

type ParentModel struct {
	state        sessionState
	mainView     tea.Model
	profilesView tea.Model
	previousView tea.Model
}

func (m ParentModel) New() tea.Model {
	return m.mainView
}

func (m ParentModel) Init() tea.Cmd {
	return nil
}

func (m ParentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		common.WindowSize = msg
	case tea.KeyMsg:
		m.state = menuView
	}
	switch m.state {
	case menuView:
		menuModel := m.mainView.Init()
		if menuModel == nil {
			l.Logger.Error("Failed to load view")
		} else {
			return m.mainView.Update(menuModel)
		}
	}
	return m, nil
}

func (m ParentModel) View() string {
	return m.mainView.View()
}
