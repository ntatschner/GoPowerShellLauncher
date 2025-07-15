package utils

import (
	"fmt"
	"strings"

	"github.com/ntatschner/GoPowerShellLauncher/cmd/launcher"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

var executePowerShellProcess = launcher.ExecutePowerShellProcess

func SplitProfiles(profiles string) []string {
	return strings.Split(profiles, ",")
}

func LaunchProfilesFromCmd(profiles string, shell string) error {
	shellPath, err := lookPath(shell)
	if err != nil {
		return fmt.Errorf("failed to find shell '%s' in your system's PATH. Please ensure that it is installed and accessible", shell)
	}

	var profileList []string
	for _, profilePath := range SplitProfiles(profiles) {
		p, err := GetProfileProperties(profilePath)
		if err != nil {
			l.Logger.Error("Failed to get profile properties", "profile", profilePath, "error", err)
			return fmt.Errorf("failed to get properties for profile: %s", profilePath)
		}

		if !p.IsValid {
			l.Logger.Warn("Skipping invalid profile", "profile", p.Name)
			continue
		}

		if p.Shell != "all" && p.Shell != shell {
			l.Logger.Warn("Skipping profile not meant for this shell", "profile", p.Name, "profile_shell", p.Shell, "target_shell", shell)
			continue
		}
		profileList = append(profileList, profilePath)
	}

	if len(profileList) == 0 {
		return fmt.Errorf("no valid profiles found for shell '%s'", shell)
	}

	merged := MergeSelectedProfiles(profileList)
	launcherErr := executePowerShellProcess(merged, shellPath)
	if launcherErr != nil {
		l.Logger.Error("Failed to launch profiles", "Error", launcherErr)
		return launcherErr
	}

	return nil
}
