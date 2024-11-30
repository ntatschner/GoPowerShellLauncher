package utils

import (
	"fmt"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

func LoadShells() ([]types.ShellItem, error) {
	var shells []types.ShellItem
	items := []types.ShellItem{
		{ItemTitle: "PowerShell", ItemDescription: "PowerShell", Name: "PowerShell", ShortName: []string{"powershell", "all"}, Path: "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"},
		{ItemTitle: "PowerShell Core", ItemDescription: "PowerShell Core", Name: "PowerShell Core", ShortName: []string{"pwsh", "all"}, Path: "C:\\Program Files\\PowerShell\\7\\pwsh.exe"},
	}
	for s := range items {
		l.Logger.Info("Processing shell", "shell", items[s])
		err := validatePath(items[s].Path)
		if err != nil {
			l.Logger.Warn("Invalid shell path", "Error", err)
		}
		l.Logger.Info("Shell path is valid", "shell", items[s])
		shells = append(shells, items[s])
	}
	if len(shells) == 0 {
		l.Logger.Error("No valid shells found")
		return nil, fmt.Errorf("no valid shells found")
	}
	return shells, nil
}
