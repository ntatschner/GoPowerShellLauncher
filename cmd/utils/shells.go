package utils

import (
	"fmt"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

var shells []types.ShellItem

func LoadShells() ([]types.ShellItem, error) {
	shells = []types.ShellItem{}
	items := []types.ShellItem{
		{Name: "PowerShell", ShortName: []string{"powershell", "all"}, Path: "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"},
		{Name: "PowerShell Core", ShortName: []string{"pwsh", "all"}, Path: "C:\\Program Files\\PowerShell\\7\\pwsh.exe"},
	}
	for s := range items {
		err := validatePath(shells[s].Path)
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
