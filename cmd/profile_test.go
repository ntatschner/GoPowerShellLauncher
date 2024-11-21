package cmd

import (
	"os"
	"testing"
)

func TestLoadProfiles(t *testing.T) {
	// Create a temporary CSV file for testing
	file, err := os.CreateTemp("", "profiles_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write test data to the CSV file
	data := `path,hash,shellVersion
/path/to/script1,hash1,pwsh
/path/to/script2,hash2,powershell
/path/to/script3,hash3,invalidShell`
	if _, err := file.Write([]byte(data)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	file.Close()

	// Test loading profiles
	err = loadProfiles(file.Name())
	if err == nil {
		t.Fatalf("Expected error for invalid shell version, got nil")
	}

	// Modify the test data to have valid shell versions
	data = `/path/to/script1,hash1,pwsh
/path/to/script2,hash2,powershell`
	if err := os.WriteFile(file.Name(), []byte(data), 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Test loading profiles again
	err = loadProfiles(file.Name())
	if err != nil {
		t.Fatalf("Failed to load profiles: %v", err)
	}

	// Check the loaded profiles
	if len(profiles) != 2 {
		t.Fatalf("Expected 2 profiles, got %d", len(profiles))
	}
	if profiles[0].path != "/path/to/script1" || profiles[0].shellVersion != "pwsh" {
		t.Errorf("Unexpected profile data: %+v", profiles[0])
	}
	if profiles[1].path != "/path/to/script2" || profiles[1].shellVersion != "powershell" {
		t.Errorf("Unexpected profile data: %+v", profiles[1])
	}
}
