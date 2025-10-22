package ui

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"postgirl/internal/app"
	"postgirl/internal/storage/sqlite"
)

// AppState represents the current state of the application
type AppState int

const (
	StateMain AppState = iota
	StateRequest
	StateResponse
	StateCollection
	StateEnvironment
)

// App represents the main application
type App struct {
	state       AppState
	request     *RequestModel
	response    *ResponseModel
	collection  *CollectionModel
	environment *EnvironmentModel
	width       int
	height      int
	service     *app.Service
}

// NewApp creates a new application instance
func NewApp() *App {
	// Initialize storage
	storage, err := sqlite.NewSQLiteStorage("postgirl.db")
	if err != nil {
		// For now, continue without storage
		storage = nil
	}
	
	// Initialize service
	service := app.NewService(storage)
	
	return &App{
		state:       StateMain,
		request:     NewRequestModel(service),
		response:    NewResponseModel(),
		collection:  NewCollectionModel(),
		environment: NewEnvironmentModel(),
		service:     service,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.request.Init(),
		a.response.Init(),
		a.collection.Init(),
		a.environment.Init(),
	)
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "1":
			a.state = StateRequest
		case "2":
			a.state = StateResponse
		case "3":
			a.state = StateCollection
		case "4":
			a.state = StateEnvironment
		case "esc":
			a.state = StateMain
		}
	}

	// Update current model based on state
	switch a.state {
	case StateRequest:
		model, cmd := a.request.Update(msg)
		a.request = model.(*RequestModel)
		return a, cmd
	case StateResponse:
		model, cmd := a.response.Update(msg)
		a.response = model.(*ResponseModel)
		return a, cmd
	case StateCollection:
		model, cmd := a.collection.Update(msg)
		a.collection = model.(*CollectionModel)
		return a, cmd
	case StateEnvironment:
		model, cmd := a.environment.Update(msg)
		a.environment = model.(*EnvironmentModel)
		return a, cmd
	}

	return a, nil
}

// View renders the application
func (a *App) View() string {
	switch a.state {
	case StateMain:
		return a.mainView()
	case StateRequest:
		return a.request.View()
	case StateResponse:
		return a.response.View()
	case StateCollection:
		return a.collection.View()
	case StateEnvironment:
		return a.environment.View()
	default:
		return "Unknown state"
	}
}

// mainView renders the main menu
func (a *App) mainView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Litepost")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A8A8A8")).
		Render("Lightweight Postman Alternative")

	menu := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(1, 2).
		Render(
			strings.Join([]string{
				"1. Request Builder",
				"2. Response Viewer",
				"3. Collections",
				"4. Environments",
				"",
				"Press 'q' to quit",
			}, "\n"),
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Render("Press a number to navigate, 'q' to quit")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		subtitle,
		"",
		menu,
		"",
		help,
	)
}
