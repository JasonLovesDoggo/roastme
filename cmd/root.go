package cmd

import (
	"fmt"
	"os"

	"github.com/jasonlovesdoggo/roastme/internal/ai"
	"github.com/jasonlovesdoggo/roastme/internal/analysis"
	"github.com/jasonlovesdoggo/roastme/internal/config"
	"github.com/jasonlovesdoggo/roastme/internal/history"
	"github.com/jasonlovesdoggo/roastme/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	deep    bool
)

var rootCmd = &cobra.Command{
	Use:   "roastme",
	Short: "Roast your command line history using AI",
	Long: `roastme is a CLI tool that analyzes your command history
and generates humorous roasts based on your terminal habits.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure the application
		cfg := config.GetConfig()

		// Create a spinner
		spinner := ui.NewSpinner("Analyzing your command history...")
		spinner.Start()

		// Get command history
		limit := 30
		if deep {
			limit = 100
		}
		commands, err := history.GetShellHistory(limit)
		if err != nil {
			spinner.Stop()
			fmt.Fprintf(os.Stderr, "Error getting shell history: %v\n", err)
			os.Exit(1)
		}

		// Analyze command patterns
		patterns := analysis.AnalyzeHistory(commands)

		// Generate roast
		roast, err := ai.GenerateRoast(cfg, patterns, commands)
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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure roastme settings",
	Run: func(cmd *cobra.Command, args []string) {
		ui.RunConfigInterface()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.roastme.yaml)")
	rootCmd.Flags().BoolVar(&deep, "deep", false, "Perform a deeper, more personal roast")

	rootCmd.AddCommand(configCmd)
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory
	config.Init(cfgFile)
}
