package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const ellipsis = "…"

type StatusBarUpdate bool

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)
)
var (
	TitleStyle = lipgloss.NewStyle().MarginLeft(2).BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#40C1AC")).Align(lipgloss.Left)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1).Faint(true).Align(lipgloss.Left)
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var (
	ViewPortTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	ViewPortInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return ViewPortTitleStyle.BorderStyle(b)
	}()
)

// ProfileSelector delegates and styles

var (
	StatusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF94F4")).Bold(true).Render
)
