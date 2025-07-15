package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/types"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Initialize the logger
	err := logger.InitLogger("", "test.log", "DEBUG")
	if err != nil {
		panic(err)
	}
	// Run the tests
	os.Exit(m.Run())
}

func TestLoadProfilesFromDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test-profiles")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create some dummy profile files
	profile1 := `
### SHELL:PowerShell:SHELL ###
### DESCRIPTION:Test Profile 1:DESCRIPTION ###
`
	profile2 := `
### SHELL:PowerShell Core:SHELL ###
### DESCRIPTION:Test Profile 2:DESCRIPTION ###
`
	profile3 := `
### SHELL:InvalidShell:SHELL ###
### DESCRIPTION:Invalid Profile:DESCRIPTION ###
`
	assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, "profile1.Profile.ps1"), []byte(profile1), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, "profile2.Profile.ps1"), []byte(profile2), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, "profile3.Profile.ps1"), []byte(profile3), 0644))
	assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, "not-a-profile.txt"), []byte("hello"), 0644))

	// Create a subdirectory for recursive testing
	subDir := filepath.Join(tmpDir, "subdir")
	assert.NoError(t, os.Mkdir(subDir, 0755))
	profile4 := `
### SHELL:PowerShell:SHELL ###
### DESCRIPTION:Test Profile 4:DESCRIPTION ###
`
	assert.NoError(t, os.WriteFile(filepath.Join(subDir, "profile4.Profile.ps1"), []byte(profile4), 0644))

	// Create a dummy config file
	configContent := `
profile:
  path: "` + tmpDir + `"
  recursive: false
`
	configPath := filepath.Join(tmpDir, "config.yaml")
	assert.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))
	// XXX: This is a hack to override the default config path.
	// A better solution would be to make the config path configurable.
	t.Setenv("GOPOWERSHELLLAUNCHER_CONFIG", configPath)

	// Test non-recursive loading
	configContent = `
profile:
  path: "` + tmpDir + `"
  recursive: false
`
	assert.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))
	profiles, err := LoadProfilesFromDir(configPath)
	assert.NoError(t, err)
	assert.Len(t, profiles, 3)

	// Test recursive loading
	configContent = `
profile:
  path: "` + tmpDir + `"
  recursive: true
`
	assert.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))
	profiles, err = LoadProfilesFromDir(configPath)
	assert.NoError(t, err)
	assert.Len(t, profiles, 4)
}

func TestGetProfileProperties(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test-profile-*.Profile.ps1")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write some content to the file
	content := `
### SHELL:PowerShell:SHELL ###
### DESCRIPTION:This is a test profile.:DESCRIPTION ###
`
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	tmpFile.Close()

	// Test GetProfileProperties
	profile, err := GetProfileProperties(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, "PowerShell", profile.Shell)
	assert.Equal(t, "This is a test profile.", profile.ItemDescription)
	assert.True(t, profile.IsValid)
}

func TestParseProfile(t *testing.T) {
	content := `
### SHELL:PowerShell:SHELL ###
### DESCRIPTION:Test Profile:DESCRIPTION ###
`
	shell, description := parseProfile(content)
	assert.Equal(t, "PowerShell", shell)
	assert.Equal(t, "Test Profile", description)
}

func TestLoadShells(t *testing.T) {
	// Mock exec.LookPath
	originalLookPath := lookPath
	defer func() { lookPath = originalLookPath }()

	// Case 1: Both shells are found
	lookPath = func(file string) (string, error) {
		if file == "powershell" {
			return "/path/to/powershell", nil
		}
		if file == "pwsh" {
			return "/path/to/pwsh", nil
		}
		return "", fmt.Errorf("not found")
	}
	shells, err := LoadShells()
	assert.NoError(t, err)
	assert.Len(t, shells, 2)

	// Case 2: Only powershell is found
	lookPath = func(file string) (string, error) {
		if file == "powershell" {
			return "/path/to/powershell", nil
		}
		return "", fmt.Errorf("not found")
	}
	shells, err = LoadShells()
	assert.NoError(t, err)
	assert.Len(t, shells, 1)
	assert.Equal(t, "PowerShell", shells[0].Name)

	// Case 3: No shells are found
	lookPath = func(file string) (string, error) {
		return "", fmt.Errorf("not found")
	}
	shells, err = LoadShells()
	assert.Error(t, err)
	assert.Nil(t, shells)
}

func TestLaunchProfilesFromCmd(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test-profiles")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create some dummy profile files
	profile1 := `
### SHELL:powershell:SHELL ###
### DESCRIPTION:Test Profile 1:DESCRIPTION ###
`
	profile2 := `
### SHELL:pwsh:SHELL ###
### DESCRIPTION:Test Profile 2:DESCRIPTION ###
`
	profile1Path := filepath.Join(tmpDir, "profile1.Profile.ps1")
	profile2Path := filepath.Join(tmpDir, "profile2.Profile.ps1")
	assert.NoError(t, os.WriteFile(profile1Path, []byte(profile1), 0644))
	assert.NoError(t, os.WriteFile(profile2Path, []byte(profile2), 0644))

	// Mock exec.LookPath
	originalLookPath := lookPath
	defer func() { lookPath = originalLookPath }()
	lookPath = func(file string) (string, error) {
		return "/path/to/" + file, nil
	}

	// Mock launcher.ExecutePowerShellProcess
	originalExecutePowerShellProcess := executePowerShellProcess
	defer func() { executePowerShellProcess = originalExecutePowerShellProcess }()
	var executedProfile string
	executePowerShellProcess = func(profileContent, shellPath string) error {
		executedProfile = profileContent
		return nil
	}

	// Test launching a single profile
	err = LaunchProfilesFromCmd(profile1Path, "powershell")
	assert.NoError(t, err)
	assert.Contains(t, executedProfile, "Test Profile 1")

	// Test launching multiple profiles
	err = LaunchProfilesFromCmd(profile1Path+","+profile2Path, "powershell")
	assert.NoError(t, err)
	assert.Contains(t, executedProfile, "Test Profile 1")
	assert.NotContains(t, executedProfile, "Test Profile 2")

	// Test launching with an invalid shell
	err = LaunchProfilesFromCmd(profile1Path, "invalid-shell")
	assert.Error(t, err)
}

func TestValidateProfile(t *testing.T) {
	originalValidatePath := ValidatePath
	originalValidateShellVersion := ValidateShellVersion
	originalValidateDescription := ValidateDescription
	defer func() {
		ValidatePath = originalValidatePath
		ValidateShellVersion = originalValidateShellVersion
		ValidateDescription = originalValidateDescription
	}()

	// Test valid profile
	p1 := &types.ProfileItem{
		Path:            "/tmp/test.ps1",
		Shell:           "PowerShell",
		ItemDescription: "A valid profile",
	}
	ValidatePath = func(path string) (bool, error) { return true, nil }
	ValidateShellVersion = func(shell string) (bool, error) { return true, nil }
	ValidateDescription = func(desc string) (bool, error) { return true, nil }
	validateProfile(p1)
	assert.True(t, p1.IsValid)

	// Test invalid profile
	p2 := &types.ProfileItem{
		Path:            "/tmp/test.ps1",
		Shell:           "InvalidShell",
		ItemDescription: "",
	}
	ValidatePath = func(path string) (bool, error) { return true, nil }
	ValidateShellVersion = func(shell string) (bool, error) { return false, nil }
	ValidateDescription = func(desc string) (bool, error) { return false, nil }
	validateProfile(p2)
	assert.False(t, p2.IsValid)
}
