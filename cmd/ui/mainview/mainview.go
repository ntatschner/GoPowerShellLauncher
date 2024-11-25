package mainview

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/menuview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/profileselector"
)

type sessionState int

const (
	menuView sessionState = iota
	profilesView
)

type mainModel struct {
	state        sessionState
	mainView     tea.Model
	profilesView tea.Model
	currentView  tea.Model
}

func NewMainModel() mainModel {
	mainView := menuview.New()
	profilesView := profileselector.New()
	return mainModel{
		state:        menuView,
		mainView:     mainView,
		profilesView: profilesView,
		currentView:  mainView,
	}
}

func (m mainModel) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "enter":
			if m.state == menuView {
				m.state = profilesView
				m.currentView = m.profilesView
			} else {
				m.state = menuView
				m.currentView = m.mainView
			}
		}
	}

	var cmd tea.Cmd
	m.currentView, cmd = m.currentView.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	return m.currentView.View()
}
