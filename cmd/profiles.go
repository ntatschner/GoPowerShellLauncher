package cmd

import (
	"encoding/csv"
	"log"
	"os"
)

type profile struct {
	path                string
	hash                string
	shellVersion        string
	isValidHash         bool
	isValidPath         bool
	isValidShellVersion bool
}

var profiles []profile

func LoadProfile(line []string) profile {
	p := profile{}
	p.path = line[0]
	p.hash = line[1]
	p.shellVersion = line[2]
	p.isValidPath = validatePath(p.path) == nil
	if p.isValidPath {
		p.isValidHash = validateHash(p.hash, p.path) == nil
	} else {
		p.isValidHash = false
	}
	p.isValidShellVersion = validateShellVersion(p.shellVersion) == nil
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

	log.Printf("Loaded %d records from CSV file", len(records)-1)

	for _, record := range records[1:] {
		profile := LoadProfile(record)
		log.Printf("Processing profile: %+v", profile)
		profiles = append(profiles, profile)
		log.Printf("Added profile: %+v", profile)
	}
	log.Printf("Total profiles loaded: %d", len(profiles))
	return nil
}
