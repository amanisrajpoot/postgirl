package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"postgirl/internal/app"
	"postgirl/internal/models"
)

// RequestModel represents the request builder UI
type RequestModel struct {
	request   *models.Request
	method    string
	url       string
	headers   map[string]string
	body      string
	bodyType  string
	selected  int
	width     int
	height    int
	service   *app.Service
	loading   bool
	response  *models.Response
	error     string
	urlInput  *InputModel
	inputMode bool
}

// NewRequestModel creates a new request model
func NewRequestModel(service *app.Service) *RequestModel {
	req := service.CreateNewRequest()
	return &RequestModel{
		request:   req,
		method:    req.Method,
		url:       req.URL,
		headers:   req.Headers,
		body:      "",
		bodyType:  "json",
		selected:  0,
		service:   service,
		loading:   false,
		response:  nil,
		error:     "",
		urlInput:  NewInputModel("Enter URL (e.g., https://httpbin.org/get)"),
		inputMode: false,
	}
}

// Init initializes the request model
func (r *RequestModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the request model
func (r *RequestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle input mode
	if r.inputMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				r.inputMode = false
				r.urlInput.Blur()
				return r, nil
			case "enter":
				r.url = r.urlInput.Value()
				r.inputMode = false
				r.urlInput.Blur()
				r.updateRequest()
				return r, nil
			}
		}
		
		// Update the input model
		model, cmd := r.urlInput.Update(msg)
		r.urlInput = model.(*InputModel)
		return r, cmd
	}

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
			if r.selected < 4 {
				r.selected++
			}
		case "enter":
			switch r.selected {
			case 0: // Method
				r.method = r.cycleMethod()
				r.updateRequest()
			case 1: // URL
				r.inputMode = true
				r.urlInput.SetValue(r.url)
				r.urlInput.Focus()
			case 2: // Headers
				// TODO: Implement header editing
			case 3: // Body
				// TODO: Implement body editing
			case 4: // Send
				if !r.loading {
					return r, r.sendRequest()
				}
			}
		}
	case RequestSentMsg:
		if msg.Response != nil {
			r.response = msg.Response
			r.error = ""
		} else {
			r.error = msg.Error
		}
		r.loading = false
	}

	return r, nil
}

// View renders the request model
func (r *RequestModel) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Request Builder")

	// Method selection
	methodStyle := lipgloss.NewStyle()
	if r.selected == 0 {
		methodStyle = methodStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	methodText := methodStyle.Render(fmt.Sprintf("Method: %s", r.method))

	// URL input
	urlStyle := lipgloss.NewStyle()
	if r.selected == 1 {
		urlStyle = urlStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	
	var urlText string
	if r.inputMode {
		urlText = urlStyle.Render("URL: " + r.urlInput.View())
	} else {
		urlText = urlStyle.Render(fmt.Sprintf("URL: %s", r.url))
	}

	// Headers
	headersStyle := lipgloss.NewStyle()
	if r.selected == 2 {
		headersStyle = headersStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	headersText := headersStyle.Render("Headers: (Press Enter to edit)")

	// Body
	bodyStyle := lipgloss.NewStyle()
	if r.selected == 3 {
		bodyStyle = bodyStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	bodyText := bodyStyle.Render(fmt.Sprintf("Body (%s): %s", r.bodyType, r.body))

	// Send button
	sendStyle := lipgloss.NewStyle()
	if r.selected == 4 {
		sendStyle = sendStyle.Bold(true).Foreground(lipgloss.Color("#7D56F4"))
	}
	
	sendText := "Send Request"
	if r.loading {
		sendText = "Sending..."
	}
	sendText = sendStyle.Render(sendText)

	content := strings.Join([]string{
		methodText,
		urlText,
		headersText,
		bodyText,
		sendText,
	}, "\n")

	// Add response/error display
	if r.response != nil {
		bodyPreview := r.response.Body
		if len(bodyPreview) > 100 {
			bodyPreview = bodyPreview[:100] + "..."
		}
		responseText := fmt.Sprintf("\n\nResponse: %d - %s", r.response.StatusCode, bodyPreview)
		content += responseText
	}
	
	if r.error != "" {
		errorText := fmt.Sprintf("\n\nError: %s", r.error)
		content += errorText
	}

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

// cycleMethod cycles through HTTP methods
func (r *RequestModel) cycleMethod() string {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for i, method := range methods {
		if method == r.method {
			next := (i + 1) % len(methods)
			return methods[next]
		}
	}
	return "GET"
}

// sendRequest sends the HTTP request
func (r *RequestModel) sendRequest() tea.Cmd {
	r.loading = true
	r.error = ""
	
	return func() tea.Msg {
		// Update request with current values
		r.updateRequest()
		
		// Execute the request
		response, err := r.service.ExecuteRequest(r.request)
		
		if err != nil {
			return RequestSentMsg{
				Request:  r.request,
				Response: nil,
				Error:    err.Error(),
			}
		}
		
		return RequestSentMsg{
			Request:  r.request,
			Response: response,
			Error:    "",
		}
	}
}

// RequestSentMsg represents a message when a request is sent
type RequestSentMsg struct {
	Request  *models.Request
	Response *models.Response
	Error    string
}

// updateRequest updates the request model with current values
func (r *RequestModel) updateRequest() {
	r.request.Method = r.method
	r.request.URL = r.url
	r.request.Headers = r.headers
	
	if r.body != "" {
		r.request.Body = &models.RequestBody{
			Type:    r.bodyType,
			Content: r.body,
		}
	}
}
