package mainview

import (
	"github.com/charmbracelet/bubbles/list"
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

type menuItem struct {
	title       string
	description string
	cmd         tea.Cmd
}

func (m menuItem) FilterValue() string {
	return m.title
}

func (m menuItem) Title() string {
	return m.title
}

func (m menuItem) Description() string {
	return m.description
}

// MainView is the main view of the application
type model struct {
	state    sessionState
	profiles tea.Model
	shells   tea.Model
	menuList list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func InitMainView() (tea.Model, tea.Cmd) {
	l.Logger.Info("Initializing MainView...")
	list := list.New([]list.Item{}, list.NewDefaultDelegate(), 150, 20)

}

func (m *model) initList(width, height int) {
	m.menuList = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.menuList.Title = "Main Menu"
	m.menuList.SetItems([]list.Item{
		menuItem{title: "Select Profiles", description: "", cmd: nil},
		menuItem{title: "Create Shortcuts", description: "", cmd: nil},
		menuItem{title: "Exit", description: "", cmd: tea.Quit},
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	var cmd tea.Cmd
	m.menuList, cmd = m.menuList.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("GoPowerShellLauncher"),
		subtleStyle.Render("Press q to quit"),
	)
}
