package styles

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)
)
var (
	TitleStyle        = lipgloss.NewStyle().MarginLeft(2).BorderStyle(lipgloss.RoundedBorder()).Align(lipgloss.Left)
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
	NormalTitle = lipgloss.NewStyle().Padding(0, 0, 0, 2).Foreground(lipgloss.Color("#78a8f5"))
	NormalDesc  = lipgloss.NewStyle().Foreground(lipgloss.Color("#0043b0"))

	SelectedTitle = lipgloss.NewStyle().Inherit(NormalTitle).Bold(true)
	SelectedDesc  = lipgloss.NewStyle().Inherit(NormalDesc).Bold(true)

	DimmedTitle = lipgloss.NewStyle().Inherit(NormalTitle).Faint(true)
	DimmedDesc  = lipgloss.NewStyle().Inherit(NormalDesc).Faint(true)

	Match = lipgloss.NewStyle().Inherit(NormalTitle).Underline(true)

	StatusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#507fcc")).Bold(true).Render
)

// ProfileSelectorItemStyles defines styling for a profile selector item.

type ProfileItemDelegate struct {
	*list.DefaultDelegate
}

func (pd ProfileItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(types.ProfileItem)
	if !ok {
		l.Logger.Errorf("Expected types.ProfileItem but got %T", listItem)
		return
	}

	msg := "Valid: "
	valid := msg + "❌"
	if i.IsValid {
		valid = msg + "✅"
	}

	title := fmt.Sprintf("%s | %s | Defined Shells: %s", i.GetName(), valid, i.GetShell())
	desc := i.GetDescription()
	outString := fmt.Sprintf("%s\n%s", title, desc)
	fn := NormalTitle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedTitle.Render("👉 " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(outString))
}

func (pd ProfileItemDelegate) Height() int  { return 1 }
func (pd ProfileItemDelegate) Spacing() int { return 0 }

func NewItemDelegate(keys *delegateKeyMap) (list.ItemDelegate, error) {
	if keys == nil {
		l.Logger.Error("keys is nil")
		return nil, fmt.Errorf("keys is nil")
	}
	d := ProfileItemDelegate{}

	d.Styles.NormalTitle = NormalTitle
	d.Styles.FilterMatch = Match

	d.Styles.SelectedTitle = SelectedTitle

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(types.ProfileItem); ok {
			title = i.GetName()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(StatusMessageStyle("Selected: " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(StatusMessageStyle("Removed: " + title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	l.Logger.Debug("Created item delegate", "delegate", d)
	return d, nil
}

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func NewDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp(" ", "Select Profile"),
		),
		remove: key.NewBinding(
			key.WithKeys("delete", " "),
			key.WithHelp("delete", "Deselect Profile"),
		),
	}
}
