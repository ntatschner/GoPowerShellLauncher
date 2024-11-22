package profileselector

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true)
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF99FF")).Italic(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF00FF")).Bold(true)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("\\u{1F449} " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type sessionState int

type profileItem struct {
	title       string
	description string
	cmd         tea.Cmd
}

func (m profileItem) FilterValue() string {
	return m.title
}

func (m profileItem) Title() string {
	return m.title
}

func (m profileItem) Description() string {
	return m.description
}

// MainView is the main view of the application
type model struct {
	profilesList list.Model
	selected     map[int]struct{}
	keymap       utils.KeyMap
}

func New() *model {
	return &model{}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		m.initList(150, 100)
		m.keymap = utils.DefaultKeyMap()
		return nil
	}
}

func (m *model) initList(width, height int) {
	l.Logger.Info("Initializing profile list")
	m.profilesList = list.New([]list.Item{}, itemDelegate{}, width, height)
	m.profilesList.Title = "Available Profiles"
	loadConfig, err := utils.LoadConfig("config.json")
	if err != nil {
		l.Logger.Error("Failed to load configuration file", "error", err)
	}
	profiles, err := utils.LoadProfiles(loadConfig.CsvPath)
	if err != nil {
		l.Logger.Error("Failed to load profiles", "error", err)
	}
	var items []list.Item
	for _, p := range profiles {
		items = append(items, p)
	}
	m.profilesList.SetItems(items)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		l.Logger.Info("Key pressed", "msg", msg.String())
		switch msg.String() {
		case "enter":
			//selectedItems := m.profilesList.SelectedItem()
		case " ":
			_, ok := m.selected[m.profilesList.Cursor()]
			if ok {
				delete(m.selected, m.profilesList.Cursor())
			} else {
				m.selected[m.profilesList.Cursor()] = struct{}{}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.profilesList.View()
}
