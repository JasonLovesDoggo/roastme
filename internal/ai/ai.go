package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/jasonlovesdoggo/roastme/internal/analysis"
	"github.com/jasonlovesdoggo/roastme/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
)

// ComplexityLevel defines how elaborate the roast should be
type ComplexityLevel int

const (
	SimpleRoast ComplexityLevel = iota
	NormalRoast
	ComplexRoast
	BrutalRoast
)

// GenerateRoast generates a roast based on the command patterns
func GenerateRoast(cfg config.Config, patterns analysis.CommandPattern, commands []string, complexity ComplexityLevel) (string, error) {
	// Use local roasts if no AI provider is configured or provider is set to "local"
	if cfg.AI.Provider == "" || cfg.AI.Provider == "local" {
		return generateLocalRoast(patterns, complexity), nil
	}

	// Try to generate a roast using the configured AI provider
	roast, err := generateAIRoast(cfg, patterns, commands, complexity)
	if err != nil {
		// Fall back to local roasts if AI fails
		return generateLocalRoast(patterns, complexity), nil
	}

	return roast, nil
}

// generateAIRoast generates a roast using the configured AI provider
func generateAIRoast(cfg config.Config, patterns analysis.CommandPattern, commands []string, complexity ComplexityLevel) (string, error) {
	// Prepare context
	ctx := context.Background()

	// Get appropriate number of commands based on complexity
	recentCmds := commands
	//cmdLimit := 10
	//
	//switch complexity {
	//case SimpleRoast:
	//	cmdLimit = 5
	//case NormalRoast:
	//	cmdLimit = 10
	//case ComplexRoast:
	//	cmdLimit = 20
	//case BrutalRoast:
	//	cmdLimit = 30
	//}
	//
	//if len(recentCmds) > cmdLimit {
	//	recentCmds = recentCmds[len(recentCmds)-cmdLimit:]
	//}

	// Create prompt based on command patterns and complexity level
	prompt := createPromptForComplexity(recentCmds, patterns, complexity)

	var llm llms.Model
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

	// Generate the roast with appropriate parameters based on complexity
	maxTokens := 150
	temperature := 0.7

	switch complexity {
	case SimpleRoast:
		maxTokens = 100
		temperature = 0.5
	case NormalRoast:
		maxTokens = 150
		temperature = 0.7
	case ComplexRoast:
		maxTokens = 300
		temperature = 0.8
	case BrutalRoast:
		maxTokens = 500
		temperature = 0.9
	}

	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(temperature),
		llms.WithMaxTokens(maxTokens),
	)

	if err != nil {
		return "", err
	}

	return completion, nil
}

// createPromptForComplexity creates a prompt based on the desired complexity level
func createPromptForComplexity(commands []string, patterns analysis.CommandPattern, complexity ComplexityLevel) string {
	systemPrompt := "You are an Arch Linux user who lives in your terminal. Your only purpose is to roast people about their command-line habits.\nBe clever, humorous, and unapologetically savage. Analyze shell history with deep technical insight, and deliver biting, hilarious roasts.\nThink of yourself as the Gordon Ramsay of the terminal—brutal but constructive. Your tone is sharp, witty, and tech-savvy, always packed\nwith programming humor."

	basePrompt := fmt.Sprintf(`
Recent commands:
%s

Patterns found:
- Repeated commands: %v
- Failed commands: %v
- Complex commands: %v
- Indecisive: %v
- Time wasters: %v
- Skill level: %s
`, formatCommands(commands), patterns.RepeatedCommands,
		patterns.FailedCommands, patterns.ComplexCommands,
		patterns.Indecisive, patterns.TimeWasters, patterns.SkillLevel)
	switch complexity {
	case SimpleRoast:
		return systemPrompt + "Roast this person based on their command line history. Be concise and mildly amusing." +
			basePrompt + "\nGenerate a short, simple roast (1-2 sentences) about their terminal habits."

	case NormalRoast:
		return systemPrompt + "Roast this person based on their command line history. Be funny but not mean." +
			basePrompt + "\nGenerate a moderate-length roast (2-3 sentences) about their terminal habits."

	case ComplexRoast:
		return systemPrompt + "Roast this person based on their command line history. Be clever, " +
			"insightful and humorous." +
			basePrompt +
			"\nGenerate a detailed roast (3-4 paragraphs) about their terminal habits. Include specific observations about their " +
			"command patterns, technical skill level, and potential personality traits that might be revealed by their commands. " +
			"Be creative and witty, using tech humor and programming references."

	case BrutalRoast:
		return systemPrompt + "Roast this person based on their command line history. Be extremely thorough, " +
			"devastatingly funny, " +
			"and borderline ruthless." +
			basePrompt +
			"\nWrite a comprehensive, brutal roast (4+ paragraphs) that thoroughly analyzes their terminal habits. " +
			"Include specific references to their commands, create an entire psychological profile based on their terminal behavior, " +
			"make wild assumptions about their coding abilities, and don't hold back on the technical humor. " +
			"Imagine this is a Comedy Central Roast but for developers. Be creative, savage but still ultimately good-natured."
	}

	return systemPrompt + "Roast this person based on their command line history." + basePrompt
}

// initOpenAI initializes the OpenAI client
func initOpenAI(cfg config.Config) (llms.Model, error) {
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
func initAnthropic(cfg config.Config) (llms.Model, error) {
	return nil, errors.New("anthropic support not implemented yet")
}

// initGemini initializes the Google Gemini client
func initGemini(cfg config.Config) (llms.Model, error) {
	if cfg.AI.Gemini.APIKey == "" {
		return nil, errors.New("google Gemini API key not configured")
	}

	ctx := context.Background()
	options := []googleai.Option{
		googleai.WithAPIKey(cfg.AI.Gemini.APIKey),
	}

	return googleai.New(ctx, options...)
}

// initCustom initializes a custom LLM client
func initCustom(cfg config.Config) (llms.Model, error) {
	return nil, errors.New("custom LLM provider support not implemented yet")
}

// formatCommands formats a slice of commands for inclusion in the prompt
func formatCommands(commands []string) string {
	result := ""
	for _, cmd := range commands {
		result += "- " + cmd + "\n"
	}
	return result
}
