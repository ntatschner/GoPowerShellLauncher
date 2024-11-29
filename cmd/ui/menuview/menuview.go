package menuview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/profileselector"
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
	l.Logger.Info("Initializing main menu")
	items := []list.Item{
		menuItem{title: "Select Profiles", description: "PowerShell profile selection screen.", pageName: "profilesView"},
		menuItem{title: "Create Shortcuts", description: "Shortcut creation screen.", pageName: "shortcutsView"},
		menuItem{title: "Exit", description: "Exit the application.", pageName: "exit"},
	}

	list := list.New(items, list.NewDefaultDelegate(), windowSize.Width, windowSize.Height)
	list.Title = "Main Menu"
	list.SetFilteringEnabled(false)

	return &model{
		menuList:    list,
		viewChanger: viewChanger,
		windowSize:  windowSize,
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
		m.menuList.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i := m.menuList.Index()
			item := m.menuList.Items()[i].(menuItem)
			switch item.PageName() {
			case "profilesView":
				l.Logger.Info("Changing view to profile selector")
				return m, m.viewChanger.ChangeView(profileselector.New(m.viewChanger, m.windowSize))
			case "shortcutsView":
				l.Logger.Info("Changing view to shortcut selector")
				// m.viewChanger.ChangeView(shortcutselector.New(m.viewChanger))
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
