package common

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var WindowSize tea.WindowSizeMsg

type KeyMap struct {
	Back key.Binding
}

var Keymap = KeyMap{
	Back: key.NewBinding(
		key.WithKeys("alt+left"),
		key.WithHelp("alt+left", "back"),
	),
}
