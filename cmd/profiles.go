package cmd

import (
	"encoding/csv"
	"log"
	"os"
)

type profile struct {
	path         string
	hash         string
	shellVersion string
}

var profiles []profile

func (p *profile) IsValid() error {
	log.Printf("Validating profile: %+v", p)
	fv := new(FieldValidator)
	fv.InArray(p.shellVersion, []string{"pwsh", "powershell"})
	if !fv.IsValid() {
		log.Printf("Validation failed for profile: %+v", p)
		return fv.Error()
	}
	log.Printf("Validation succeeded for profile: %+v", p)
	return nil
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
		profile := profile{
			path:         record[0],
			hash:         record[1],
			shellVersion: record[2],
		}
		log.Printf("Processing profile: %+v", profile)
		if err := profile.IsValid(); err != nil {
			log.Printf("Invalid profile: %v", err)
			return err
		}
		profiles = append(profiles, profile)
		log.Printf("Added profile: %+v", profile)
	}
	log.Printf("Total profiles loaded: %d", len(profiles))
	return nil
}
