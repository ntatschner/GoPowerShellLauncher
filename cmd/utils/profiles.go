package utils

import (
	"encoding/csv"
	"fmt"
	"os"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
)

// LoadProfile loads a profile from a CSV line.
// It validates the path, hash, shell version, and description of the profile.
// Parameters:
// line: a slice of strings containing the profile data.
// Returns:
// - ProfileItem: a ProfileItem struct with the loaded and validated data.
func LoadProfile(line []string, hashValidator HashValidator) types.ProfileItem {
	l.Logger.Info("Loading profile", "line", line)
	p := types.ProfileItem{
		Path:            line[0],
		Hash:            line[1],
		Shell:           line[2],
		ItemDescription: line[3],
	}
	p.ItemTitle = p.GetName()
	isValidPath, patherr := ValidatePath(p.Path)
	if patherr != nil {
		l.Logger.Error(fmt.Sprintf("Failed to validate path %s", p.Path), "error", patherr)
	}
	p.IsValidPath = isValidPath
	isValidHash, hasherr := hashValidator.ValidateHash(p.Hash, p.Path)
	if hasherr != nil {
		l.Logger.Error(fmt.Sprintf("Failed to validate hash for path %s", p.Path), "error", hasherr)
	}
	p.IsValidHash = isValidHash
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
	p.IsValid = p.IsValidPath && p.IsValidHash && p.IsValidShellVersion && p.IsValidDescription
	l.Logger.Info("Profile loaded", "profile", p)
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
	var hashValidator HashValidator
	for i, record := range records[1:] {
		if len(record) != len(expectedHeaders) {
			l.Logger.Error("Wrong number of fields", "line", i+2, "record", record)
			continue
		}
		profile := LoadProfile(record, hashValidator)
		l.Logger.Info(fmt.Sprintf("Processing profile: %+v", profile))
		profiles = append(profiles, profile)
		l.Logger.Info(fmt.Sprintf("Added profile: %+v", profile))
	}
	l.Logger.Info(fmt.Sprintf("Total profiles loaded: %d", len(profiles)))
	return profiles, nil
}
