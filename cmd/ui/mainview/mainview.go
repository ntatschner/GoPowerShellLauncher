package mainview

import (
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/menuview"
)

type mainModel struct {
	currentView  tea.Model
	previousView tea.Model
	windowSize   tea.WindowSizeMsg
}

func NewMainModel() *mainModel {
	l.Logger.Info("Creating a new main view")
	mainModel := &mainModel{
		windowSize: tea.WindowSizeMsg{},
	}
	mainView := menuview.New(mainModel)
	mainModel.currentView = mainView
	return mainModel
}

func (m *mainModel) Init() tea.Cmd {
	return m.currentView.Init()
}

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		// Pass the window size to the current view
		m.currentView, _ = m.currentView.Update(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "backspace":
			if m.previousView != nil {
				m.currentView = m.previousView
				m.previousView = nil
			}
		}
	case ChangeViewMsg:
		return m.handleChangeViewMsg(msg)
	}

	var cmd tea.Cmd
	m.currentView, cmd = m.currentView.Update(msg)
	return m, cmd
}

func (m *mainModel) View() string {
	return m.currentView.View()
}

type ChangeViewMsg struct {
	NewView tea.Model
}

func (m *mainModel) handleChangeViewMsg(msg ChangeViewMsg) (tea.Model, tea.Cmd) {
	m.previousView = m.currentView
	m.currentView = msg.NewView
	return m, nil
}

func (m *mainModel) ChangeView(newView tea.Model) {
	m.previousView = m.currentView
	m.currentView = newView
	tea.NewProgram(m).Send(ChangeViewMsg{NewView: newView})
}
