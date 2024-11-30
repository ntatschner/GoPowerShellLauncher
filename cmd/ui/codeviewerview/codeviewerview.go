package codeviewerview

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/styles"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/utils"
)

type model struct {
	codeviewer  viewport.Model
	profilePath string
	viewChanger view.ViewChanger
	windowSize  tea.WindowSizeMsg
	help        string
}

const useHighPerformanceRenderer = false

func New(path string, windowSize tea.WindowSizeMsg, viewChanger view.ViewChanger) model {
	l.Logger.Info("Initializing code viewer")
	content, err := utils.LoadProfileContent(path)
	if err != nil {
		l.Logger.Error("Failed to load profile content", "error", err)
		content = "Failed to load profile content"
	}
	l.Logger.Info("Loaded profile content", "content", content)

	title := styles.TitleStyle.Render(path)
	vp := viewport.New(windowSize.Width, windowSize.Height)
	line := strings.Repeat("─", max(0, vp.Width-lipgloss.Width(title)))
	header := lipgloss.JoinHorizontal(lipgloss.Center, title, line)
	headerHeight := lipgloss.Height(header)
	info := styles.ViewPortInfoStyle.Render(fmt.Sprintf("%3.f%%", vp.ScrollPercent()*100))
	footer := lipgloss.JoinHorizontal(lipgloss.Center, line, info)
	footerHeight := lipgloss.Height(footer)
	help := styles.HelpStyle.Render("↑/k: up, ↓/j: down, u: ½ page up, d: ½ page down, backspace: back")
	helpHeight := lipgloss.Height(help)
	verticalMarginHeight := headerHeight + footerHeight + helpHeight
	vp.YPosition = headerHeight
	vp.HighPerformanceRendering = useHighPerformanceRenderer
	vp.SetContent(content)
	vp.VisibleLineCount()
	vp.KeyMap = viewport.DefaultKeyMap()
	vp.MouseWheelEnabled = true

	vp.Height = windowSize.Height - verticalMarginHeight
	return model{
		codeviewer:  vp,
		profilePath: path,
		viewChanger: viewChanger,
		windowSize:  windowSize,
		help:        help,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.codeviewer.Width = msg.Width
		m.codeviewer.Height = msg.Height - m.codeviewer.YPosition - 1
	}
	var cmd tea.Cmd
	m.codeviewer, cmd = m.codeviewer.Update(msg)
	return m, cmd
}

func (m model) View() string {
	title := styles.TitleStyle.Render(m.profilePath)
	line := strings.Repeat("─", max(0, m.codeviewer.Width-lipgloss.Width(title)))
	header := lipgloss.JoinHorizontal(lipgloss.Center, title, line)
	finfo := styles.ViewPortTitleStyle.Render(fmt.Sprintf("%3.f%%", m.codeviewer.ScrollPercent()*100))
	fline := strings.Repeat("─", max(0, m.codeviewer.Width-lipgloss.Width(finfo)))
	footer := lipgloss.JoinVertical(lipgloss.Left, lipgloss.JoinHorizontal(lipgloss.Center, fline, finfo), m.help)
	return lipgloss.JoinVertical(lipgloss.Left, header, m.codeviewer.View(), footer)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
