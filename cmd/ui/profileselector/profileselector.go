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
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/styles"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type model struct {
	profilesList list.Model
	selected     map[int]struct{}
	csvPath      string
	windowSize   tea.WindowSizeMsg
	viewChanger  view.ViewChanger
}

func New(viewChanger view.ViewChanger, windowsSize tea.WindowSizeMsg) *model {
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
		msg := "Validated Import:"
		valid := msg + "❌"
		if p.IsValid {
			valid = msg + "✅"
		}
		item := types.ProfileItem{
			ItemTitle:       p.ItemTitle,
			ItemDescription: fmt.Sprintf("%s\n%s\n%s\n%s\n%s", p.GetDescription(), valid, p.GetPath(), p.GetHash(), p.GetShell()),
			Valid:           valid,
			IsValid:         p.IsValid,
			Path:            p.Path,
			Hash:            p.Hash,
			Shell:           p.Shell,
		}
		items = append(items, item)
	}

	profilesList := list.New(items, list.NewDefaultDelegate(), windowsSize.Width, windowsSize.Height)
	profilesList.Title = "Available Profiles"
	profilesList.Styles.Title = styles.TitleStyle
	profilesList.Styles.PaginationStyle = styles.PaginationStyle
	profilesList.Styles.Title.Align(lipgloss.Center)
	profilesList.Styles.HelpStyle = styles.HelpStyle
	profilesList.SetFilteringEnabled(true)
	profilesList.SetShowStatusBar(true)

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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		h, v := styles.AppStyle.GetFrameSize()
		m.profilesList.SetSize(msg.Width-h, msg.Height-v)
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
