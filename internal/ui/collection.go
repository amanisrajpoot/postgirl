package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"postgirl/internal/models"
)

// CollectionModel represents the collection manager UI
type CollectionModel struct {
	collections []*models.Collection
	selected    int
	width       int
	height      int
}

// NewCollectionModel creates a new collection model
func NewCollectionModel() *CollectionModel {
	return &CollectionModel{
		collections: []*models.Collection{
			{
				ID:          "1",
				Name:        "Sample Collection",
				Description: "A sample collection for testing",
				Requests:    []models.Request{},
				Folders:    []models.Folder{},
				Variables:   make(map[string]string),
			},
		},
		selected: 0,
	}
}

// Init initializes the collection model
func (c *CollectionModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the collection model
func (c *CollectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return c, nil
		case "up", "k":
			if c.selected > 0 {
				c.selected--
			}
		case "down", "j":
			if c.selected < len(c.collections)-1 {
				c.selected++
			}
		case "enter":
			// TODO: Implement collection selection
		case "n":
			// TODO: Implement new collection creation
		}
	}

	return c, nil
}

// View renders the collection model
func (c *CollectionModel) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Collections")

	var items []string
	for i, collection := range c.collections {
		style := lipgloss.NewStyle()
		if i == c.selected {
			style = style.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
		}
		item := style.Render(fmt.Sprintf("%s - %s", collection.Name, collection.Description))
		items = append(items, item)
	}

	content := strings.Join(items, "\n")
	if len(items) == 0 {
		content = "No collections found"
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
