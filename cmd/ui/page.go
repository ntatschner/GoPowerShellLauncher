package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ntatschner/GoPowerShellLauncher/cmd/ui/view"
)

type Page interface {
	New(view.ViewChanger, tea.WindowSizeMsg) tea.Model
	Title() string
	Description() string
}

var pageRegistry = make(map[string]Page)

func RegisterPage(name string, page Page) {
	pageRegistry[name] = page
}

func GetPage(name string) Page {
	return pageRegistry[name]
}

func GetPages() map[string]Page {
	return pageRegistry
}
