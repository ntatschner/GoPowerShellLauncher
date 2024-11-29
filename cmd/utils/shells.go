package utils

import (
	"fmt"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type shellItem struct {
	title       string
	description string
	name        string
	path        string
	shortName   []string
}

func (m shellItem) Name() string {
	return m.name
}

func (m shellItem) Path() string {
	return m.path
}
func (m shellItem) ShortName() []string {
	return m.shortName
}
func (m shellItem) Title() string       { return m.name }
func (m shellItem) Description() string { return "Shell for " + m.name }

var shells []shellItem

func LoadShells() ([]shellItem, error) {
	shells = []shellItem{}
	items := []shellItem{
		{name: "PowerShell", shortName: []string{"powershell", "all"}, path: "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"},
		{name: "PowerShell Core", shortName: []string{"pwsh", "all"}, path: "C:\\Program Files\\PowerShell\\7\\pwsh.exe"},
	}
	for s := range items {
		err := validatePath(shells[s].path)
		if err != nil {
			l.Logger.Warn("Invalid shell path", "Error", err)
		}
		shells = append(shells, items[s])
	}
	if len(shells) == 0 {
		l.Logger.Error("No valid shells found")
		return nil, fmt.Errorf("no valid shells found")
	}
	return shells, nil
}
