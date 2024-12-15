package menuview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/profileselector"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/shortcutview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/styles"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
)

type menuItem struct {
	title       string
	description string
	pageName    string
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Description() string { return m.description }
func (m menuItem) FilterValue() string { return m.title }
func (m menuItem) PageName() string    { return m.pageName }

type model struct {
	menuList    list.Model
	viewChanger view.ViewChanger
	windowSize  tea.WindowSizeMsg
}

func New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) *model {
	l.Logger.Debug("Initializing main menu")
	items := []list.Item{
		menuItem{title: "Select Profiles", description: "PowerShell profile selection screen.", pageName: "profilesView"},
		menuItem{title: "Create Shortcuts", description: "Shortcut creation screen.", pageName: "shortcutsView"},
		menuItem{title: "Exit", description: "Exit the application.", pageName: "exit"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#008A74", Dark: "#40C1AC"}).
		Padding(0, 0, 0, 2)
	delegate.Styles.NormalDesc = delegate.Styles.NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#FF94F4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#FF94F4", Dark: "#FF94F4"}).
		Padding(0, 0, 0, 1)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	list := list.New(items, delegate, windowSize.Width, windowSize.Height)
	list.Title = "Main Menu"
	list.Styles.Title = styles.TitleStyle
	list.Styles.HelpStyle = styles.HelpStyle
	list.SetFilteringEnabled(false)

	return &model{
		menuList:    list,
		viewChanger: viewChanger,
		windowSize:  windowSize,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.SetWindowTitle("PowerShell Profile Launcher")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.menuList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i := m.menuList.Index()
			item := m.menuList.Items()[i].(menuItem)
			switch item.PageName() {
			case "profilesView":
				l.Logger.Debug("Changing view to profile selector")
				return m, m.viewChanger.ChangeView(profileselector.New(m.viewChanger, m.windowSize), true)
			case "shortcutsView":
				l.Logger.Debug("Changing view to shortcut selector")
				return m, m.viewChanger.ChangeView(shortcutview.New(m.viewChanger, m.windowSize), true)
			case "exit":
				l.Logger.Info("Exiting application")
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.menuList, cmd = m.menuList.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return m.menuList.View()
}
