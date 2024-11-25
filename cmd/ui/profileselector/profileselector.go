package profileselector

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF")).Italic(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF00FF")).Bold(true)
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
	items []table.Row
}

func (m model) Init() tea.Cmd { return nil }

func InitTable() {
	l.Logger.Info("Initializing profile list")

	loadConfig, err := utils.LoadConfig("config.json")
	if err != nil {
		l.Logger.Error("Failed to load configuration file", "error", err)
	}
	profiles, err := utils.LoadProfiles(loadConfig.CsvPath)
	if err != nil {
		l.Logger.Error("Failed to load profiles", "error", err)
	}
	var rows []table.Row
	for _, p := range profiles {
		items = append(rows, p)
	}
	
}

type profileItem struct {
	title       string
	description string
}

func (p profileItem) Title() string       { return p.title }
func (p profileItem) Description() string { return p.description }
func (p profileItem) FilterValue() string { return p.title }

func New() model {


	return model{table: }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetHeight(msg.Height)
		m.table.SetWidth(msg.Width)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.table.View()
}
