package utils

import (
	"encoding/csv"
	"fmt"
	"os"

	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type profile struct {
	path                string
	hash                string
	shellVersion        string
	description         string
	isValidHash         bool
	isValidPath         bool
	isValidShellVersion bool
	isValidDescription  bool
}

var profiles []profile

func LoadProfile(line []string) profile {
	l.Logger.Info("Loading profile", "line", line)
	p := profile{}
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

func loadProfiles(filePath string) error {
	// clear current profiles
	profiles = []profile{}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	l.Logger.Info(fmt.Sprintf("Loaded %d records from CSV file", len(records)-1))

	for _, record := range records[1:] {
		profile := LoadProfile(record)
		l.Logger.Info(fmt.Sprintf("Processing profile: %+v", profile))
		profiles = append(profiles, profile)
		l.Logger.Info(fmt.Sprintf("Added profile: %+v", profile))
	}
	l.Logger.Info(fmt.Sprintf("Total profiles loaded: %d", len(profiles)))
	return nil
}
