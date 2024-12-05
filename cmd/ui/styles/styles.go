package styles

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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

const ellipsis = "…"

// ProfileSelectorItemStyles defines styling for a profile selector item.

type ProfileItemDelegate struct {
	Styles          DefaultItemStyles
	ShowDescription bool
	UpdateFunc      func(msg tea.Msg, m *list.Model) tea.Cmd
	ShortHelpFunc   func() []key.Binding
	ShortHelp       func() []key.Binding
	FullHelpFunc    func() [][]key.Binding
	FullHelp        func() [][]key.Binding
	height          int
	spacing         int
}

func (d ProfileItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc  string
		matchedRunes []int
		s            = &d.Styles
	)

	i, ok := item.(types.ProfileItem)
	if !ok {
		l.Logger.Errorf("Expected types.ProfileItem but got %T", item)
		return
	}

	msg := "Valid: "
	valid := msg + "❌"
	if i.IsValid {
		valid = msg + "✅"
	}

	if i, ok := item.(types.ProfileItem); ok {
		title = fmt.Sprintf("%s | %s | Defined Shells: %s", i.GetName(), valid, i.GetShell())
		desc = i.GetDescription()
	} else {
		return
	}

	if m.Width() <= 0 {
		// short-circuit
		return
	}

	// Prevent text from exceeding list width
	textwidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.height-1 {
				break
			}
			lines = append(lines, ansi.Truncate(line, textwidth, ellipsis))
		}
		desc = strings.Join(lines, "\n")
	}

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if emptyFilter {
		title = s.DimmedTitle.Render("   " + title)
		desc = s.DimmedDesc.Render("   " + desc)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render("👉 " + title)
		desc = s.SelectedDesc.Render("   " + desc)
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render("   " + title)
		desc = s.NormalDesc.Render("   " + desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc) //nolint: errcheck
		return
	}
	fmt.Fprintf(w, "%s", title) //nolint: errcheck
}

type StatusBarUpdate bool

func NewDefaultProfileStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#008A74", Dark: "#40C1AC"}).
		Padding(0, 0, 0, 2) //nolint:mnd

	s.NormalDesc = s.NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#FF94F4"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#FF94F4", Dark: "#FF94F4"}).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).
		Padding(0, 0, 0, 2) //nolint:mnd

	s.DimmedDesc = s.DimmedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

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
	d.Styles = NewDefaultProfileStyles()
	d.height = 2
	d.spacing = 1

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

func (pd ProfileItemDelegate) Height() int {
	if pd.ShowDescription {
		return pd.height
	}
	return 1
}

// SetSpacing sets the delegate's spacing.
func (pd *ProfileItemDelegate) SetSpacing(i int) {
	pd.spacing = i
}

// Spacing returns the delegate's spacing.
func (pd ProfileItemDelegate) Spacing() int {
	return pd.spacing
}

// Update checks whether the delegate's UpdateFunc is set and calls it.
func (pd ProfileItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if pd.UpdateFunc == nil {
		return nil
	}
	return pd.UpdateFunc(msg, m)
}

func (pd *ProfileItemDelegate) SetHeight(i int) {
	pd.height = i
}
