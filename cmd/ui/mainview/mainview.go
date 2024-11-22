package mainview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

var (
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
	subtleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF")).Italic(true)
)

type menuItem struct {
	title       string
	description string
	screen      utils.SwitchViewMsg
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
	keymap   utils.KeyMap
}

func New() *model {
	return &model{}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		m.initList(150, 100)
		m.keymap = utils.DefaultKeyMap()
		return nil
	}
}

func (m *model) initList(width, height int) {
	m.menuList = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.menuList.Title = "Main Menu"
	m.menuList.SetFilteringEnabled(false)
	menuItems := []list.Item{
		menuItem{title: "Select Profiles", description: "PowerShell profile selection screen.", screen: "profileView"},
		menuItem{title: "Create Shortcuts", description: "Shortcut creation screen.", screen: "shortcutView"},
	}

	var items []list.Item
	items = append(items, menuItems...)
	m.menuList.SetItems(items)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			selectedItem := m.menuList.SelectedItem()
			if menuItem, ok := selectedItem.(menuItem); ok {
				l.Logger.Info("Switching Screen", "Screen", m.menuList.SelectedItem())
				return m, func() tea.Msg {
					return utils.SwitchViewMsg(menuItem.screen)
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.menuList.View()
}
