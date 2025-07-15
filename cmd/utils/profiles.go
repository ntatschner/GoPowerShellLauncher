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

func findProfileFiles(root string, recursive bool) ([]string, error) {
	var files []string
	if recursive {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(d.Name(), ".Profile.ps1") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".Profile.ps1") {
				files = append(files, filepath.Join(root, entry.Name()))
			}
		}
	}
	return files, nil
}

func LoadProfilesFromDir(configPath ...string) ([]types.ProfileItem, error) {
	configData, err := LoadConfig(configPath...)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	directory := configData.Profile.Path
	l.Logger.Info("Loading profiles from config directory", "dir", directory)
	recursive := configData.Profile.Recursive
	l.Logger.Info("Recursive search", "recursive", recursive)

	profileFiles, err := findProfileFiles(directory, recursive)
	if err != nil {
		return nil, fmt.Errorf("failed to find profile files: %w", err)
	}

	var profiles []types.ProfileItem
	for _, file := range profileFiles {
		l.Logger.Info("Loading file", "file", file)
		profile, err := GetProfileProperties(file)
		if err != nil {
			l.Logger.Error("Failed to get profile properties", "file", file, "error", err)
			continue // Skip invalid profiles
		}
		l.Logger.Info("Profile loaded", "profile", profile)
		profiles = append(profiles, profile)
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

const (
	shellPattern      = `### SHELL:(.*):SHELL ###`
	descriptionPattern = `### DESCRIPTION:(.*):DESCRIPTION ###`
)

func parseProfile(content string) (shell, description string) {
	shell, err := ExtractString(content, shellPattern)
	if err != nil {
		l.Logger.Error("Failed to extract shell", "error", err)
		shell = "InvalidShell"
	}

	description, err = ExtractString(content, descriptionPattern)
	if err != nil {
		l.Logger.Error("Failed to extract description", "error", err)
		description = ""
	}

	return shell, description
}

func validateProfile(p *types.ProfileItem) {
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
}

func GetProfileProperties(path string) (types.ProfileItem, error) {
	l.Logger.Info("Getting profile properties", "path", path)

	content, err := os.ReadFile(path)
	if err != nil {
		return types.ProfileItem{}, fmt.Errorf("failed to read file: %w", err)
	}

	shell, description := parseProfile(string(content))

	p := types.ProfileItem{
		Path:            path,
		Shell:           shell,
		ItemDescription: description,
	}
	p.ItemTitle = p.GetName()
	p.Name = p.GetName()
	p.IsSelected = false

	validateProfile(&p)

	l.Logger.Info("Profile loaded", "profile", p)
	return p, nil
}
