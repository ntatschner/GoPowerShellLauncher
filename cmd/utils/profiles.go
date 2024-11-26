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
	name := strings.Split(p.path, "\\")
	return name[len(name)-1]
}
func (p profile) Description() string       { return p.description }
func (p profile) Hash() string              { return p.hash }
func (p profile) Shell() string             { return p.shellVersion }
func (p profile) IsValidHash() bool         { return p.isValidHash }
func (p profile) IsValidPath() bool         { return p.isValidPath }
func (p profile) IsValidDescription() bool  { return p.isValidDescription }
func (p profile) IsValidShellVersion() bool { return p.isValidShellVersion }
func (p profile) Valid() bool {
	return p.isValidPath && p.isValidHash && p.isValidShellVersion && p.isValidDescription
}
func (p profile) FilterValue() string { return p.name }

var profiles []profile

func LoadProfile(line []string) profile {
	l.Logger.Info("Loading profile", "line", line)
	p := profile{}
	p.name = strings.Split(line[0], string(os.PathSeparator))[len(strings.Split(line[0], string(os.PathSeparator)))-1]
	p.path = line[0]
	p.hash = line[1]
	p.shellVersion = line[2]
	p.description = line[3]
	p.isValidPath = validatePath(p.path) == nil
	if p.isValidPath {
		p.isValidHash = validateHash(p.hash, p.path) == nil
	} else {
		p.isValidHash = false
	}
	p.isValidShellVersion = validateShellVersion(p.shellVersion) == nil
	p.isValidDescription = validateDescription(p.description) == nil
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
		return nil, err
	}
	// Validate csv headers
	headers := []string{"path", "hash", "shellversion", "description"}
	if err := validateHeaders(records[0], headers); err != nil {
		l.Logger.Error("Invalid CSV headers 🤯:", "expected", headers, "found", records[0])
		return nil, err
	}
	l.Logger.Info(fmt.Sprintf("Loaded %d records from CSV file", len(records)-1))
	var profiles []profile
	for _, record := range records[1:] {
		profile := LoadProfile(record)
		l.Logger.Info(fmt.Sprintf("Processing profile: %+v", profile))
		profiles = append(profiles, profile)
		l.Logger.Info(fmt.Sprintf("Added profile: %+v", profile))
	}
	l.Logger.Info(fmt.Sprintf("Total profiles loaded: %d", len(profiles)))
	return profiles, nil
}
