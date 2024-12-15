package utils

import (
	"fmt"
	"os"
	"strings"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func CreateShortcut(profilepaths []string, name string, path string) error {
	l.Logger.Info("Creating shortcut", "name", name, "path", path)
	// Check if the path exists
	_, err := os.Stat(path)
	if err != nil {
		l.Logger.Error("Path doesn't exist", "error", err)
		return err
	}
	// Validate each profile path
	for _, profilepath := range profilepaths {
		_, err := os.Stat(profilepath)
		if err != nil {
			l.Logger.Error("Profile path doesn't exist", "error", err)
			return err
		}
	}
	// join the profile paths by a comma
	profilePaths := strings.Join(profilepaths, ",")
	// Get the current working directory to for the application path
	cwd, err := os.Getwd()
	if err != nil {
		l.Logger.Error("Failed to get working directory", "error", err)
		return err
	}
	var appName string
	// Get the name of the application
	if len(os.Args) > 0 {
		appName = os.Args[0]
	} else {
		appName = "GoPowerShellLauncher"
	}
	launchCommand := "profiles --path "

	appFullPath := fmt.Sprintf("%s%s%s", cwd, os.PathSeparator, appName)

	finalCommand := fmt.Sprintf("%s %s", launchCommand, profilePaths)
	_, perr := os.Stat(path)
	if perr != nil {
		l.Logger.Error("Destination path doesn't exist", "error", perr)
		return perr
	}
	shortcutPath := fmt.Sprintf("%s%s%s.lnk", path, os.PathSeparator, name)
	createShortcutCommand := fmt.Sprintf(
		"$ws = New-Object -ComObject WScript.Shell; "+
			"$s = $ws.CreateShortcut('%s'); "+
			"$s.TargetPath = '%s'; "+
			"$s.Arguments = '%s'; "+
			"$s.Description = 'Shortcut to launch GoPowerShellLauncher with selected profiles'; "+
			"$s.Save()",
		shortcutPath,
		appFullPath,
		finalCommand,
	)
	encodedCommand, encodeerr := EncodeCommand(createShortcutCommand)
	if encodeerr != nil {
		l.Logger.Error("Failed to encode command", "error", encodeerr)
		return encodeerr
	}
	launcherr := ExecuteCommandWithPowershell(encodedCommand)
	if launcherr != nil {
		l.Logger.Error("Failed to create shortcut", "error", launcherr)
		return launcherr
	}
	return nil
}
