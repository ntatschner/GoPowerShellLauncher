package shellview

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type model struct {
	shellsList     list.Model
	selected       map[int]struct{}
	windowSize     tea.WindowSizeMsg
	viewChanger    view.ViewChanger
	loadedProfiles []types.ProfileItem
}

func New(profiles []types.ProfileItem, windowSize tea.WindowSizeMsg, viewChanger view.ViewChanger) *model {
	l.Logger.Info("Initializing shell list", "profiles", profiles)
	shells, err := utils.LoadShells()
	if err != nil {
		l.Logger.Error("Failed to load shells", "error", err)
	}
	if len(shells) == 0 {
		l.Logger.Error("No shells loaded")
		return &model{
			shellsList:     list.New([]list.Item{}, list.NewDefaultDelegate(), windowSize.Width, windowSize.Height),
			selected:       make(map[int]struct{}),
			windowSize:     windowSize,
			viewChanger:    viewChanger,
			loadedProfiles: profiles,
		}
	}
	// Load shell items based on profiles
	var items []list.Item
	for _, shell := range shells {
		l.Logger.Info("Processing shell", "shell", shell.ItemTitle)
		// get the profiles that use this shell
		var profilesForShell []string
		for _, profile := range profiles {
			l.Logger.Info("Processing profile", "profile", profile.ItemTitle)
			for _, shortName := range shell.ShortName {
				l.Logger.Info("Processing short name", "shortName", shortName)
				profileShell := utils.NormalizeString(profile.Shell)
				shortNameTrimmed := utils.NormalizeString(shortName)
				l.Logger.Info("Comparing", "profile.Shell", profileShell, "shortName", shortNameTrimmed)
				if profileShell == shortNameTrimmed {
					l.Logger.Info("Profile uses shell", "profile", profile.ItemTitle, "shell", shortNameTrimmed)
					profilesForShell = append(profilesForShell, profile.Path)
					break
				}
			}
		}
		l.Logger.Info("Profiles for shell", "shell", shell.ItemTitle, "profilesForShell", profilesForShell)
		shellItem := types.ShellItem{
			ItemTitle:       shell.ItemTitle,
			ItemDescription: shell.ItemDescription + ": loaded profiles: " + strconv.Itoa(len(profilesForShell)),
			Name:            shell.Name,
			ShortName:       shell.ShortName,
			Path:            shell.Path,
			ProfilePaths:    profilesForShell,
		}
		items = append(items, shellItem)
	}

	shellsList := list.New(items, list.NewDefaultDelegate(), windowSize.Width, windowSize.Height)
	shellsList.Title = "Available Shells"
	shellsList.SetFilteringEnabled(false)
	shellsList.SetShowStatusBar(true)
	return &model{
		shellsList:     shellsList,
		selected:       make(map[int]struct{}),
		windowSize:     windowSize,
		viewChanger:    viewChanger,
		loadedProfiles: profiles,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.shellsList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			return m, m.viewChanger.ChangeView(nil, false)
		}
	}

	var cmd tea.Cmd
	m.shellsList, cmd = m.shellsList.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return m.shellsList.View()
}

// CountProfilesMatchingShell counts the profiles that match the shell to the shortName of the shellItem
func (m *model) CountProfilesMatchingShell(shortName string) int {
	count := 0
	for _, profile := range m.loadedProfiles {
		if profile.Shell == shortName {
			count++
		}
	}
	return count
}
