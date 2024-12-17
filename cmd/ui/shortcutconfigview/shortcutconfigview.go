package shortcutconfigview

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Create ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Create"))
)

type model struct {
	focusIndex  int
	inputs      []textinput.Model
	windowSize  tea.WindowSizeMsg
	viewChanger view.ViewChanger
	profiles    []types.ProfileItem
	shell       types.ShellItem
}

func New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg, profiles []types.ProfileItem, shell types.ShellItem) *model {
	m := &model{
		inputs:      make([]textinput.Model, 2),
		shell:       shell,
		profiles:    profiles,
		viewChanger: viewChanger,
		windowSize:  windowSize,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Name of the shortcut"
			t.Prompt = "Name: "
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Destination Path"
			t.Prompt = "Destination: "
			t.Validate = func(s string) error {
				if _, err := os.Stat(s); os.IsNotExist(err) {
					return fmt.Errorf("destination Path is not valid")
				}
				return nil
			}
			t.CharLimit = 64
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Set focus to next input
		case "enter", "up", "down":
			s := msg.String()
			if s == "enter" && m.focusIndex == len(m.inputs) {
				name := m.inputs[0].Value()
				destination := m.inputs[1].Value()
				l.Logger.Info("Creating shortcut", "name", name, "destination", destination, "profiles", m.shell.ProfilePaths)

				err := utils.CreateShortcut(m.shell.ProfilePaths, name, destination)
				if err != nil {
					l.Logger.Error("Failed to create shortcut", "Error", err)
					return m, nil
				}
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

var (
	titleStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderBottom(true).
			Padding(0, 2).
			Align(lipgloss.Center).
			Render

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Render
)

func (m *model) View() string {
	var b strings.Builder

	var errString string
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if m.inputs[i].Err != nil {
			errString += m.inputs[i].Err.Error() + " "
		}
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	fmt.Fprintf(&b, "\n\n%s\n%s\n", *button, errString)

	// Get the window size
	width, height := m.windowSize.Width, m.windowSize.Height

	// Create the title
	title := titleStyle("Enter Shortcut Details")

	// Combine the title and the content
	content := lipgloss.JoinVertical(lipgloss.Left, title, b.String())

	// Add a border around the content
	borderedContent := borderStyle(content)

	// Center the content horizontally and vertically
	centeredContent := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, borderedContent)

	return centeredContent
}
