package utils

import (
	"fmt"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type shellItem struct {
	title       string
	description string
	path        string
}

func (m shellItem) FilterValue() string {
	return m.title
}

func (m shellItem) Title() string {
	return m.title
}

func (m shellItem) Description() string {
	return m.description
}

func (m shellItem) Path() string {
	return m.path
}

var shells []shellItem

func LoadShells() ([]shellItem, error) {
	shells = []shellItem{}
	items := []shellItem{
		{title: "PowerShell", description: "Integrated PowerShell", path: "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"},
		{title: "PowerShell Core", description: "PowerShell Core (pwsh)", path: "C:\\Program Files\\PowerShell\\7\\pwsh.exe"},
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
