package mainview

import (
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
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
	windowSize   tea.WindowSizeMsg
}

func NewMainModel() mainModel {
	l.Logger.Info("Creating a new main view")
	mainView := menuview.New()
	profilesView := profileselector.New()
	return mainModel{
		state:        menuView,
		mainView:     mainView,
		profilesView: profilesView,
		currentView:  mainView,
		windowSize:   tea.WindowSizeMsg{},
	}
}

func (m mainModel) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		// Pass the window size to the current view if needed
		m.currentView, _ = m.currentView.Update(msg)
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
