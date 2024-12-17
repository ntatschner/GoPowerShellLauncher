package utils

import (
	"fmt"
	"os"
	"strings"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func CreateShortcut(profilepaths []string, name string, path string, shell string) error {
	l.Logger.Info("Creating shortcut", "name", name, "path", path)
	if name == "" {
		l.Logger.Error("Shortcut name is null")
		return fmt.Errorf("shortcut name cannot be null")
	}
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		l.Logger.Error("Destination path is not valid", "destination", path)
		return fmt.Errorf("destination path is not valid")
	}
	l.Logger.Info("Path exists", "path", path)

	// Validate each profile path
	for _, profilepath := range profilepaths {
		_, err := os.Stat(profilepath)
		if err != nil {
			l.Logger.Error("Profile path doesn't exist", "error", err)
			return err
		}
		l.Logger.Info("Profile path exists", "profilepath", profilepath)
	}

	// Join the profile paths by a comma
	profilePaths := strings.Join(profilepaths, ",")
	l.Logger.Info("Joined profile paths", "profilePaths", profilePaths)

	// Get the current working directory for the application path
	cwd, err := os.Getwd()
	if err != nil {
		l.Logger.Error("Failed to get working directory", "error", err)
		return err
	}
	l.Logger.Info("Current working directory", "cwd", cwd)

	var appName string
	// Get the name of the application
	if len(os.Args) > 0 {
		appName = os.Args[0]
	} else {
		appName = "GoPowerShellLauncher"
	}
	l.Logger.Info("Application name", "appName", appName)

	profilesCommand := "profiles --path"
	shellCommand := " --shell"
	finalCommand := fmt.Sprintf("%s %s %s %s", profilesCommand, profilePaths, shellCommand, shell)
	l.Logger.Info("Final command", "finalCommand", finalCommand)

	_, perr := os.Stat(path)
	if perr != nil {
		l.Logger.Error("Destination path doesn't exist", "error", perr)
		return perr
	}
	l.Logger.Info("Destination path exists", "path", path)

	shortcutPath := fmt.Sprintf("%s%c%s.lnk", path, os.PathSeparator, name)
	createShortcutCommand := fmt.Sprintf(
		"$ws = New-Object -ComObject WScript.Shell; "+
			"$s = $ws.CreateShortcut('%s'); "+
			"$s.TargetPath = '%s'; "+
			"$s.Arguments = '%s'; "+
			"$s.Description = 'Shortcut to launch GoPowerShellLauncher with selected profiles'; "+
			"$s.Save()",
		shortcutPath,
		appName,
		finalCommand,
	)
	l.Logger.Info("Create shortcut command", "createShortcutCommand", createShortcutCommand)

	encodedCommand, encodeerr := EncodeCommand(createShortcutCommand)
	if encodeerr != nil {
		l.Logger.Error("Failed to encode command", "error", encodeerr)
		return encodeerr
	}
	l.Logger.Info("Encoded command", "encodedCommand", encodedCommand)

	launcherr := ExecuteCommandWithPowershell(encodedCommand)
	if launcherr != nil {
		l.Logger.Error("Failed to create shortcut", "error", launcherr)
		return launcherr
	}
	l.Logger.Info("Shortcut created successfully", "shortcutPath", shortcutPath)

	return nil
}
