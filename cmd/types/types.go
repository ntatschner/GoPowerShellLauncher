package types

import (
	"strings"
)

// ProfileItem represents a profile item in the list
type ProfileItem struct {
	ItemTitle           string
	ItemDescription     string
	Valid               string
	IsValid             bool
	Path                string
	Hash                string
	Shell               string
	Name                string
	ShellVersion        string
	IsValidHash         bool
	IsValidPath         bool
	IsValidShellVersion bool
	IsValidDescription  bool
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
func (p ProfileItem) GetHash() string              { return p.Hash }
func (p ProfileItem) GetShell() string             { return strings.ToLower(p.Shell) }
func (p ProfileItem) GetIsValidHash() bool         { return p.IsValidHash }
func (p ProfileItem) GetIsValidPath() bool         { return p.IsValidPath }
func (p ProfileItem) GetIsValidDescription() bool  { return p.IsValidDescription }
func (p ProfileItem) GetIsValidShellVersion() bool { return p.IsValidShellVersion }
func (p ProfileItem) IsValidProfile() bool {
	return p.IsValidPath && p.IsValidHash && p.IsValidShellVersion && p.IsValidDescription
}

// ShellItem represents a shell item in the list

type ShellItem struct {
	ItemTitle       string
	Name            string
	Path            string
	ItemDescription string
	ShortName       []string
	ProfilePaths    []string
}

// Implement the list.Item interface for ShellItem
func (m ShellItem) GetName() string {
	return strings.ToLower(m.Name)
}

func (m ShellItem) GetPath() string {
	return m.Path
}
func (m ShellItem) GetShortName() []string {
	lowerShortNames := make([]string, len(m.ShortName))
	for i, name := range m.ShortName {
		lowerShortNames[i] = strings.ToLower(name)
	}
	return lowerShortNames
}
func (m ShellItem) Title() string       { return m.ItemTitle }
func (m ShellItem) Description() string { return m.ItemTitle }
func (m ShellItem) FilterValue() string { return m.ItemTitle }
