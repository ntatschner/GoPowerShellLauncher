package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type profile struct {
	name                string
	path                string
	hash                string
	shellVersion        string
	description         string
	isValidHash         bool
	isValidPath         bool
	isValidShellVersion bool
	isValidDescription  bool
}

func (p profile) Path() string { return p.path }
func (p profile) Name() string {
	n := strings.Split(p.path, "\\")
	p.name = n[len(n)-1]
	return p.name
}
func (p profile) Description() string       { return strings.TrimLeft(p.description, " ") }
func (p profile) Hash() string              { return p.hash }
func (p profile) Shell() string             { return p.shellVersion }
func (p profile) IsValidHash() bool         { return p.isValidHash }
func (p profile) IsValidPath() bool         { return p.isValidPath }
func (p profile) IsValidDescription() bool  { return p.isValidDescription }
func (p profile) IsValidShellVersion() bool { return p.isValidShellVersion }
func (p profile) Valid() bool {
	return p.isValidPath && p.isValidHash && p.isValidShellVersion && p.isValidDescription
}
func (p profile) FilterValue() string { return p.path }

var profiles []profile

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
// - profile: a profile struct with the loaded and validated data.
func LoadProfile(line []string) profile {
	l.Logger.Info("Loading profile", "line", line)
	p := profile{
		path:         line[0],
		hash:         line[1],
		shellVersion: line[2],
		description:  line[3],
	}

	p.isValidPath = validateField(p.path, validatePath, "path")
	p.isValidHash = validateField(p.hash, func(hash string) error { return validateHash(hash, p.path) }, "hash")
	p.isValidShellVersion = validateField(p.shellVersion, validateShellVersion, "shell version")
	p.isValidDescription = validateField(p.description, validateDescription, "description")
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

func LoadProfiles(filePath string) ([]profile, error) {
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
	var profiles []profile
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
