package profileselector

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/codeviewerview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/shellview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#2f80c2")).Bold(true).Align(lipgloss.Center).Underline(true)
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF")).Italic(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Background(lipgloss.Color("#FF99FF"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF00FF")).Bold(true)
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	profilesList list.Model
	selected     map[int]struct{}
	csvPath      string
	windowSize   tea.WindowSizeMsg
	viewChanger  view.ViewChanger
}

func New(viewChanger view.ViewChanger) *model {
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
		if p.IsValid {
			valid = "✅"
		}
		item := types.ProfileItem{
			Title:       p.Title,
			Description: fmt.Sprintf("%s %s %s %s %s", p.Description, valid, p.Path, p.Hash, p.Shell),
			Valid:       valid,
			IsValid:     p.IsValid,
			Path:        p.Path,
			Hash:        p.Hash,
			Shell:       p.Shell,
		}
		items = append(items, item)
	}

	profilesList := list.New(items, list.NewDefaultDelegate(), 50, 30)
	profilesList.Title = "Available Profiles"
	profilesList.SetFilteringEnabled(true)
	profilesList.SetShowStatusBar(true)
	profilesList.Styles.Title = titleStyle

	return &model{
		profilesList: profilesList,
		selected:     make(map[int]struct{}),
		csvPath:      loadConfig.CsvPath,
		viewChanger:  viewChanger,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l.Logger.Info("Update called", "msg", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.profilesList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			i := m.profilesList.Index()
			if i < 0 || i >= len(m.profilesList.Items()) {
				l.Logger.Error("Invalid index", "index", i)
				break
			}
			item := m.profilesList.Items()[i].(types.ProfileItem)
			if !item.IsValid {
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
		case "enter":
			// load selected profiles
			var selectedProfiles []types.ProfileItem
			if len(m.selected) == 0 {
				l.Logger.Info("No Profiles selected, using currently highlighted profile")
				i := m.profilesList.Index()
				if i < 0 || i >= len(m.profilesList.Items()) {
					l.Logger.Error("Invalid index", "index", i)
					break
				}
				item := m.profilesList.Items()[i].(types.ProfileItem)
				selectedProfiles = append(selectedProfiles, item)
			}
			for i := range m.selected {
				item := m.profilesList.Items()[i].(types.ProfileItem)
				selectedProfiles = append(selectedProfiles, item)
			}
			// open shellview with profiles selected
			l.Logger.Info("Selected profiles", "profiles", selectedProfiles)
			return m, m.viewChanger.ChangeView(shellview.New(selectedProfiles, m.windowSize, m.viewChanger))
		case "v":
			// view profile content
			i := m.profilesList.Index()
			if i < 0 || i >= len(m.profilesList.Items()) {
				l.Logger.Error("Invalid index", "index", i)
				break
			}
			item := m.profilesList.Items()[i].(types.ProfileItem)
			return m, m.viewChanger.ChangeView(codeviewerview.New(item.Path, m.windowSize, m.viewChanger))
		}
	}

	var cmd tea.Cmd
	m.profilesList, cmd = m.profilesList.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return m.profilesList.View()
}
