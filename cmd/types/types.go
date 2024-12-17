package types

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	Checked     lipgloss.Style
}

// ProfileItem represents a profile item in the list
type ProfileItem struct {
	ItemTitle           string
	ItemDescription     string
	IsValid             bool
	Path                string
	Shell               string
	Name                string
	ShellVersion        string
	IsValidPath         bool
	IsValidShellVersion bool
	IsValidDescription  bool
	IsSelected          bool
}

func (p ProfileItem) Title() string       { return p.ItemTitle }
func (p ProfileItem) Description() string { return p.ItemDescription }
func (p ProfileItem) FilterValue() string { return p.Name }

func (p ProfileItem) GetPath() string { return p.Path }
func (p ProfileItem) GetName() string {
	n := strings.Split(p.Path, "\\")
	p.Name = n[len(n)-1]
	return p.Name
}
func (p ProfileItem) GetDescription() string       { return strings.TrimLeft(p.ItemDescription, " ") }
func (p ProfileItem) GetShell() string             { return strings.ToLower(p.Shell) }
func (p ProfileItem) GetIsValidPath() bool         { return p.IsValidPath }
func (p ProfileItem) GetIsValidDescription() bool  { return p.IsValidDescription }
func (p ProfileItem) GetIsValidShellVersion() bool { return p.IsValidShellVersion }
func (p ProfileItem) IsValidProfile() bool {
	return p.IsValidPath && p.IsValidShellVersion && p.IsValidDescription
}
func (p ProfileItem) IsSelectedProfile() bool { return p.IsSelected }

// ShellItem represents a shell item in the list

type ShellItem struct {
	ItemTitle       string
	Name            string
	Path            string
	ItemDescription string
	ShortName       string
	ShortNames      []string
	ProfilePaths    []string
	IsSelected      bool
}

// Implement the list.Item interface for ShellItem
func (m ShellItem) GetName() string {
	return strings.ToLower(m.Name)
}

func (m ShellItem) GetPath() string {
	return m.Path
}
func (m ShellItem) GetShortNames() []string {
	lowerShortNames := make([]string, len(m.ShortName))
	for i, name := range m.ShortNames {
		lowerShortNames[i] = strings.ToLower(name)
	}
	return lowerShortNames
}
func (m ShellItem) GetShortName() string {
	return strings.ToLower(m.ShortName)
}
func (m ShellItem) IsSelectedShell() bool { return m.IsSelected }
func (m ShellItem) Title() string         { return m.ItemTitle }
func (m ShellItem) Description() string   { return m.ShortName }
func (m ShellItem) FilterValue() string   { return m.ItemTitle }
