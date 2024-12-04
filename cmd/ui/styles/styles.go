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
	NormalTitle = lipgloss.NewStyle().Padding(0, 0, 3, 0).Foreground(lipgloss.Color("#78a8f5"))
	NormalDesc  = lipgloss.NewStyle().Padding(0, 0, 2, 0).Foreground(lipgloss.Color("#0043b0"))

	SelectedTitle = lipgloss.NewStyle().Inherit(NormalTitle).Bold(true)
	SelectedDesc  = lipgloss.NewStyle().Inherit(NormalDesc).Bold(true)

	DimmedTitle = lipgloss.NewStyle().Inherit(NormalTitle).Faint(true)
	DimmedDesc  = lipgloss.NewStyle().Inherit(NormalDesc).Faint(true)

	Match = lipgloss.NewStyle().Inherit(NormalTitle).Underline(true)

	StatusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#507fcc")).Bold(true).Render
)

type DefaultItemStyles struct {
	// The Normal state.
	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	// The selected item state.
	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	// The dimmed state, for when the filter input is initially activated.
	DimmedTitle lipgloss.Style
	DimmedDesc  lipgloss.Style

	// Characters matching the current filter, if any.
	FilterMatch lipgloss.Style
}

// ProfileSelectorItemStyles defines styling for a profile selector item.

type ProfileItemDelegate struct {
	Styles        DefaultItemStyles
 ShowDescription bool
	UpdateFunc    func(msg tea.Msg, m *list.Model) tea.Cmd
	ShortHelpFunc func() []key.Binding
	ShortHelp     func() []key.Binding
	FullHelpFunc  func() [][]key.Binding
	FullHelp      func() [][]key.Binding
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

	fn := NormalTitle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return SelectedTitle.Render("👉 " + strings.Join(s, " "))
		}
	} else if index != m.Index() && m.FilterState() != list.Filtering {
		fn = func(s ...string) string {
			return SelectedTitle.Render("    " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(title, desc))
}

func (pd ProfileItemDelegate) Height() int  { return 2 }
func (pd ProfileItemDelegate) Spacing() int { return 1 }
func (pd ProfileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if pd.UpdateFunc != nil {
		return pd.UpdateFunc(msg, m)
	}
	return nil
}

type StatusBarUpdate bool

func NewDefaultProfileStyles() (s DefaultItemStyles) {
 s.NormalTitle = lipgloss.NewStyle().Padding(0, 0, 3, 0).Foreground(lipgloss.Color("#78a8f5"))
	s.NormalDesc  = lipgloss.NewStyle().Padding(0, 0, 2, 0).Foreground(lipgloss.Color("#0043b0"))

	s.SelectedTitle = lipgloss.NewStyle().Inherit(NormalTitle).Bold(true)
	s.SelectedDesc  = lipgloss.NewStyle().Inherit(NormalDesc).Bold(true)

	s.DimmedTitle = lipgloss.NewStyle().Inherit(NormalTitle).Faint(true)
	s.DimmedDesc  = lipgloss.NewStyle().Inherit(NormalDesc).Faint(true)

	s.Match = lipgloss.NewStyle().Inherit(NormalTitle).Underline(true)

 return s
}

func NewItemDelegate(keys *delegateKeyMap) (*ProfileItemDelegate, error) {
	l.Logger.Debug("Creating item delegate", "keys", keys)
	if keys == nil {
		l.Logger.Error("keys is nil")
		return nil, fmt.Errorf("keys is nil")
	}
	//d := list.NewDefaultDelegate()
	d := &ProfileItemDelegate{}
	l.Logger.Debug("Created instance of ProfileItemDelegate item delegate", "delegate", d)
 d.ShowDescription = true
	d.Styles.NormalTitle = NormalTitle
	d.Styles.NormalDesc = NormalDesc
	d.Styles.SelectedTitle = SelectedTitle
	d.Styles.SelectedDesc = SelectedDesc
	d.Styles.DimmedTitle = DimmedTitle
	d.Styles.DimmedDesc = DimmedDesc
	d.Styles.FilterMatch = Match

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
			case key.Matches(msg, keys.selected):
				return m.NewStatusMessage(StatusMessageStyle("Selected: " + title))

			case key.Matches(msg, keys.unselected):
				return m.NewStatusMessage(StatusMessageStyle("Unselected: " + title))
			}
		case StatusBarUpdate:
			if bool(msg) {
				return m.NewStatusMessage(StatusMessageStyle("Selected: " + title))
			} else {
				return m.NewStatusMessage(StatusMessageStyle("Unselected: " + title))
			}
		}

		return nil
	}
	l.Logger.Debug("Created item delegate UpdateFunc")

	help := []key.Binding{keys.selected, keys.unselected}
	l.Logger.Debug("Created item delegate help", "help", help)

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	d.ShortHelp = d.ShortHelpFunc
	d.FullHelp = d.FullHelpFunc
	l.Logger.Debug("Created item delegate", "delegate", d)
	return d, nil
}

type delegateKeyMap struct {
	selected   key.Binding
	unselected key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.selected,
		d.unselected,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.selected,
			d.unselected,
		},
	}
}

func NewDelegateKeyMap() (*delegateKeyMap, error) {
	d := &delegateKeyMap{
		selected: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "Select Profile"),
		),
		unselected: key.NewBinding(
			key.WithKeys("delete", " "),
			key.WithHelp("delete", "Deselect Profile"),
		),
	}
	l.Logger.Debug("Created delegate key map", "delegateKeyMap", d)
	return d, nil
}
