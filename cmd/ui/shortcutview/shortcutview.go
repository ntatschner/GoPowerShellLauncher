package shortcutview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/shellview"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/styles"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

func init() {
	ui.RegisterPage("shortcutView", Page{})
}

type Page struct{}

func (p Page) New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) tea.Model {
	return New(viewChanger, windowSize)
}

func (p Page) Title() string {
	return "Create Shortcuts"
}

func (p Page) Description() string {
	return "Shortcut creation screen."
}

type model struct {
	shortcutsList list.Model
	windowSize    tea.WindowSizeMsg
	viewChanger   view.ViewChanger
}

func New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) *model {
	l.Logger.Debug("Initializing shortcut list")
	config, err := utils.LoadConfig()
	if err != nil {
		l.Logger.Error("Failed to load configuration file", "error", err)
	}

	var items []list.Item
	for _, s := range config.Shortcuts {
		items = append(items, s)
	}

	delegate := list.NewDefaultDelegate()
	shortcutsList := list.New(items, delegate, windowSize.Width, windowSize.Height)
	shortcutsList.Title = "Available Shortcuts"
	shortcutsList.Styles.Title = styles.TitleStyle

	return &model{
		shortcutsList: shortcutsList,
		viewChanger:   viewChanger,
		windowSize:    windowSize,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.SetWindowTitle("Shortcut Selection")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.shortcutsList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i := m.shortcutsList.Index()
			item := m.shortcutsList.Items()[i].(utils.Shortcut)
			var profilePaths []string
			for _, p := range item.Profiles {
				profilePaths = append(profilePaths, p.Path)
			}
			// ToDo: Get shell from shortcut
			return m, m.viewChanger.ChangeView(shellview.New(nil, m.windowSize, m.viewChanger, true), true)
		}
	}

	m.shortcutsList, cmd = m.shortcutsList.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return m.shortcutsList.View()
}
