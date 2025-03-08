package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	AI struct {
		Provider  string           `mapstructure:"provider"`
		OpenAI    AIProviderConfig `mapstructure:"openai"`
		Anthropic AIProviderConfig `mapstructure:"anthropic"`
		Custom    AIProviderConfig `mapstructure:"custom"`
	} `mapstructure:"ai"`
	UI struct {
		ColorTheme string `mapstructure:"colorTheme"`
		Style      string `mapstructure:"style"`
	} `mapstructure:"ui"`
}

type AIProviderConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

var config Config

func Init(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".roastme" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".roastme")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Set defaults
	viper.SetDefault("ai.provider", "local")
	viper.SetDefault("ai.openai.model", "gpt-3.5-turbo")
	viper.SetDefault("ai.anthropic.model", "claude-2")
	viper.SetDefault("ui.colorTheme", "dark")
	viper.SetDefault("ui.style", "rounded")

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// Create default config file if it doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, _ := homedir.Dir()
			configPath := filepath.Join(home, ".roastme.yaml")

			// Create the default configuration
			defaultConfig := `ai:
  provider: local
  openai:
    api_key: ""
    base_url: "https://api.openai.com/v1"
    model: "gpt-3.5-turbo"
  anthropic:
    api_key: ""
    base_url: "https://api.anthropic.com"
    model: "claude-2"
  custom:
    api_key: ""
    base_url: ""
    model: ""
ui:
  colorTheme: dark
  style: rounded
`
			// Write the default config to file
			if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
				fmt.Println("Warning: Could not create default config file:", err)
			} else {
				fmt.Println("Created default config file:", configPath)
				viper.SetConfigFile(configPath)
				viper.ReadInConfig()
			}
		} else {
			fmt.Println("Warning: Error reading config file:", err)
		}
	}

	// Unmarshal the config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Unable to decode config:", err)
	}
}

func GetConfig() Config {
	return config
}

func SaveConfig() error {
	return viper.WriteConfig()
}

func UpdateConfig(newConfig Config) error {
	config = newConfig

	// Update viper values from the config struct
	viper.Set("ai.provider", config.AI.Provider)
	viper.Set("ai.openai.api_key", config.AI.OpenAI.APIKey)
	viper.Set("ai.openai.base_url", config.AI.OpenAI.BaseURL)
	viper.Set("ai.openai.model", config.AI.OpenAI.Model)
	viper.Set("ai.anthropic.api_key", config.AI.Anthropic.APIKey)
	viper.Set("ai.anthropic.base_url", config.AI.Anthropic.BaseURL)
	viper.Set("ai.anthropic.model", config.AI.Anthropic.Model)
	viper.Set("ai.custom.api_key", config.AI.Custom.APIKey)
	viper.Set("ai.custom.base_url", config.AI.Custom.BaseURL)
	viper.Set("ai.custom.model", config.AI.Custom.Model)
	viper.Set("ui.colorTheme", config.UI.ColorTheme)
	viper.Set("ui.style", config.UI.Style)

	return SaveConfig()
}
