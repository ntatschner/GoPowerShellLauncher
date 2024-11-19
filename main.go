package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var logger = log.New("app")

// ValidateProfileFile checks if the given PowerShell profile file is valid.
func ValidateProfileFile(fileBytes []byte) error {
	h := sha256.New()
	h.Write(fileBytes)
	return nil // Replace with your validation logic here!
}

type shellClosed struct{}

type Shell struct {
	Name string
	Path string
}

type profile struct {
	profile_path string
	hash         string
	ps_version   string
}

type Config struct {
	CsvPath string `json:"csv_path"`
}

var profiles []profile
var shells = []Shell{
	{
		Name: "Pwsh",
		Path: "C:\\Program Files\\PowerShell\\7\\pwsh.exe",
	},
	{
		Name: "Powershell",
		Path: "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
	},
}

func loadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadProfiles(filePath string) error {
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

	for _, record := range records[1:] { // Skip header row
		profiles = append(profiles, profile{
			profile_path: record[0],
			hash:         record[1],
			ps_version:   record[2],
		})
	}

	return nil
}

func getProfilePath(hash string) string {
	for _, p := range profiles {
		if p.hash == hash {
			return p.profile_path
		}
	}
	return ""
}

func getProfileContent(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func mergeSelectedProfiles(selected map[int]struct{}) string {
	var merged string
	for i := range selected {
		content, err := getProfileContent(profiles[i].profile_path)
		if err != nil {
			logger.Errorf("Error reading profile content", "Error", err)
			continue
		}
		merged += content + "\n"
	}
	return merged
}

func launchPowerShell(m model) tea.Cmd {
	return func() tea.Msg {
		merged := mergeSelectedProfiles(m.selected)
		// create temporary file to store the merged profile
		tmpFile, err := os.CreateTemp("", "merged_profile_*.ps1")
		if err != nil {
			logger.Errorf("Error creating temporary file", "Error", err)
			return shellClosed{}
		}
		_, err = tmpFile.Write([]byte(merged))
		if err != nil {
			logger.Errorf("Error writing to temporary file", "Error", err)
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return shellClosed{}
		}
		tmpFile.Close() // Close the file before running the PowerShell command

		cmd := exec.Command("cmd", "/C", "start", "/wait", "powershell", "-NoExit", "-File", tmpFile.Name())
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}

		err = cmd.Start()
		if err != nil {
			logger.Errorf("Error starting PowerShell command", "Error", err)
			os.Remove(tmpFile.Name())
			return shellClosed{}
		}
		cmd.Wait()
		// Remove the temporary file after starting the command
		os.Remove(tmpFile.Name())
		return shellClosed{}
	}
}

//////////////////////////////

// Var

var (
	// Styles
	profileNameStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	profileDescriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	profileVersionStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF"))
	profileValidStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	profileInvalidStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	profileSelectedStyle    = lipgloss.NewStyle().Italic(true)
	profileCursorStyle      = lipgloss.NewStyle().Underline(true)
)

// Main Model

type model struct {
	selected     map[int]struct{}
	cursor       int
	profiles     []profile
	defaultStyle lipgloss.Style
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			return m, launchPowerShell(m)
		}
	}
	return m, nil
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Rothesay IT - PowerShell Profile Selector")
}

func initialModel() model {
	config, err := loadConfig("config.json")
	if err != nil {
		logger.Fatal(err)
	}

	err = loadProfiles(config.CsvPath)
	if err != nil {
		logger.Fatal(err)
	}

	return model{
		selected: make(map[int]struct{}),
		profiles: profiles,
	}
}

func (m model) View() string {
	s := "Available PowerShell profiles:\n"

	for i, p := range m.profiles {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		selected := "[ ]"
		if _, ok := m.selected[i]; ok {
			selected = "[x]"
		}
		s += fmt.Sprintf("%s %s %s\n", cursor, selected, p.profile_path)
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)

	}
}
