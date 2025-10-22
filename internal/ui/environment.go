package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"litepost/internal/models"
)

// EnvironmentModel represents the environment manager UI
type EnvironmentModel struct {
	environments []*models.Environment
	selected     int
	width        int
	height       int
}

// NewEnvironmentModel creates a new environment model
func NewEnvironmentModel() *EnvironmentModel {
	return &EnvironmentModel{
		environments: []*models.Environment{
			{
				ID:        "1",
				Name:      "Development",
				Variables: map[string]string{"base_url": "http://localhost:3000"},
				IsActive:  true,
			},
			{
				ID:        "2",
				Name:      "Production",
				Variables: map[string]string{"base_url": "https://api.example.com"},
				IsActive:  false,
			},
		},
		selected: 0,
	}
}

// Init initializes the environment model
func (e *EnvironmentModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the environment model
func (e *EnvironmentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return e, nil
		case "up", "k":
			if e.selected > 0 {
				e.selected--
			}
		case "down", "j":
			if e.selected < len(e.environments)-1 {
				e.selected++
			}
		case "enter":
			// TODO: Implement environment selection
		case "n":
			// TODO: Implement new environment creation
		}
	}

	return e, nil
}

// View renders the environment model
func (e *EnvironmentModel) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Environments")

	var items []string
	for i, env := range e.environments {
		style := lipgloss.NewStyle()
		if i == e.selected {
			style = style.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
		}
		
		status := "Inactive"
		if env.IsActive {
			status = "Active"
		}
		
		item := style.Render(fmt.Sprintf("%s (%s) - %d variables", env.Name, status, len(env.Variables)))
		items = append(items, item)
	}

	content := strings.Join(items, "\n")
	if len(items) == 0 {
		content = "No environments found"
	}

	menu := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2).
		Render(content)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Use arrow keys to navigate, Enter to select, 'n' for new, Esc to go back")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		menu,
		"",
		help,
	)
}
