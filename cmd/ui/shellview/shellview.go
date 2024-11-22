package shellview

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

type sessionState int

type shellItem struct {
	title       string
	description string
	cmd         tea.Cmd
}

type model struct {
	shellsList list.Model
}

func New() *model {
	return &model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m shellItem) FilterValue() string {
	return m.title
}

func (m shellItem) Title() string {
	return m.title
}

func (m shellItem) Description() string {
	return m.description
}

func (m *model) initList(width, height int) {
	l.Logger.Info("Initializing shell list")
	m.shellsList = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.shellsList.Title = "Available Shells"
	loadShellItems, err := utils.LoadShells()
	if err != nil {
		l.Logger.Error("Failed to load shells", "error", err)
	}
	shells := []list.Item{}
	for s := range loadShellItems {
		shells = append(shells, loadShellItems[s])
	}
	m.shellsList.SetItems(shells)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	return m.shellsList.View()
}
