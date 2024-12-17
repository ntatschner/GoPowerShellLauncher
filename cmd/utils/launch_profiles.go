package utils

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ntatschner/GoPowerShellLauncher/cmd/launcher"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

func SplitProfiles(profiles string) []string {
	return strings.Split(profiles, ",")
}

func LaunchProfilesFromCmd(profiles string, shell string) error {
	var profileList []string
	shellPath, err := exec.LookPath(shell)
	if err != nil {
		l.Logger.Error("Failed to find shell", "Error", err)
	}

	for _, profile := range SplitProfiles(profiles) {
		p, errProfile := GetProfileProperties(profile)
		if errProfile != nil {
			l.Logger.Error("Failed to get profile properties", "Error", errProfile)
			return errProfile
		}
		if p.Shell == shell {
			profileList = append(profileList, profile)
		}
		if profileList == nil {
			l.Logger.Warn("No profiles passed were validated for the shell", "shell", shell)
			return fmt.Errorf("no profiles passed were validated for the shell")
		}
		merged := MergeSelectedProfiles(profileList)
		encodedcommand, encodeErr := EncodeCommand(merged)
		if encodeErr != nil {
			l.Logger.Error("Failed to encode command", "Error", encodeErr)
			return encodeErr
		}
		launcherErr := launcher.ExecutePowerShellProcess(encodedcommand, shellPath)
		if launcherErr != nil {
			l.Logger.Error("Failed to launch profiles", "Error", launcherErr)
			return launcherErr
		}
	}
	return nil
}
