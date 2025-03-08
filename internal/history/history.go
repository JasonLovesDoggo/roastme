package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
)

// CommandEntry represents a command with its timestamp
type CommandEntry struct {
	Command   string
	Timestamp time.Time
}

// GetShellHistory returns the shell command history
func GetShellHistory(limit int) ([]string, error) {
	// Detect shell
	shell := os.Getenv("SHELL")

	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var historyFile string
	var parseHistoryFile func(string, int) ([]string, error)

	// Determine history file and parse function based on shell
	if strings.Contains(shell, "zsh") {
		historyFile = filepath.Join(home, ".zsh_history")
		parseHistoryFile = parseZshHistory
	} else if strings.Contains(shell, "bash") {
		historyFile = filepath.Join(home, ".bash_history")
		parseHistoryFile = parseBashHistory
	} else if strings.Contains(shell, "fish") {
		historyFile = filepath.Join(home, ".local", "share", "fish", "fish_history")
		parseHistoryFile = parseFishHistory
	} else {
		// Fallback to bash history format
		historyFile = filepath.Join(home, ".bash_history")
		parseHistoryFile = parseBashHistory
	}

	// Check if history file exists
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return []string{"No history file found"}, nil
	}

	// Parse history file with the appropriate function
	commands, err := parseHistoryFile(historyFile, limit)
	if err != nil {
		return nil, fmt.Errorf("error parsing history file: %v", err)
	}

	return commands, nil
}

// parseBashHistory efficiently parses bash history file
func parseBashHistory(historyFile string, limit int) ([]string, error) {
	file, err := os.Open(historyFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a circular buffer-like structure of fixed capacity
	commands := make([]string, 0, limit)

	scanner := bufio.NewScanner(file)
	lineCount := 0

	// For very large history files, we might need to set a custom buffer
	// This allows scanning lines longer than the default buffer size
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, 64*1024)    // 64KB initial buffer
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		// Skip timestamp lines if present (they start with #)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the command line
		cmd := strings.TrimSpace(line)
		if cmd != "" {
			commands = append(commands, cmd)

			// If we exceed the limit, remove oldest commands
			if limit > 0 && len(commands) > limit {
				// Shift elements by removing the oldest
				commands = commands[1:]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return commands, fmt.Errorf("error reading history file: %v", err)
	}

	return commands, nil
}

// parseZshHistory efficiently parses zsh history file
func parseZshHistory(historyFile string, limit int) ([]string, error) {
	file, err := os.Open(historyFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a circular buffer-like structure of fixed capacity
	commands := make([]string, 0, limit)

	// ZSH history format regexp: ": TIMESTAMP:0;COMMAND"
	// or simply "COMMAND" without timestamp
	re := regexp.MustCompile(`: (\d+):\d+;(.*)`)

	scanner := bufio.NewScanner(file)

	// For very large history files, we might need to set a custom buffer
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, 64*1024)    // 64KB initial buffer
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()

		// Check for the timestamp format
		matches := re.FindStringSubmatch(line)
		if len(matches) > 2 {
			// Extract command from timestamp format
			cmd := strings.TrimSpace(matches[2])
			if cmd != "" {
				commands = append(commands, cmd)
			}
		} else {
			// Try direct approach if regex fails
			parts := strings.SplitN(line, ";", 2)
			if len(parts) > 1 {
				cmd := strings.TrimSpace(parts[1])
				if cmd != "" {
					commands = append(commands, cmd)
				}
			} else if strings.TrimSpace(line) != "" {
				// If it's not a timestamp format but not empty either
				commands = append(commands, strings.TrimSpace(line))
			}
		}

		// If we exceed the limit, remove oldest commands
		if limit > 0 && len(commands) > limit {
			// Shift elements by removing the oldest
			commands = commands[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return commands, fmt.Errorf("error reading history file: %v", err)
	}

	return commands, nil
}

// parseFishHistory parses fish shell history which is stored in a different format
func parseFishHistory(filePath string, limit int) ([]string, error) {
	// Fish history is stored in a more complex format
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []string{"No fish history file found"}, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a circular buffer-like structure of fixed capacity
	commands := make([]string, 0, limit)

	// Fish history entry looks like:
	// - cmd: the actual command
	// - when: timestamp
	type fishEntry struct {
		Cmd  string `json:"cmd"`
		When int64  `json:"when"`
	}

	// Fish history can be quite large, so we need a robust scanner
	scanner := bufio.NewScanner(file)

	// For very large history files, we might need to set a custom buffer
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, 64*1024)    // 64KB initial buffer
	scanner.Buffer(buf, maxCapacity)

	// Regular expression for extracting command from non-JSON format
	cmdRegex := regexp.MustCompile(`- cmd: (.+)`)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Try to parse as JSON
		var entry fishEntry
		if err := json.Unmarshal([]byte(line), &entry); err == nil && entry.Cmd != "" {
			commands = append(commands, entry.Cmd)
		} else {
			// Fall back to regex
			matches := cmdRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				cmd := strings.TrimSpace(matches[1])
				if cmd != "" {
					commands = append(commands, cmd)
				}
			}
		}

		// If we exceed the limit, remove oldest commands
		if limit > 0 && len(commands) > limit {
			// Shift elements by removing the oldest
			commands = commands[1:]
		}
	}

	if err := scanner.Err(); err != nil {
		return commands, fmt.Errorf("error reading fish history file: %v", err)
	}

	return commands, nil
}

// GetMostRecentCommands returns only the most recent commands
// with a maximum count specified by the limit parameter
func GetMostRecentCommands(allCommands []string, limit int) []string {
	if len(allCommands) <= limit || limit <= 0 {
		return allCommands
	}

	// Return only the most recent commands
	return allCommands[len(allCommands)-limit:]
}

// GetCommandsWithTimestamps attempts to get commands with timestamps from history
// This is a more advanced function that tries to extract timing information
// Only works for shells that store timestamp information (zsh, fish)
func GetCommandsWithTimestamps(limit int) ([]CommandEntry, error) {
	// This is a placeholder implementation - would need specific shell handling
	// Similar to GetShellHistory but with timestamp extraction
	shell := os.Getenv("SHELL")

	// For now just return empty
	return []CommandEntry{}, fmt.Errorf("timestamp extraction not implemented for %s", shell)
}

// FilterCommands returns commands matching the given pattern
func FilterCommands(commands []string, pattern string) []string {
	var filtered []string

	re, err := regexp.Compile(pattern)
	if err != nil {
		// Fall back to simple string matching if regex is invalid
		for _, cmd := range commands {
			if strings.Contains(cmd, pattern) {
				filtered = append(filtered, cmd)
			}
		}
	} else {
		for _, cmd := range commands {
			if re.MatchString(cmd) {
				filtered = append(filtered, cmd)
			}
		}
	}

	return filtered
}
