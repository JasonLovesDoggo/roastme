package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/jasonlovesdoggo/roastme/internal/ai"
	"github.com/jasonlovesdoggo/roastme/internal/analysis"
	"github.com/jasonlovesdoggo/roastme/internal/config"
	"github.com/jasonlovesdoggo/roastme/internal/history"
	"github.com/jasonlovesdoggo/roastme/internal/ui"
	"github.com/spf13/cobra"
)

const (
	LowerLimit = 100
	UpperLimit = 500
)

var (
	cfgFile      string
	deep         bool
	continuous   bool
	interval     int
	maxRoasts    int
	complexity   string
	notifySystem bool
)

var rootCmd = &cobra.Command{
	Use:   "roastme",
	Short: "Roast your command line history using AI",
	Long: `roastme is a CLI tool that analyzes your command history
and generates humorous roasts based on your terminal habits.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure the application
		cfg := config.GetConfig()

		// If continuous mode is enabled
		if continuous {
			runContinuousMode(cfg)
			return
		}

		// Create a spinner
		spinner := ui.NewSpinner("Analyzing your command history...")
		spinner.Start()

		// Get command history
		limit := LowerLimit
		if deep {
			limit = UpperLimit
		}
		commands, err := history.GetShellHistory(limit)
		if err != nil {
			spinner.Stop()
			fmt.Fprintf(os.Stderr, "Error getting shell history: %v\n", err)
			os.Exit(1)
		}

		// Analyze command patterns
		patterns := analysis.AnalyzeHistory(commands)

		// Generate roast with selected complexity
		roast, err := ai.GenerateRoast(cfg, patterns, commands, getComplexityLevel())
		if err != nil {
			spinner.Stop()
			fmt.Fprintf(os.Stderr, "Error generating roast: %v\n", err)
			os.Exit(1)
		}

		spinner.Stop()

		// Display the roast
		ui.DisplayRoast(roast)
	},
}

// continuousCmd represents the continuous command
var continuousCmd = &cobra.Command{
	Use:   "continuous",
	Short: "Run GoRoastMe in continuous mode",
	Long: `Continuously monitors your command history and generates
new roasts at specified intervals.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		runContinuousMode(cfg)
	},
}

func runContinuousMode(cfg config.Config) {
	fmt.Println("Starting continuous roast mode...")
	fmt.Printf("Generating a new roast every %d minutes\n", interval)
	fmt.Println("Press Ctrl+C to exit")

	count := 0
	for {
		if maxRoasts > 0 && count >= maxRoasts {
			fmt.Println("Reached maximum number of roasts. Exiting.")
			return
		}

		// Clear screen between roasts
		if count > 0 {
			time.Sleep(time.Duration(interval) * time.Minute)
		}

		// Get updated command history
		limit := LowerLimit
		if deep {
			limit = UpperLimit
		}
		commands, err := history.GetShellHistory(limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting shell history: %v\n", err)
			continue
		}

		// Analyze command patterns
		patterns := analysis.AnalyzeHistory(commands)

		// Generate roast with selected complexity
		roast, err := ai.GenerateRoast(cfg, patterns, commands, getComplexityLevel())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating roast: %v\n", err)
			continue
		}

		// Send system notification if enabled
		if notifySystem {
			ui.SendNotification("GoRoastMe", "New roast generated!")
		}

		// Display the roast with timestamp
		ui.DisplayContinuousRoast(roast, count+1)

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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure goroastme settings",
	Run: func(cmd *cobra.Command, args []string) {
		ui.RunConfigInterface()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goroastme.yaml)")
	rootCmd.Flags().BoolVar(&deep, "deep", false, "Perform a deeper, more personal roast")
	rootCmd.Flags().BoolVar(&continuous, "continuous", false, "Run in continuous mode, generating roasts at intervals")
	rootCmd.Flags().IntVar(&interval, "interval", 5, "Interval in minutes between roasts in continuous mode")
	rootCmd.Flags().IntVar(&maxRoasts, "max", 0, "Maximum number of roasts to generate (0 = unlimited)")
	rootCmd.Flags().StringVar(&complexity, "complexity", "normal", "Roast complexity: simple, normal, complex, or brutal")
	rootCmd.Flags().BoolVar(&notifySystem, "notify", false, "Send system notifications when new roasts are generated")

	// Add the continuous command separately
	continuousCmd.Flags().IntVar(&interval, "interval", 5, "Interval in minutes between roasts")
	continuousCmd.Flags().IntVar(&maxRoasts, "max", 0, "Maximum number of roasts to generate (0 = unlimited)")
	continuousCmd.Flags().StringVar(&complexity, "complexity", "normal", "Roast complexity: simple, normal, complex, or brutal")
	continuousCmd.Flags().BoolVar(&notifySystem, "notify", false, "Send system notifications when new roasts are generated")
	continuousCmd.Flags().BoolVar(&deep, "deep", false, "Perform a deeper, more personal roast")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(continuousCmd)
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory
	config.Init(cfgFile)
}
