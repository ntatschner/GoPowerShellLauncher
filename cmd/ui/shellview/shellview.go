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
	// Load shell items based on profiles
	var items []list.Item
	for _, shell := range shells {
		// get the profiles that use this shell
		var profilesForShell []string
		for _, profile := range profiles {
			for _, shortName := range shell.ShortName() {
				if profile.Shell == shortName {
					profilesForShell = append(profilesForShell, profile.Path)
					break
				}
			}
		}
		shellItem := types.ShellItem{
			ItemTitle:       shell.Name(),
			ItemDescription: shell.Description() + ": loaded profiles: " + strconv.Itoa(len(profilesForShell)),
			ShortName:       shell.ShortName(),
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
	l.Logger.Info("Update called", "msg", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.shellsList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Launch the selected shell with the profiles that use it
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
