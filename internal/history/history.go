package history

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// GetShellHistory returns the shell command history
func GetShellHistory(limit int) ([]string, error) {
	// Detect shell
	shell := os.Getenv("SHELL")

	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	var historyFile string
	var parseFunc func(string) (string, bool)

	// Determine history file and parse function based on shell
	if strings.Contains(shell, "zsh") {
		historyFile = filepath.Join(home, ".zsh_history")
		parseFunc = parseZshHistoryLine
	} else if strings.Contains(shell, "bash") {
		historyFile = filepath.Join(home, ".bash_history")
		parseFunc = parseBashHistoryLine
	} else {
		// Fallback to bash history format
		historyFile = filepath.Join(home, ".bash_history")
		parseFunc = parseBashHistoryLine
	}

	// Check if history file exists
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return []string{"No history file found"}, nil
	}

	// Read history file
	file, err := os.Open(historyFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if cmd, ok := parseFunc(line); ok && cmd != "" {
			commands = append(commands, cmd)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Get the most recent commands
	if len(commands) > limit {
		commands = commands[len(commands)-limit:]
	}

	return commands, nil
}

// parseBashHistoryLine extracts command from a bash history line
func parseBashHistoryLine(line string) (string, bool) {
	// Bash history is straightforward, just trim spaces
	return strings.TrimSpace(line), true
}

// parseZshHistoryLine extracts command from a zsh history line
// zsh history format is more complex with timestamps
func parseZshHistoryLine(line string) (string, bool) {
	// Example zsh history line:
	// : 1617985604:0;ls -la

	re := regexp.MustCompile(`: \d+:\d+;(.*)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) > 1 {
		return strings.TrimSpace(matches[1]), true
	}

	// Try direct approach if regex fails
	parts := strings.SplitN(line, ";", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1]), true
	}

	return strings.TrimSpace(line), true
}
