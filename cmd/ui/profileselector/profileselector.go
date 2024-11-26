package profileselector

import (
	"os"

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
	table      table.Model
	items      []table.Row
	configPath string
}

var configPath string

func (m model) ConfigPath() string { return m.configPath }

func (m model) Init() tea.Cmd {
	l.Logger.Info("Initializing profile list")
	path, err := os.Getwd()
	l.Logger.Info("Getting working directory", "path", path)
	if err != nil {
		l.Logger.Error("Failed to get working directory", "error", err)
	}
	path = path + string(os.PathSeparator) + "config"
	loadConfig, err := utils.LoadConfig(path)
	if err != nil {
		l.Logger.Error("Failed to load configuration file", "error", err)
	} else {
		l.Logger.Info("Loaded configuration file", "config", loadConfig)
		configPath = loadConfig.CsvPath
	}
	return nil
}

func New() model {
	var valid string
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Valid", Width: 20},
		{Title: "Description", Width: 20},
		{Title: "Path", Width: 20},
		{Title: "Hash", Width: 20},
		{Title: "Shell", Width: 20},
	}
	profiles, err := utils.LoadProfiles(configPath)
	if err != nil {
		l.Logger.Error("Failed to load profiles", "error", err)
	}
	var items []table.Row
	for _, p := range profiles {
		if p.Valid() {
			valid = "✅"
		} else {
			valid = "❌"
		}
		items = append(items, table.Row{
			p.Name(),
			valid,
			p.Description(),
			p.Path(),
			p.Hash(),
			p.Shell(),
		})
	}
	table := table.New()
	table.Focus()
	table.SetColumns(columns)
	table.SetRows(items)

	return model{table: table}
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
