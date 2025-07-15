package utils

import (
	"fmt"
	"os/exec"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

var lookPath = exec.LookPath

func LoadShells() ([]types.ShellItem, error) {
	var shells []types.ShellItem
	items := []types.ShellItem{
		{ItemTitle: "PowerShell", ItemDescription: "Windows PowerShell", Name: "PowerShell", ShortName: "powershell", ShortNames: []string{"powershell", "all"}},
		{ItemTitle: "PowerShell Core", ItemDescription: "PowerShell Core (pwsh)", Name: "PowerShell Core", ShortName: "pwsh", ShortNames: []string{"pwsh", "all"}},
	}

	for i := range items {
		path, err := lookPath(items[i].ShortName)
		if err != nil {
			l.Logger.Warn("Could not find shell in PATH", "shell", items[i].ShortName)
			continue // Skip shell if not found
		}
		items[i].Path = path
		l.Logger.Info("Found shell", "shell", items[i].Name, "path", path)
		shells = append(shells, items[i])
	}

	if len(shells) == 0 {
		l.Logger.Error("No valid shells found in PATH")
		return nil, fmt.Errorf("no valid shells found in your system's PATH. Please ensure that either 'powershell' or 'pwsh' is installed and accessible")
	}

	return shells, nil
}
