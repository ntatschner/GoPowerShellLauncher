package utils

import (
	"encoding/csv"
	"fmt"
	"os"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

func validateField(field string, validateFunc func(string) error, fieldName string) bool {
	if err := validateFunc(field); err != nil {
		l.Logger.Error(fmt.Sprintf("Failed to validate %s", fieldName), "error", err)
		return false
	}
	return true
}

// LoadProfile loads a profile from a CSV line.
// It validates the path, hash, shell version, and description of the profile.
// Parameters:
// line: a slice of strings containing the profile data.
// Returns:
// - ProfileItem: a ProfileItem struct with the loaded and validated data.
func LoadProfile(line []string) types.ProfileItem {
	l.Logger.Info("Loading profile", "line", line)
	p := types.ProfileItem{
		Path:            line[0],
		Hash:            line[1],
		Shell:           line[2],
		ItemDescription: line[3],
	}
	p.ItemTitle = p.GetName()

	p.IsValidPath = validateField(p.Path, validatePath, "path")
	p.IsValidHash = validateField(p.Hash, func(hash string) error { return validateHash(hash, p.Path) }, "hash")
	p.IsValidShellVersion = validateField(p.Shell, validateShellVersion, "shell version")
	p.IsValidDescription = validateField(p.ItemDescription, validateDescription, "description")
	p.IsValid = p.IsValidPath && p.IsValidHash && p.IsValidShellVersion && p.IsValidDescription
	return p
}

func validateHeaders(headers []string, expectedHeaders []string) error {
	if len(headers) != len(expectedHeaders) {
		return fmt.Errorf("invalid number of headers: got %d, expected %d", len(headers), len(expectedHeaders))
	}
	for i, header := range headers {
		if header != expectedHeaders[i] {
			return fmt.Errorf("invalid header: got %s, expected %s", header, expectedHeaders[i])
		}
	}
	return nil
}

func LoadProfiles(filePath string) ([]types.ProfileItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		l.Logger.Error("Failed to open CSV file", "path", filePath, "error", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		l.Logger.Error("Failed to read CSV file", "error", err)
		return nil, err
	}

	// Validate csv headers
	expectedHeaders := []string{"path", "hash", "shellversion", "description"}
	if err := validateHeaders(records[0], expectedHeaders); err != nil {
		l.Logger.Error("Invalid CSV headers", "error", err)
		return nil, err
	}

	l.Logger.Info(fmt.Sprintf("Loaded %d records from CSV file", len(records)-1))
	var profiles []types.ProfileItem
	for i, record := range records[1:] {
		if len(record) != len(expectedHeaders) {
			l.Logger.Error("Wrong number of fields", "line", i+2, "record", record)
			continue
		}
		profile := LoadProfile(record)
		l.Logger.Info(fmt.Sprintf("Processing profile: %+v", profile))
		profiles = append(profiles, profile)
		l.Logger.Info(fmt.Sprintf("Added profile: %+v", profile))
	}
	l.Logger.Info(fmt.Sprintf("Total profiles loaded: %d", len(profiles)))
	return profiles, nil
}
