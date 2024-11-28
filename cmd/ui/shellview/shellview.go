package shellview

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type shellItem struct {
	title       string
	description string
}

func (s shellItem) Title() string       { return s.title }
func (s shellItem) Description() string { return s.description }
func (s shellItem) FilterValue() string { return s.title }

type model struct {
	shellsList  list.Model
	selected    map[int]struct{}
	windowSize  tea.WindowSizeMsg
	viewChanger view.ViewChanger
}

func New(profiles []string, windowsSize tea.WindowSizeMsg, viewChanger view.ViewChanger) *model {
	l.Logger.Info("Initializing shell list")
	ws := windowsSize
	shellsList := list.New([]list.Item{}, list.NewDefaultDelegate(), ws.Width, ws.Height)
	shellsList.Title = "Available Shells"

	loadShellItems, err := utils.LoadShells()
	if err != nil {
		l.Logger.Error("Failed to load shells", "error", err)
	}

	var items []list.Item
	for _, s := range loadShellItems {
		items = append(items, shellItem{
			title:       s.Title(),
			description: s.Description(),
		})
	}
	shellsList.SetItems(items)
	return &model{
		shellsList:  shellsList,
		selected:    make(map[int]struct{}),
		viewChanger: viewChanger,
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
		}
	}

	var cmd tea.Cmd
	m.shellsList, cmd = m.shellsList.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return m.shellsList.View()
}
