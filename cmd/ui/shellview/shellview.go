package shellview

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/launcher"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/shortcutconfigview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type model struct {
	shellsList     list.Model
	selected       map[int]struct{}
	windowSize     tea.WindowSizeMsg
	viewChanger    view.ViewChanger
	loadedProfiles []types.ProfileItem
	shortcut       bool
}

func New(profiles []types.ProfileItem, windowSize tea.WindowSizeMsg, viewChanger view.ViewChanger, createShortcut bool) *model {
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
			shortcut:       createShortcut,
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
		shortcut:       createShortcut,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.SetWindowTitle("Select Shell")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.shellsList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			i := m.shellsList.Index()
			if i < 0 || i >= len(m.shellsList.Items()) {
				l.Logger.Error("Invalid index", "index", i)
				break
			}
			if _, ok := m.selected[i]; ok {
				delete(m.selected, i)
				l.Logger.Info("Deselected shell", "index", i)
			} else {
				m.selected[i] = struct{}{}
				l.Logger.Info("Selected shell", "index", i)
			}

		case "enter":
			if len(m.selected) == 0 {
				l.Logger.Warn("No shell selected, using currently highlighted.")
				i := m.shellsList.Index()
				if i < 0 || i >= len(m.shellsList.Items()) {
					l.Logger.Error("Invalid index", "index", i)
					break
				}
				m.selected[i] = struct{}{}
			}
			if m.shortcut {
				i := m.shellsList.Index()
				item := m.shellsList.Items()[i].(types.ShellItem)
				return m, m.viewChanger.ChangeView(shortcutconfigview.New(m.viewChanger, m.windowSize, m.loadedProfiles, item), false)
				// l.Logger.Info("Creating shortcut", "selected", m.selected, "profiles", m.loadedProfiles)
				// //
				// err := utils.CreateShortcut(m.shellsList.Items()[0].(types.ShellItem).ProfilePaths, "test", "c:\\nerd_stuff")
				// if err != nil {
				// 	l.Logger.Error("Failed to create shortcut", "Error", err)
				// }

			} else {
				l.Logger.Info("Launching selected shells", "selected", m.selected, "profiles", m.loadedProfiles)
				for i := range m.selected {
					merged := launcher.MergeSelectedProfiles(m.shellsList.Items()[i].(types.ShellItem).ProfilePaths)
					// tempFilePath, err := launcher.CreateTempFile(merged)
					encodedCommand, err := utils.EncodeCommand(merged)
					if err != nil {
						l.Logger.Error("Failed to encode profiles", "Error", err)
						continue
					}
					item := m.shellsList.Items()[i].(types.ShellItem)
					err = launcher.ExecutePowerShellProcess(encodedCommand, item.Path)
					if err != nil {
						l.Logger.Error("Failed to execute PowerShell process", "Error", err)
					}
				}
			}
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
