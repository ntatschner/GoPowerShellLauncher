package aboutview

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
)

func init() {
	ui.RegisterPage("aboutView", Page{})
}

type Page struct{}

func (p Page) New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) tea.Model {
	return New(viewChanger, windowSize)
}

func (p Page) Title() string {
	return "About"
}

func (p Page) Description() string {
	return "Information about the application."
}

type model struct {
	viewChanger view.ViewChanger
	windowSize  tea.WindowSizeMsg
}

func New(viewChanger view.ViewChanger, windowSize tea.WindowSizeMsg) *model {
	return &model{
		viewChanger: viewChanger,
		windowSize:  windowSize,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *model) View() string {
	return fmt.Sprintf("GoPowerShellLauncher v%s", "1.0.0")
}
