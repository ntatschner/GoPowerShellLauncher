package shellview

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/launcher"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/shortcutconfigview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/styles"
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
			for _, shortName := range shell.ShortNames {
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
			ShortNames:      shell.ShortNames,
			Path:            shell.Path,
			ProfilePaths:    profilesForShell,
		}
		items = append(items, shellItem)
	}
	delegateKeyMap, err := styles.NewShellDelegateKeyMap()
	if err != nil {
		l.Logger.Fatal("Failed to create delegate key map", "error", err)
		return nil
	}
	itemDelegate, delerr := styles.NewShellItemDelegate(delegateKeyMap)
	if delerr != nil {
		l.Logger.Fatal("Failed to create item delegate", "error", delerr)
		return nil
	}
	shellsList := list.New(items, itemDelegate, windowSize.Width, windowSize.Height)
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.shellsList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			items := m.shellsList.Items()
			i := m.shellsList.Index()
			if i < 0 || i >= len(m.shellsList.Items()) {
				l.Logger.Error("Invalid index", "index", i)
				break
			}
			item := items[i].(types.ShellItem)
			if _, ok := m.selected[i]; ok {
				delete(m.selected, i)
				l.Logger.Debug("Deselected shell", "index", i)
				cmd = tea.Batch(func() tea.Msg {
					return styles.StatusBarUpdate(false)
				})
				item.IsSelected = false
				items[i] = item
				return m, cmd
			} else {
				m.selected[i] = struct{}{}
				l.Logger.Debug("Selected shelli", "index", i)
				cmd = tea.Batch(func() tea.Msg {
					return styles.StatusBarUpdate(true)
				})
				item.IsSelected = true
				items[i] = item
				return m, cmd
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
			var selectedShells []types.ShellItem
			for i := range m.selected {
				item := m.shellsList.Items()[i].(types.ShellItem)
				selectedShells = append(selectedShells, item)
			}
			if m.shortcut {
				return m, m.viewChanger.ChangeView(shortcutconfigview.New(m.viewChanger, m.windowSize, m.loadedProfiles, selectedShells), false)
			} else {
				l.Logger.Info("Launching selected shells", "selected", m.selected, "profiles", m.loadedProfiles)
				for i := range m.selected {
					merged := utils.MergeSelectedProfiles(m.shellsList.Items()[i].(types.ShellItem).ProfilePaths)
					item := m.shellsList.Items()[i].(types.ShellItem)
					err := launcher.ExecutePowerShellProcess(merged, item.Path)
					if err != nil {
						l.Logger.Error("Failed to execute PowerShell process", "Error", err)
					}
				}
			}
		}
	}

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
