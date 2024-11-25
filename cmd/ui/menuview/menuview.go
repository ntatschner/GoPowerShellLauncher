package menuview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

type menuItem struct {
	title       string
	description string
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Description() string { return m.description }
func (m menuItem) FilterValue() string { return m.title }

type model struct {
	menuList list.Model
}

func New() model {
	items := []list.Item{
		menuItem{title: "Select Profiles", description: "PowerShell profile selection screen."},
		menuItem{title: "Create Shortcuts", description: "Shortcut creation screen."},
		menuItem{title: "Exit", description: "Exit the application."},
	}

	list := list.New(items, list.NewDefaultDelegate(), 20, 10)
	list.Title = "Main Menu"
	list.SetFilteringEnabled(false)

	return model{menuList: list}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.menuList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.menuList, cmd = m.menuList.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.menuList.View()
}
