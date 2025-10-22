package ui

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputModel represents a text input component
type InputModel struct {
	value    string
	placeholder string
	focused  bool
	width    int
}

// NewInputModel creates a new input model
func NewInputModel(placeholder string) *InputModel {
	return &InputModel{
		value:       "",
		placeholder: placeholder,
		focused:     false,
		width:       50,
	}
}

// Init initializes the input model
func (i *InputModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the input model
func (i *InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return i, tea.Quit
		case "enter":
			return i, tea.Quit
		case "backspace":
			if len(i.value) > 0 {
				i.value = i.value[:len(i.value)-1]
			}
		default:
			if len(msg.String()) == 1 {
				i.value += msg.String()
			}
		}
	}

	return i, nil
}

// View renders the input model
func (i *InputModel) View() string {
	displayValue := i.value
	if displayValue == "" {
		displayValue = i.placeholder
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(i.width)

	if i.focused {
		style = style.BorderForeground(lipgloss.Color("#FAFAFA"))
	}

	return style.Render(displayValue)
}

// Value returns the current input value
func (i *InputModel) Value() string {
	return i.value
}

// SetValue sets the input value
func (i *InputModel) SetValue(value string) {
	i.value = value
}

// Focus sets the input as focused
func (i *InputModel) Focus() {
	i.focused = true
}

// Blur removes focus from the input
func (i *InputModel) Blur() {
	i.focused = false
}
