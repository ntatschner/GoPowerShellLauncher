package profileselector

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF")).Italic(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Background(lipgloss.Color("#FF99FF"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF00FF")).Bold(true)
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type profileItem struct {
	title       string
	description string
	valid       string
	isValid     bool
	path        string
	hash        string
	shell       string
}

func (p profileItem) Title() string       { return p.title }
func (p profileItem) Description() string { return p.description }
func (p profileItem) FilterValue() string { return p.title }

type model struct {
	profilesList list.Model
	selected     map[int]struct{}
	csvPath      string
}

var configPath string

func (m model) CsvPath() string { return m.csvPath }

func New() model {
	l.Logger.Info("Initializing profile list")
	path, err := os.Getwd()
	l.Logger.Info("Getting working directory", "path", path)
	if err != nil {
		l.Logger.Error("Failed to get working directory", "error", err)
	}
	path = path + string(os.PathSeparator) + "config.json"
	loadConfig, err := utils.LoadConfig(path)
	if err != nil {
		l.Logger.Error("Failed to load configuration file", "error", err)
	} else {
		l.Logger.Info("Loaded configuration file", "config", loadConfig)
	}

	profiles, err := utils.LoadProfiles(loadConfig.CsvPath)
	if err != nil {
		l.Logger.Error("Failed to load profiles", "error", err)
	}

	var items []list.Item
	for _, p := range profiles {
		valid := "❌"
		if p.Valid() {
			valid = "✅"
		}
		item := profileItem{
			title: p.Name(),
			description: fmt.Sprintf("%s %s %s %s %s",
				p.Description(), valid, p.Path(), p.Hash(), p.Shell()),
			valid:   valid,
			isValid: p.Valid(),
			path:    p.Path(),
			hash:    p.Hash(),
			shell:   p.Shell(),
		}
		items = append(items, item)
	}

	profilesList := list.New(items, list.NewDefaultDelegate(), 50, 30)
	profilesList.Title = "Available Profiles"
	profilesList.SetFilteringEnabled(true)
	profilesList.StatusBarItemName()
	profilesList.SetShowStatusBar(true)
	profilesList.Styles.Title = titleStyle
	profilesList.Styles.StatusBar = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)

	return model{
		profilesList: profilesList,
		selected:     make(map[int]struct{}),
		csvPath:      loadConfig.CsvPath,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.profilesList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			i := m.profilesList.Index()
			item := m.profilesList.Items()[i].(profileItem)
			if !item.isValid {
				l.Logger.Warn("Selected item is not valid")
			} else {
				if _, ok := m.selected[i]; ok {
					delete(m.selected, i)
					l.Logger.Info("Deselected profile", "index", i)
				} else {
					m.selected[i] = struct{}{}
					l.Logger.Info("Selected profile", "index", i)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.profilesList, cmd = m.profilesList.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.profilesList.View()
}
