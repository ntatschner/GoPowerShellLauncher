package profileselector

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

func New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) *model {
	l.Logger.Debug("Initializing profile list")
	path, err := os.Getwd()
	l.Logger.Debug("Getting working directory", "path", path)
	if err != nil {
		l.Logger.Error("Failed to get working directory", "error", err)
	}
	path = path + string(os.PathSeparator) + "config.json"
	loadConfig, err := utils.LoadConfig()
	if err != nil {
		l.Logger.Error("Failed to load configuration file", "error", err)
	} else {
		l.Logger.Info("Loaded configuration file", "config", loadConfig)
	}

	profiles, err := utils.LoadProfilesFromDir()
	if err != nil {
		l.Logger.Error("Failed to load profiles", "error", err)
	}

	var items []list.Item
	for _, p := range profiles {
		item := types.ProfileItem{
			ItemTitle:       p.ItemTitle,
			ItemDescription: p.ItemDescription,
			IsValid:         p.IsValid,
			Path:            p.Path,
			Shell:           p.Shell,
			Name:            p.Name,
		}
		items = append(items, item)
	}
	delegateKeyMap, err := styles.NewDelegateKeyMap()
	if err != nil {
		l.Logger.Fatal("Failed to create delegate key map", "error", err)
		return nil
	}
	itemDelegate, delerr := styles.NewItemDelegate(delegateKeyMap)
	if delerr != nil {
		l.Logger.Fatal("Failed to create item delegate", "error", delerr)
		return nil
	}
	profilesList := list.New(items, itemDelegate, windowSize.Width+50, windowSize.Height)
	profilesList.Title = "Available PowerShell Profiles"

	profilesList.Styles.Title = styles.TitleStyle
	profilesList.Styles.PaginationStyle = styles.PaginationStyle

	profilesList.SetFilteringEnabled(true)
	profilesList.FilterValue()
	profilesList.SetShowStatusBar(true)
	profilesList.SetShowTitle(true)

	return &model{
		profilesList: profilesList,
		selected:     make(map[int]struct{}),
		viewChanger:  viewChanger,
		windowSize:   windowSize,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.SetWindowTitle("Profile Selection")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.profilesList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			items := m.profilesList.Items()
			i := m.profilesList.Index()
			if i < 0 || i >= len(m.profilesList.Items()) {
				l.Logger.Error("Invalid index", "index", i)
				break
			}
			item := items[i].(types.ProfileItem)
			if !item.IsValid {
				l.Logger.Warn("Selected item is not valid")
			} else {
				if _, ok := m.selected[i]; ok {
					delete(m.selected, i)
					l.Logger.Debug("Deselected profile", "index", i)
					cmd = tea.Batch(func() tea.Msg {
						return styles.StatusBarUpdate(false)
					})
					item.IsSelected = false
					items[i] = item
					return m, cmd
				} else {
					m.selected[i] = struct{}{}
					l.Logger.Debug("Selected profile", "index", i)
					cmd = tea.Batch(func() tea.Msg {
						return styles.StatusBarUpdate(true)
					})
					item.IsSelected = true
					items[i] = item
					return m, cmd
				}
			}
		case "enter":
			// load selected profiles
			var selectedProfiles []types.ProfileItem
			if len(m.selected) == 0 {
				l.Logger.Warn("No Profiles selected, using currently highlighted profile")
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
			return m, m.viewChanger.ChangeView(shellview.New(selectedProfiles, m.windowSize, m.viewChanger, false), true)
		case "v":
			// view profile content
			i := m.profilesList.Index()
			if i < 0 || i >= len(m.profilesList.Items()) {
				l.Logger.Error("Invalid index", "index", i)
				break
			}
			item := m.profilesList.Items()[i].(types.ProfileItem)
			return m, m.viewChanger.ChangeView(codeviewerview.New(item.Path, m.windowSize, m.viewChanger), true)
		}
	}

	m.profilesList, cmd = m.profilesList.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return m.profilesList.View()
}

func (m *model) ClearSelectedItems() {
	m.selected = make(map[int]struct{})
}

func (m *model) FilterState() list.FilterState {
	return m.profilesList.FilterState()
}

// Ensure model implements view.Clearable
var _ view.Clearable = (*model)(nil)
