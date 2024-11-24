package parentview

import (
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/menuview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/common"
)

type ParentModel struct {
	state        common.SessionState
	menuView     tea.Model
	profilesView tea.Model
	previousView tea.Model
}

func New() ParentModel {
	return menuview.New()
}

func (m ParentModel) Init() tea.Cmd {
	return nil
}

func (m ParentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		common.WindowSize = msg
	case tea.KeyMsg:
		m.state = menuView
	}
	switch m.state {
	case menuView:
		menuModel, ok := menuview.New()
		if !ok {
			l.Logger.Error("Failed to load view")
		} else {
			return menuModel, nil
		}
	}
	return m, nil
}

func (m ParentModel) View() string {
	modelView := m.state[m.state].View()
	return modelView
}
