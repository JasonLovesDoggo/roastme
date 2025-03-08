package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonlovesdoggo/roastme/internal/config"
	"strings"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5F87")).
			MarginBottom(1).
			MarginTop(1)

	roastBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF875F")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1).
			Width(70)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5F87FF"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5FFF87")).
			Bold(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5F87FF"))
)

// Spinner represents a loading spinner
type Spinner struct {
	model  spinner.Model
	text   string
	active bool
}

// NewSpinner creates a new spinner with the given text
func NewSpinner(text string) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87"))
	return &Spinner{
		model:  s,
		text:   text,
		active: false,
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.active = true
	go func() {
		p := tea.NewProgram(spinnerModel{s.model, s.text})
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running spinner: %v\n", err)
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.active = false
}

// spinnerModel is the model for the spinner
type spinnerModel struct {
	spinner spinner.Model
	text    string
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m spinnerModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.text)
}

// DisplayRoast displays the roast in a pretty format
func DisplayRoast(roast string) {
	// Title
	fmt.Println(titleStyle.Render("🔥 COMMAND HISTORY ROAST 🔥"))

	// Roast content in a fancy box
	fmt.Println(roastBoxStyle.Render(strings.TrimSpace(roast)))

	// Footer
	fmt.Println(infoStyle.Render("Run with --deep for a deeper analysis or use 'config' to set up AI integration"))
}

// configModel represents the state for the config UI
type configModel struct {
	inputs        []textinput.Model
	currentIndex  int
	providerIndex int
	saved         bool
	cfg           config.Config
	err           string
}

func (m configModel) Init() tea.Cmd {
	return textinput.Blink
}
