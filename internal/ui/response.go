package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"litepost/internal/models"
)

// ResponseModel represents the response viewer UI
type ResponseModel struct {
	response *models.Response
	selected int
	width    int
	height   int
}

// NewResponseModel creates a new response model
func NewResponseModel() *ResponseModel {
	return &ResponseModel{
		response: &models.Response{
			ID:         "new",
			RequestID:  "new",
			StatusCode: 200,
			Headers:    make(map[string]string),
			Body:       "No response yet",
			Size:       0,
			Duration:   0,
		},
		selected: 0,
	}
}

// Init initializes the response model
func (r *ResponseModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the response model
func (r *ResponseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return r, nil
		case "up", "k":
			if r.selected > 0 {
				r.selected--
			}
		case "down", "j":
			if r.selected < 2 {
				r.selected++
			}
		}
	}

	return r, nil
}

// View renders the response model
func (r *ResponseModel) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Response Viewer")

	// Status code
	_ = lipgloss.Color("#00FF00") // statusColor for future use

	statusStyle := lipgloss.NewStyle()
	if r.selected == 0 {
		statusStyle = statusStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	statusText := statusStyle.Render(fmt.Sprintf("Status: %d", r.response.StatusCode))

	// Headers
	headersStyle := lipgloss.NewStyle()
	if r.selected == 1 {
		headersStyle = headersStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	headersText := headersStyle.Render("Headers: (Press Enter to view)")

	// Body
	bodyStyle := lipgloss.NewStyle()
	if r.selected == 2 {
		bodyStyle = bodyStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	bodyText := bodyStyle.Render(fmt.Sprintf("Body: %s", r.response.Body))

	content := strings.Join([]string{
		statusText,
		headersText,
		bodyText,
	}, "\n")

	menu := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2).
		Render(content)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Use arrow keys to navigate, Enter to select, Esc to go back")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		menu,
		"",
		help,
	)
}
