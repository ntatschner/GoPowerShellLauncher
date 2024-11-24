package mainview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/common"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/profileselector"
)

var (
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF")).Italic(true)
)

type menuItem struct {
	title       string
	description string
	screen      string
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
	menuList list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) New() (tea.Model, tea.Cmd) {
	return m.initList()
}

func (m *model) initList() (tea.Model, tea.Cmd) {
	m.menuList = list.New([]list.Item{}, list.NewDefaultDelegate(), common.WindowSize.Height, common.WindowSize.Width)
	m.menuList.Title = "Main Menu"
	m.menuList.SetFilteringEnabled(false)
	menuItems := []list.Item{
		menuItem{title: "Select Profiles", description: "PowerShell profile selection screen.", screen: "profilesView"},
		menuItem{title: "Create Shortcuts", description: "Shortcut creation screen.", screen: "profilesView"},
	}

	var items []list.Item
	items = append(items, menuItems...)
	m.menuList.SetItems(items)
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		common.WindowSize = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			newItem := m.menuList.SelectedItem().(menuItem).screen
			if newItem == "profilesView" {
				newModel := profileselector.New()
				return newModel, nil
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.menuList.View()
}
