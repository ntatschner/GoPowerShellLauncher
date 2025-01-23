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

type ProfileItemDelegate struct {
	ShowDescription bool
	Styles          types.DefaultItemStyles
	UpdateFunc      func(msg tea.Msg, m *list.Model) tea.Cmd
	ShortHelpFunc   func() []key.Binding
	FullHelpFunc    func() [][]key.Binding
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

	var selectedProfile string
	if i.IsSelectedProfile() {
		selectedProfile = "✓"
	} else {
		selectedProfile = ""
	}

	selectedProfile = s.Checked.Render(selectedProfile)

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
		fmt.Fprintf(w, "%s %s\n%s", title, selectedProfile, desc) //nolint: errcheck
		return
	}
	fmt.Fprintf(w, "%s %s", title, selectedProfile) //nolint: errcheck
}

func NewDefaultProfileStyles() (s types.DefaultItemStyles) {
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

	s.Checked = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF94F4")).Bold(true)

	return s
}

func NewProfileItemDelegate(keys *profiledelegateKeyMap) (*ProfileItemDelegate, error) {
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
				statusMsgText := "Selected: " + title

				return m.NewStatusMessage(StatusMessageStyle(statusMsgText))

			case key.Matches(msg, keys.unselected):
				statusMsgText := "Unselected: " + title

				return m.NewStatusMessage(StatusMessageStyle(statusMsgText))
			}
		case StatusBarUpdate:
			if bool(msg) {
				statusMsgText := "Selected: " + title

				return m.NewStatusMessage(StatusMessageStyle(statusMsgText))
			} else {
				statusMsgText := "Unselected: " + title

				return m.NewStatusMessage(StatusMessageStyle(statusMsgText))
			}
		}

		return nil
	}
	l.Logger.Debug("Created item delegate UpdateFunc")

	help := []key.Binding{keys.selected, keys.unselected, keys.view}
	l.Logger.Debug("Created item delegate help", "help", help)

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	l.Logger.Debug("Created item delegate", "delegate", d)
	return d, nil
}

type profiledelegateKeyMap struct {
	selected   key.Binding
	unselected key.Binding
	view       key.Binding
	navigation key.Binding
}

func (d profiledelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.selected,
		d.view,
		d.navigation,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d profiledelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.selected,
			d.unselected,
			d.view,
			d.navigation,
		},
	}
}

func NewProfileDelegateKeyMap() (*profiledelegateKeyMap, error) {
	d := &profiledelegateKeyMap{
		selected: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "(De)Select Profile"),
		),
		view: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "View Profile"),
		),
		navigation: key.NewBinding(
			key.WithKeys("ctrl+left", "ctrl+right"),
			key.WithHelp("ctrl+←/→", "Navigate"),
		),
	}
	l.Logger.Debug("Created delegate key map", "profiledelegateKeyMap", d)
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

// ShortHelp returns the delegate's short help.
func (d ProfileItemDelegate) ShortHelp() []key.Binding {
	if d.ShortHelpFunc != nil {
		return d.ShortHelpFunc()
	}
	return nil
}

// FullHelp returns the delegate's full help.
func (d ProfileItemDelegate) FullHelp() [][]key.Binding {
	if d.FullHelpFunc != nil {
		return d.FullHelpFunc()
	}
	return nil
}
