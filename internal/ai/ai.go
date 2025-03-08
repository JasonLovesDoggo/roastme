package ai

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jasonlovesdoggo/roastme/internal/analysis"
	"github.com/jasonlovesdoggo/roastme/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
)

// GenerateRoast generates a roast based on the command patterns
func GenerateRoast(cfg config.Config, patterns analysis.CommandPattern, commands []string) (string, error) {
	// Use local roasts if no AI provider is configured or provider is set to "local"
	if cfg.AI.Provider == "" || cfg.AI.Provider == "local" {
		return generateLocalRoast(patterns), nil
	}
	// Try to generate a roast using the configured AI provider
	roast, err := generateAIRoast(cfg, patterns, commands)
	if err != nil {
		// Fall back to local roasts if AI fails
		return generateLocalRoast(patterns), nil
	}

	return roast, nil
}

// generateAIRoast generates a roast using the configured AI provider
func generateAIRoast(cfg config.Config, patterns analysis.CommandPattern, commands []string) (string, error) {
	// Prepare context
	ctx := context.Background()

	// Get the last 10 commands for the prompt
	recentCmds := commands
	if len(recentCmds) > 10 {
		recentCmds = recentCmds[len(recentCmds)-10:]
	}

	// Create prompt based on command patterns
	prompt := fmt.Sprintf(`
Roast this person based on their command line history. Be funny but not mean.

Recent commands:
%s

Patterns found:
- Repeated commands: %v
- Failed commands: %v
- Complex commands: %v
- Indecisive: %v
- Time wasters: %v
- Skill level: %s

Generate a short, funny roast (1-3 sentences) about their terminal habits.
`, formatCommands(recentCmds), patterns.RepeatedCommands,
		patterns.FailedCommands, patterns.ComplexCommands,
		patterns.Indecisive, patterns.TimeWasters, patterns.SkillLevel)

	var llm llms.LLM
	var err error

	// Initialize the appropriate LLM based on provider
	switch cfg.AI.Provider {
	case "openai":
		llm, err = initOpenAI(cfg)
	case "anthropic":
		llm, err = initAnthropic(cfg)
	case "gemini":
		llm, err = initGemini(cfg)
	case "custom":
		llm, err = initCustom(cfg)
	default:
		return "", errors.New("unsupported AI provider")
	}

	if err != nil {
		return "", err
	}

	// Generate the roast
	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(0.7),
		llms.WithMaxTokens(150),
	)

	if err != nil {
		return "", err
	}

	return completion, nil
}

// initOpenAI initializes the OpenAI client
func initOpenAI(cfg config.Config) (llms.LLM, error) {
	if cfg.AI.OpenAI.APIKey == "" {
		return nil, errors.New("OpenAI API key not configured")
	}

	options := []openai.Option{
		openai.WithToken(cfg.AI.OpenAI.APIKey),
	}

	if cfg.AI.OpenAI.BaseURL != "" {
		options = append(options, openai.WithBaseURL(cfg.AI.OpenAI.BaseURL))
	}

	if cfg.AI.OpenAI.Model != "" {
		options = append(options, openai.WithModel(cfg.AI.OpenAI.Model))
	}

	return openai.New(options...)
}

// initAnthropic initializes the Anthropic client
func initAnthropic(cfg config.Config) (llms.LLM, error) {
	// Note: As of this implementation, langchaingo might not have direct Anthropic support
	// This is a placeholder for when it's available or you could implement a custom connector
	return nil, errors.New("Anthropic support not implemented yet")
}

// initGemini initializes the Google Gemini client
func initGemini(cfg config.Config) (llms.LLM, error) {
	if cfg.AI.Gemini.APIKey == "" {
		return nil, errors.New("Google Gemini API key not configured")
	}

	ctx := context.Background()
	options := []googleai.Option{
		googleai.WithAPIKey(cfg.AI.Gemini.APIKey),
	}

	// Note: Gemini's API is different, we'll set the model using the option directly
	if cfg.AI.Gemini.Model != "" {
		// The model needs to be set separately during Call()
	}

	// googleai.New() requires a context as first parameter
	return googleai.New(ctx, options...)
}

// initCustom initializes a custom LLM client
func initCustom(cfg config.Config) (llms.LLM, error) {
	// This would implement a custom LLM client
	return nil, errors.New("Custom LLM provider support not implemented yet")
}

// generateLocalRoast generates a roast without using an external AI service
func generateLocalRoast(patterns analysis.CommandPattern) string {
	roasts := []string{}

	if len(patterns.RepeatedCommands) > 0 {
		cmd := patterns.RepeatedCommands[0].Command
		count := patterns.RepeatedCommands[0].Count
		roasts = append(roasts, fmt.Sprintf("I see you've used '%s' %d times. Having memory issues or just really, really in love with that command?", cmd, count))
	}

	if len(patterns.FailedCommands) > 0 {
		roasts = append(roasts, fmt.Sprintf("Nice typos! Maybe typing lessons should be in your future before attempting '%s'.", patterns.FailedCommands[0]))
	}

	if len(patterns.ComplexCommands) > 0 {
		roasts = append(roasts, "Wow, those complex commands! Trying to impress an invisible audience or just afraid of using separate lines?")
	}

	if patterns.Indecisive {
		roasts = append(roasts, "All those cd's and ls's... are you exploring your filesystem or just completely lost in there?")
	}

	if len(patterns.TimeWasters) > 0 {
		waster := patterns.TimeWasters[0]
		roasts = append(roasts, fmt.Sprintf("I see you're visiting %s. Working hard or hardly working, eh?", waster))
	}

	if len(roasts) == 0 {
		skill := patterns.SkillLevel
		if skill == "beginner" {
			roasts = append(roasts, "Your command history screams 'I just discovered what a terminal is'. How adorable.")
		} else if skill == "intermediate" {
			roasts = append(roasts, "Your command history is like a mediocre pizza - it gets the job done but nobody's impressed.")
		} else {
			roasts = append(roasts, "Fancy commands! Overcompensating for something or just showing off to nobody?")
		}
	}

	// Make sure we have at least one roast
	if len(roasts) == 0 {
		roasts = append(roasts, "I can't even roast your command history - it's that boring. Try doing something interesting first!")
	}

	// Return a random roast
	rand.Seed(time.Now().UnixNano())
	return roasts[rand.Intn(len(roasts))]
}

// formatCommands formats a slice of commands for inclusion in the prompt
func formatCommands(commands []string) string {
	result := ""
	for _, cmd := range commands {
		result += "- " + cmd + "\n"
	}
	return result
}
