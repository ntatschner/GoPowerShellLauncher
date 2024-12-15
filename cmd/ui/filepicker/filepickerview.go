package filepicker

import (
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
)

type model struct {
	filepicker   filepicker.Model
	selectedFile string
	windowSize   tea.WindowSizeMsg
	viewChanger  view.ViewChanger
	err          error
}

func New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) *model {
	fp := filepicker.New()
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.CurrentDirectory, _ = os.UserHomeDir()
	return &model{
		filepicker:   fp,
		selectedFile: "",
		windowSize:   windowSize,
		viewChanger:  viewChanger,
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
		m.filepicker.Width = msg.Width
		m.filepicker.Height = msg.Height
	case clearErrorMsg:
		m.err = nil
	}

	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}
	return m, cmd
}

func (m *model) View() string {
	return m.filepicker.View()
}
