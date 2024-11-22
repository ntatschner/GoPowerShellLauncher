package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	l "github.com/ntatschner/GoPowerShellLauncher/cmd/logger"
)

type SwitchViewMsg string

type KeyMap struct {
	// Keybindings used when browsing the list.
	CursorUp    key.Binding
	CursorDown  key.Binding
	NextPage    key.Binding
	PrevPage    key.Binding
	GoToStart   key.Binding
	GoToEnd     key.Binding
	Filter      key.Binding
	ClearFilter key.Binding

	// Keybindings used when setting a filter.
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding

	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding

	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding

	// The quit-no-matter-what keybinding. This will be caught when filtering.
	ForceQuit key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Browsing.
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left", "h", "pgup", "b", "u"),
			key.WithHelp("←/h/pgup", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right", "l", "pgdown", "f", "d"),
			key.WithHelp("→/l/pgdn", "next page"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear filter"),
		),

		// Filtering.
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
			key.WithHelp("enter", "apply filter"),
		),

		// Toggle help.
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		// Quitting.
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

func validatePath(path string) error {
	l.Logger.Info("Validating path", "Path", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}
	if err != nil {
		return fmt.Errorf("error accessing path: %s", err)
	}
	return nil
}

func validateHash(hash string, filePath string) error {
	l.Logger.Info("Validating hash", "Hash", hash, "Path", filePath)
	// Check if the path is valid before validating the hash
	if err := validatePath(filePath); err != nil {
		return err
	}

	// Decode the base64-encoded hash
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return fmt.Errorf("invalid base64 hash: %v", err)
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Compute the SHA-256 hash of the file
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("error hashing file: %v", err)
	}
	computedHash := hasher.Sum(nil)

	// Compare the computed hash with the provided hash
	if !compareHashes(decodedHash, computedHash) {
		return fmt.Errorf("hash mismatch: expected %x, got %x", decodedHash, computedHash)
	}

	return nil
}

func compareHashes(hash1, hash2 []byte) bool {
	l.Logger.Info("Comparing hashes", "Hash1", hash1, "Hash2", hash2)
	if len(hash1) != len(hash2) {
		return false
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}

func validateShellVersion(shellVersion string) error {
	l.Logger.Info("Validating shell version", "ShellVersion", shellVersion)
	switch shellVersion {
	case "pwsh", "powershell":
		return nil
	default:
		return fmt.Errorf("invalid shell version: %s", shellVersion)
	}
}

func validateDescription(description string) error {
	l.Logger.Info("Validating description", "Description", description)
	if len(description) > 100 {
		return fmt.Errorf("description is too long (max 100 characters)")
	}
	return nil
}

func getProfileContent(path string) (string, error) {
	l.Logger.Info("Getting profile content", "Path", path)
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

func MergeSelectedProfiles(selected map[int]struct{}) string {
	l.Logger.Info("Merging selected profiles", "Selected", selected)
	var merged string
	for i := range selected {
		content, err := getProfileContent(profiles[i].path)
		if err != nil {
			l.Logger.Warn("Error reading profile content", "Error", err)
			continue
		}
		merged += content + "\n"
	}
	return merged
}

func loadAvailableShells() []list.Item {
	l.Logger.Info("Loading available shells")
	var items []list.Item
	for _, shell := range shells {
		items = append(items, shell)
	}
	return items
}
