package analysis

import (
	"strings"
)

// CommandPattern represents patterns found in command history
type CommandPattern struct {
	RepeatedCommands []CommandCount
	FailedCommands   []string
	ComplexCommands  []string
	Indecisive       bool
	TimeWasters      []string
	SkillLevel       string
}

// CommandCount represents a command and its frequency
type CommandCount struct {
	Command string
	Count   int
}

// AnalyzeHistory analyzes command patterns in history
func AnalyzeHistory(commands []string) CommandPattern {
	patterns := CommandPattern{
		RepeatedCommands: []CommandCount{},
		FailedCommands:   []string{},
		ComplexCommands:  []string{},
		Indecisive:       false,
		TimeWasters:      []string{},
		SkillLevel:       "beginner", // Default
	}

	// Count command frequencies
	commandCounts := make(map[string]int)
	for _, cmd := range commands {
		// Extract base command
		baseCmdParts := strings.Fields(cmd)
		if len(baseCmdParts) > 0 {
			baseCmd := baseCmdParts[0]
			commandCounts[baseCmd]++
		}
	}

	// Find repeated commands
	for cmd, count := range commandCounts {
		if count > 3 && cmd != "" {
			patterns.RepeatedCommands = append(patterns.RepeatedCommands, CommandCount{
				Command: cmd,
				Count:   count,
			})
		}
	}

	// Look for failed commands and corrections
	for i := 1; i < len(commands); i++ {
		prevCmd := commands[i-1]
		currCmd := commands[i]

		// Check for potential corrections
		if (strings.HasPrefix(prevCmd, "cd ") ||
			strings.HasPrefix(prevCmd, "mkdir ") ||
			strings.HasPrefix(prevCmd, "git ")) &&
			strings.HasPrefix(currCmd, prevCmd[:min(len(prevCmd), 4)]) {
			patterns.FailedCommands = append(patterns.FailedCommands, prevCmd)
		}
	}

	// Check for complex commands
	for _, cmd := range commands {
		if strings.Count(cmd, "|") > 2 || strings.Count(cmd, ";") > 2 || len(cmd) > 80 {
			patterns.ComplexCommands = append(patterns.ComplexCommands, cmd)
		}
	}

	// Check for indecisiveness (lots of cd, ls in sequence)
	cdLsCount := 0
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, "cd ") || strings.HasPrefix(cmd, "ls ") || cmd == "ls" {
			cdLsCount++
		}
	}

	if float64(cdLsCount) > float64(len(commands))*0.4 { // If more than 40% are cd/ls
		patterns.Indecisive = true
	}

	// Look for "time waster" commands
	timeWasters := []string{"reddit", "youtube", "twitter", "facebook", "instagram"}
	for _, cmd := range commands {
		cmdLower := strings.ToLower(cmd)
		for _, waster := range timeWasters {
			if strings.Contains(cmdLower, waster) && !contains(patterns.TimeWasters, waster) {
				patterns.TimeWasters = append(patterns.TimeWasters, waster)
			}
		}
	}

	// Guess skill level based on command complexity
	advancedCmds := []string{"awk", "sed", "grep -E", "xargs", "find -exec", "docker", "kubernetes", "k8s", "kubectl"}
	advancedCount := 0

	for _, cmd := range commands {
		cmdLower := strings.ToLower(cmd)
		for _, advCmd := range advancedCmds {
			if strings.Contains(cmdLower, advCmd) {
				advancedCount++
				break
			}
		}
	}

	if advancedCount > 5 {
		patterns.SkillLevel = "advanced"
	} else if advancedCount > 2 {
		patterns.SkillLevel = "intermediate"
	}

	return patterns
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
