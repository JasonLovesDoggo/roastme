package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jasonlovesdoggo/roastme/internal/ai"
	"github.com/jasonlovesdoggo/roastme/internal/analysis"
	"github.com/jasonlovesdoggo/roastme/internal/config"
	"github.com/jasonlovesdoggo/roastme/internal/history"
	"github.com/jasonlovesdoggo/roastme/internal/ui"
	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	deep         bool
	complexity   string
	commandLimit int
)

var rootCmd = &cobra.Command{
	Use:   "roastme",
	Short: "Endless roasts of your command line history using AI",
	Long: `GoRoastMe is a CLI tool that analyzes your command history
and generates endless, hilarious roasts about your terminal habits.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure the application
		cfg := config.GetConfig()

		// Run the interactive roasting mode
		runInteractiveMode(cfg)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure GoRoastMe settings",
	Run: func(cmd *cobra.Command, args []string) {
		ui.RunConfigInterface()
	},
}

func runInteractiveMode(cfg config.Config) {
	fmt.Println("GoRoastMe - Terminal History Roaster")
	fmt.Println("Press Enter for a new roast, Ctrl+C to exit")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	count := 0

	for {
		// Generate initial roast or wait for Enter for subsequent roasts
		if count > 0 {
			fmt.Print("\nPress Enter for another roast... ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			// Exit if user types "exit" or "quit"
			if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
				fmt.Println("Exiting GoRoastMe. Your terminal is safe... for now.")
				return
			}
		}

		// Create a spinner
		spinner := ui.NewSpinner("Analyzing your command history...")
		spinner.Start()

		// Calculate actual command limit
		actualLimit := commandLimit
		if deep {
			actualLimit *= 5 // Multiply by 5 in deep mode
		}

		// Get command history - using actualLimit
		commands, err := history.GetShellHistory(actualLimit)
		if err != nil {
			spinner.Stop()
			fmt.Fprintf(os.Stderr, "Error getting shell history: %v\n", err)
			continue
		}

		// Analyze command patterns
		patterns := analysis.AnalyzeHistory(commands)

		// Generate roast with selected complexity
		level := getComplexityLevel()
		roast, err := ai.GenerateRoast(cfg, patterns, commands, level)
		if err != nil {
			spinner.Stop()
			fmt.Fprintf(os.Stderr, "Error generating roast: %v\n", err)
			continue
		}

		spinner.Stop()

		// Display the roast with current count
		ui.DisplayInteractiveRoast(roast, count+1, len(commands))

		count++
	}
}

func getComplexityLevel() ai.ComplexityLevel {
	switch complexity {
	case "simple":
		return ai.SimpleRoast
	case "normal":
		return ai.NormalRoast
	case "complex":
		return ai.ComplexRoast
	case "brutal":
		return ai.BrutalRoast
	default:
		return ai.NormalRoast
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goroastme.yaml)")
	rootCmd.Flags().BoolVar(&deep, "deep", false, "Analyze 5x more commands than the default limit")
	rootCmd.Flags().StringVar(&complexity, "complexity", "normal", "Roast complexity: simple, normal, complex, or brutal")
	rootCmd.Flags().IntVar(&commandLimit, "limit", 500, "Number of commands to analyze (multiplied by 5 in deep mode)")

	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory
	config.Init(cfgFile)
}
