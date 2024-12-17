package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

func LoadProfilesFromDir() ([]types.ProfileItem, error) {
	var profiles []types.ProfileItem
	configData, _ := LoadConfig()
	directory := configData.ProfilePath
	l.Logger.Info("Loading profiles from config directory", "dir", directory)
	recursive := configData.Recursive
	l.Logger.Info("Recursive search", "recursive", recursive)

	var processedFiles []string

	if recursive {
		err := filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				l.Logger.Error("Failed to access path", "path", path, "error", err)
				return err
			}
			if !d.IsDir() && strings.Contains(d.Name(), ".Profile.ps1") {
				fullPath := path
				processedFiles = append(processedFiles, fullPath)
				l.Logger.Info("File processed", "file", fullPath)
			}
			return nil
		})
		if err != nil {
			l.Logger.Error("Failed to walk directory", "dir", directory, "error", err)
			return nil, err
		}
	} else {
		files, err := os.ReadDir(directory)
		if err != nil {
			l.Logger.Error("Failed to read directory", "dir", directory, "error", err)
			return nil, err
		}
		for _, file := range files {
			l.Logger.Info("Processing file", "file", file.Name())
			if !file.IsDir() && strings.Contains(file.Name(), ".Profile.ps1") {
				fullPath := filepath.Join(directory, file.Name())
				processedFiles = append(processedFiles, fullPath)
				l.Logger.Info("File processed", "file", fullPath)
			}
		}
	}

	for _, file := range processedFiles {
		l.Logger.Info("Loading file", "file", file)
		profile, profileerr := GetProfileProperties(file)
		if profileerr != nil {
			l.Logger.Error("Failed to get profile properties", "error", profileerr)
		} else {
			l.Logger.Info("Profile loaded", "profile", profile)
			profiles = append(profiles, profile)
		}
	}
	return profiles, nil
}

func ExtractString(input string, pattern string) (string, error) {
	// Compile the regex pattern
	re := regexp.MustCompile(pattern)

	// Find the submatch
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", fmt.Errorf("no match found")
	}

	// Return the extracted string
	return matches[1], nil
}

func GetProfileProperties(path string) (types.ProfileItem, error) {
	l.Logger.Info("Getting profile properties", "path", path)
	// Get the .Profile.ps1 content, parse the file to get the required SHELL and DESCRIPTION using regex with these patterns:
	// ### SHELL:<SHELL>:SHELL ### and ### DESCRIPTION:<DESCRIPTION>:DESCRIPTION ###
	// Read the file content
	content, readerr := os.ReadFile(path)
	if readerr != nil {
		l.Logger.Error("Failed to read file", "path", path, "error", readerr)
		return types.ProfileItem{}, readerr
	}
	shell, shellerr := ExtractString(string(content), `### SHELL:(.*):SHELL ###`)
	if shellerr != nil {
		l.Logger.Error("Failed to extract shell", "error", shellerr)
		shell = "InvalidShell"
	}
	description, descerr := ExtractString(string(content), `### DESCRIPTION:(.*):DESCRIPTION ###`)
	if descerr != nil {
		l.Logger.Error("Failed to extract description", "error", descerr)
		description = ""
	}
	p := types.ProfileItem{
		Path:            path,
		Shell:           shell,
		ItemDescription: description,
	}
	p.ItemTitle = p.GetName()
	p.Name = p.GetName()
	p.IsSelected = false
	isValidPath, patherr := ValidatePath(p.Path)
	if patherr != nil {
		l.Logger.Error(fmt.Sprintf("Failed to validate path %s", p.Path), "error", patherr)
	}
	p.IsValidPath = isValidPath
	isValidShell, shellerr := ValidateShellVersion(p.Shell)
	if shellerr != nil {
		l.Logger.Error(fmt.Sprintf("Failed to validate shell version %s", p.Shell), "error", shellerr)
	}
	p.IsValidShellVersion = isValidShell
	isValidDescription, descerr := ValidateDescription(p.ItemDescription)
	if descerr != nil {
		l.Logger.Error(fmt.Sprintf("Failed to validate description %s", p.ItemDescription), "error", descerr)
	}
	p.IsValidDescription = isValidDescription
	p.IsValid = p.IsValidPath && p.IsValidShellVersion && p.IsValidDescription
	l.Logger.Info("Profile loaded", "profile", p)
	return p, nil
}
