package filepicker

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type model struct {
	filepicker   filepicker.Model
	selectedPath string
	windowSize   tea.WindowSizeMsg
	viewChanger  view.ViewChanger
	err          error
}

type KeyMap struct {
	GoToTop  key.Binding
	GoToLast key.Binding
	Down     key.Binding
	Up       key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Back     key.Binding
	Open     key.Binding
	Select   key.Binding
	Enter    key.Binding
}

// DefaultKeyMap defines the default keybindings.
func NewKeyMap() KeyMap {
	return KeyMap{
		GoToTop:  key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "first")),
		GoToLast: key.NewBinding(key.WithKeys("G"), key.WithHelp("G", "last")),
		Down:     key.NewBinding(key.WithKeys("j", "down", "ctrl+n"), key.WithHelp("j", "down")),
		Up:       key.NewBinding(key.WithKeys("k", "up", "ctrl+p"), key.WithHelp("k", "up")),
		PageUp:   key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("pgup", "page up")),
		PageDown: key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("pgdown", "page down")),
		Back:     key.NewBinding(key.WithKeys("h", "backspace", "left", "esc"), key.WithHelp("h", "back")),
		Open:     key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "open")),
		Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open")),
		Select:   key.NewBinding(key.WithKeys("select"), key.WithHelp("select", "select")),
	}
}

func New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) *model {
	fp := filepicker.New()
	fp.DirAllowed = true
	fp.ShowPermissions = true
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.AutoHeight = true
	_, height := utils.GetWindowSize()
	fp.Height = height - 10
	fp.KeyMap = NewKeyMap()
	return &model{
		filepicker:  fp,
		windowSize:  windowSize,
		viewChanger: viewChanger,
	}
}

func (m *model) Init() tea.Cmd {
	return m.filepicker.Init()
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.filepicker.Height = msg.Height
		m.filepicker, cmd = m.filepicker.Update(msg)
		return m, cmd
	case clearErrorMsg:
		m.err = nil
	}

	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedPath = path

	}
	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedPath = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}
	return m, cmd
}

func (m *model) View() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedPath == "" {
		s.WriteString("Pick a Path:")
	} else {
		s.WriteString("Selected Path: " + m.filepicker.Styles.Selected.Render(m.selectedPath))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")

	return s.String()
}
