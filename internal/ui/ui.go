package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jasonlovesdoggo/roastme/internal/config"
)

var (
	// Styles - more professional color scheme
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#61AFEF")).
			MarginBottom(1).
			MarginTop(1)

	roastBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#61AFEF")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1).
			Width(70)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E06C75")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#56B6C2"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")).
			Bold(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C678DD"))

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5C07B")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")).
			Bold(true)

	radioUnselectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ABB2BF"))
)

// Spinner represents a loading spinner
type Spinner struct {
	model  spinner.Model
	text   string
	active bool
	done   chan struct{}
}

// NewSpinner creates a new spinner with the given text
func NewSpinner(text string) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF"))
	return &Spinner{
		model:  s,
		text:   text,
		active: false,
		done:   make(chan struct{}),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.active = true
	go func() {
		p := tea.NewProgram(spinnerModel{s.model, s.text})
		go func() {
			<-s.done
			p.Quit()
		}()
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running spinner: %v\n", err)
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	if s.active {
		s.active = false
		close(s.done)
		// Clear the spinner line
		fmt.Print("\033[2K\r")
	}
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
		if msg.String() == "q" || msg.String() == "ctrl+c" {
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
	// Clear screen
	fmt.Print("\033[H\033[2J")

	// Clean professional header (no emoji)
	fmt.Println(titleStyle.Render("───── COMMAND HISTORY ROAST ─────"))
	fmt.Println()

	// Roast content in a fancy box
	fmt.Println(roastBoxStyle.Render(strings.TrimSpace(roast)))
	fmt.Println()

	// Footer
	fmt.Println(infoStyle.Render("Run with --deep for a deeper analysis"))
	fmt.Println(infoStyle.Render("Use 'roastme config' to customize AI settings"))
}

// Msg types for our program
type providerSelectMsg string

// configModel represents the state for the config UI
type configModel struct {
	inputs          []textinput.Model
	currentIndex    int
	currentProvider string
	saved           bool
	cfg             config.Config
	err             string
	successMsg      string
	page            string // "provider" or "settings"
}

func (m configModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+v": // Handle paste explicitly
			// Continue to standard handling which will include paste

		case "enter":
			if m.page == "provider" {
				// Switch to settings page for selected provider
				m.page = "settings"
				m.setupInputsForCurrentProvider()
				m.currentIndex = 0
				if len(m.inputs) > 0 {
					m.inputs[m.currentIndex].Focus()
				}
				return m, nil

			} else if m.page == "settings" {
				if m.currentIndex == len(m.inputs)-1 {
					// Save button
					m.saveConfig()
					return m, nil
				} else if m.currentIndex == len(m.inputs)-2 {
					// Back button
					m.page = "provider"
					return m, nil
				}

				// Move to next field
				m.currentIndex++
				if m.currentIndex >= len(m.inputs) {
					m.currentIndex = 0
				}

				for i := range m.inputs {
					if i == m.currentIndex {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			}
			return m, nil

		case "tab":
			if m.page == "provider" {
				// Cycle through providers
				providers := []string{"local", "gemini", "openai", "anthropic", "custom"}
				for i, p := range providers {
					if p == m.currentProvider {
						nextIndex := (i + 1) % len(providers)
						m.currentProvider = providers[nextIndex]
						break
					}
				}
				return m, nil
			} else if m.page == "settings" {
				// Move to next input
				m.currentIndex = (m.currentIndex + 1) % len(m.inputs)
				for i := range m.inputs {
					if i == m.currentIndex {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			}
			return m, nil

		case "shift+tab":
			if m.page == "provider" {
				// Cycle through providers backward
				providers := []string{"local", "gemini", "openai", "anthropic", "custom"}
				for i, p := range providers {
					if p == m.currentProvider {
						nextIndex := (i - 1)
						if nextIndex < 0 {
							nextIndex = len(providers) - 1
						}
						m.currentProvider = providers[nextIndex]
						break
					}
				}
				return m, nil
			} else if m.page == "settings" {
				// Move to previous input
				m.currentIndex--
				if m.currentIndex < 0 {
					m.currentIndex = len(m.inputs) - 1
				}

				for i := range m.inputs {
					if i == m.currentIndex {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
			}
			return m, nil

		case "1", "2", "3", "4", "5":
			if m.page == "provider" {
				// Direct provider selection using number keys
				providers := []string{"local", "gemini", "openai", "anthropic", "custom"}
				idx := int(msg.Runes[0] - '1')
				if idx >= 0 && idx < len(providers) {
					m.currentProvider = providers[idx]
				}
			}
			return m, nil
		}

	case providerSelectMsg:
		// Handle provider selection message
		m.currentProvider = string(msg)
		return m, nil
	}

	// Only update inputs on settings page
	if m.page == "settings" && m.currentIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.currentIndex], cmd = m.inputs[m.currentIndex].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m configModel) View() string {
	// Clear screen
	//s := "\033[H\033[2J"

	if m.page == "provider" {
		return m.renderProviderSelection()
	} else {
		return m.renderSettingsPage()
	}
}

func (m configModel) renderProviderSelection() string {
	s := titleStyle.Render("ROASTME CONFIGURATION") + "\n\n"
	s += promptStyle.Render("Select AI Provider:") + "\n\n"

	providers := []string{"local", "gemini", "openai", "anthropic", "custom"}
	descriptions := map[string]string{
		"local":     "No API key required - uses built-in rules",
		"gemini":    "Google's Gemini AI (API key required)",
		"openai":    "OpenAI's GPT models (API key required)",
		"anthropic": "Anthropic's Claude models (API key required)",
		"custom":    "Custom LLM provider (API key required)",
	}

	for i, provider := range providers {
		prefix := "   "
		if provider == m.currentProvider {
			prefix = " > "
			s += highlightStyle.Render(fmt.Sprintf("%s(%d) [%s] %s\n", prefix, i+1, selectedStyle.Render("●"), provider))
			s += "      " + infoStyle.Render(descriptions[provider]) + "\n\n"
		} else {
			s += fmt.Sprintf("%s(%d) [%s] %s\n", prefix, i+1, radioUnselectedStyle.Render("○"), provider)
			s += "      " + descriptions[provider] + "\n\n"
		}
	}

	s += "\n" + promptStyle.Render("Press Enter to configure, Tab to navigate, or Esc to quit") + "\n"

	if m.err != "" {
		s += "\n" + errorStyle.Render(m.err)
	}

	if m.successMsg != "" {
		s += "\n" + successStyle.Render(m.successMsg)
	}

	return s
}

func (m configModel) renderSettingsPage() string {
	s := titleStyle.Render("CONFIGURE "+strings.ToUpper(m.currentProvider)) + "\n\n"

	if m.currentProvider == "local" {
		s += infoStyle.Render("Local mode doesn't require any API keys or configuration.\n")
		s += infoStyle.Render("Roasts are generated using built-in rules.") + "\n\n"
	} else {
		// Display the inputs
		for i, input := range m.inputs {
			// Skip buttons for now
			if i >= len(m.inputs)-2 {
				continue
			}

			label := ""
			switch m.currentProvider {
			case "gemini":
				if i == 0 {
					label = "API Key:"
				} else if i == 1 {
					label = "Model:"
				}
			case "openai":
				if i == 0 {
					label = "API Key:"
				} else if i == 1 {
					label = "Model:"
				}
			case "anthropic":
				if i == 0 {
					label = "API Key:"
				} else if i == 1 {
					label = "Model:"
				}
			case "custom":
				if i == 0 {
					label = "API Key:"
				} else if i == 1 {
					label = "Base URL:"
				} else if i == 2 {
					label = "Model:"
				}
			}

			// Add cursor indicator for focused input
			prefix := "  "
			if i == m.currentIndex {
				prefix = highlightStyle.Render("> ")
			}

			s += fmt.Sprintf("%s%s %s\n", prefix, promptStyle.Render(label), input.View())
			s += "\n"
		}
	}

	// Add Back and Save buttons
	for i := len(m.inputs) - 2; i < len(m.inputs); i++ {
		if i < 0 || i >= len(m.inputs) {
			continue
		}

		prefix := "  "
		if i == m.currentIndex {
			prefix = highlightStyle.Render("> ")
		}

		s += fmt.Sprintf("%s%s\n", prefix, m.inputs[i].View())
	}

	s += "\n" + promptStyle.Render("Tab/Shift+Tab to navigate, Enter to select") + "\n"

	if m.err != "" {
		s += "\n" + errorStyle.Render(m.err)
	}

	if m.successMsg != "" {
		s += "\n" + successStyle.Render(m.successMsg)
	}

	return s
}

// RunConfigInterface runs the configuration interface
func RunConfigInterface() {
	cfg := config.GetConfig()

	// Create initial model for provider selection
	model := configModel{
		inputs:          []textinput.Model{},
		cfg:             cfg,
		currentProvider: cfg.AI.Provider,
		page:            "provider", // Start with provider selection
	}

	// Run the program
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running config interface: %v\n", err)
		os.Exit(1)
	}
}

// Set up inputs for current provider
func (m *configModel) setupInputsForCurrentProvider() {
	// Clear existing inputs
	m.inputs = []textinput.Model{}

	switch m.currentProvider {
	case "gemini":
		// API Key
		apiKeyInput := textinput.New()
		apiKeyInput.Placeholder = "Enter your Gemini API Key"
		apiKeyInput.SetValue(m.cfg.AI.Gemini.APIKey)
		apiKeyInput.Width = 60
		apiKeyInput.EchoMode = textinput.EchoPassword
		apiKeyInput.EchoCharacter = '•'

		// Model
		modelInput := textinput.New()
		modelInput.Placeholder = "gemini-pro"
		modelInput.SetValue(m.cfg.AI.Gemini.Model)
		if modelInput.Value() == "" {
			modelInput.SetValue("gemini-pro")
		}

		m.inputs = append(m.inputs, apiKeyInput, modelInput)

	case "openai":
		// API Key
		apiKeyInput := textinput.New()
		apiKeyInput.Placeholder = "Enter your OpenAI API Key"
		apiKeyInput.SetValue(m.cfg.AI.OpenAI.APIKey)
		apiKeyInput.Width = 60
		apiKeyInput.EchoMode = textinput.EchoPassword
		apiKeyInput.EchoCharacter = '•'

		// Model
		modelInput := textinput.New()
		modelInput.Placeholder = "gpt-3.5-turbo"
		modelInput.SetValue(m.cfg.AI.OpenAI.Model)
		if modelInput.Value() == "" {
			modelInput.SetValue("gpt-3.5-turbo")
		}

		m.inputs = append(m.inputs, apiKeyInput, modelInput)

	case "anthropic":
		// API Key
		apiKeyInput := textinput.New()
		apiKeyInput.Placeholder = "Enter your Anthropic API Key"
		apiKeyInput.SetValue(m.cfg.AI.Anthropic.APIKey)
		apiKeyInput.Width = 60
		apiKeyInput.EchoMode = textinput.EchoPassword
		apiKeyInput.EchoCharacter = '•'

		// Model
		modelInput := textinput.New()
		modelInput.Placeholder = "claude-2"
		modelInput.SetValue(m.cfg.AI.Anthropic.Model)
		if modelInput.Value() == "" {
			modelInput.SetValue("claude-2")
		}

		m.inputs = append(m.inputs, apiKeyInput, modelInput)

	case "custom":
		// API Key
		apiKeyInput := textinput.New()
		apiKeyInput.Placeholder = "Enter your API Key"
		apiKeyInput.SetValue(m.cfg.AI.Custom.APIKey)
		apiKeyInput.Width = 60
		apiKeyInput.EchoMode = textinput.EchoPassword
		apiKeyInput.EchoCharacter = '•'

		// Base URL
		baseURLInput := textinput.New()
		baseURLInput.Placeholder = "https://api.example.com"
		baseURLInput.SetValue(m.cfg.AI.Custom.BaseURL)
		baseURLInput.Width = 60

		// Model
		modelInput := textinput.New()
		modelInput.Placeholder = "model-name"
		modelInput.SetValue(m.cfg.AI.Custom.Model)

		m.inputs = append(m.inputs, apiKeyInput, baseURLInput, modelInput)

	case "local":
		// No inputs needed for local
	}

	// Add Back and Save buttons
	backButton := textinput.New()
	backButton.SetValue("[ Back to Provider Selection ]")
	backButton.Blur()

	saveButton := textinput.New()
	saveButton.SetValue("[ Save Configuration ]")
	saveButton.Blur()

	m.inputs = append(m.inputs, backButton, saveButton)
}

// Save the configuration
func (m *configModel) saveConfig() {
	// Update config based on current provider
	m.cfg.AI.Provider = m.currentProvider

	// Update provider-specific settings
	if m.currentProvider != "local" && len(m.inputs) >= 2 {
		switch m.currentProvider {
		case "gemini":
			m.cfg.AI.Gemini.APIKey = m.inputs[0].Value()
			m.cfg.AI.Gemini.Model = m.inputs[1].Value()
		case "openai":
			m.cfg.AI.OpenAI.APIKey = m.inputs[0].Value()
			m.cfg.AI.OpenAI.Model = m.inputs[1].Value()
		case "anthropic":
			m.cfg.AI.Anthropic.APIKey = m.inputs[0].Value()
			m.cfg.AI.Anthropic.Model = m.inputs[1].Value()
		case "custom":
			if len(m.inputs) >= 3 {
				m.cfg.AI.Custom.APIKey = m.inputs[0].Value()
				m.cfg.AI.Custom.BaseURL = m.inputs[1].Value()
				m.cfg.AI.Custom.Model = m.inputs[2].Value()
			}
		}
	}

	// Save the configuration
	if err := config.UpdateConfig(m.cfg); err != nil {
		m.err = fmt.Sprintf("Failed to save configuration: %v", err)
	} else {
		m.successMsg = "Configuration saved successfully!"
		go func() {
			time.Sleep(1 * time.Second)
			m.page = "provider" // Return to provider selection after saving
		}()
	}
}
