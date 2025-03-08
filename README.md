<div align="center">

# GoRoastMe

### Get roasted by AI for your terminal habits

![img.png](img.png)
[![Go Reference](https://pkg.go.dev/badge/github.com/jasonlovesdoggo/roastme.svg)](https://pkg.go.dev/github.com/jasonlovesdoggo/roastme)
[![Go Report Card](https://goreportcard.com/badge/github.com/jasonlovesdoggo/roastme)](https://goreportcard.com/report/github.com/jasonlovesdoggo/roastme)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

## 📖 Overview

**RoastMe** is a fun CLI tool that analyzes your command history and generates humorous "roasts" based on your terminal habits. 

Whether you're making the same typos, using excessively complex commands, or spending too much time on social media sites through the terminal, RoastMe will call you out with style.

<div align="center">
<img src="https://raw.githubusercontent.com/jasonlovesdoggo/roastme/main/docs/demo.gif" width="700">
</div>

## ✨ Features

- 🔍 Analyzes your shell command history (supports Bash, Zsh, and Fish)
- 🤖 Generates personalized roasts using AI (Google Gemini, OpenAI, or Anthropic)
- 🏠 Works locally without API keys or internet connection
- 🎨 Beautiful TUI using Bubble Tea and Lip Gloss
- 🧩 Extensible with custom AI providers

## 🚀 Installation

### Using Go

```bash
go install github.com/jasonlovesdoggo/roastme@latest
```

### From source

```bash
git clone https://github.com/jasonlovesdoggo/roastme.git
cd roastme
go build -o roastme
sudo mv roastme /usr/local/bin/
```

### Pre-built binaries

Download the [latest release](https://github.com/jasonlovesdoggo/roastme/releases/latest) for your platform.

## 💻 Usage

### Basic Usage

Simply run:

```bash
roastme
```

This will analyze your recent command history and generate a humorous roast using the local (non-AI) engine.

### Advanced Usage

```bash
# Get a deeper analysis of your command history
roastme --deep

# Configure your AI provider settings
roastme config
```

## ⚙️ Configuration

RoastMe supports multiple AI providers, with Google Gemini set as the default. You can configure your preferred provider in two ways:

### 1. Interactive Configuration

Run the configuration UI:

```bash
roastme config
```

This will open an interactive terminal UI where you can:
- Select your preferred AI provider (local, Google Gemini, OpenAI, Anthropic, or custom)
- Enter your API credentials
- Set model preferences

### 2. Manual Configuration

Edit `~/.roastme.toml`:

```toml
[ai]
provider = "gemini" # Options: local, gemini, openai, anthropic, custom

[ai.openai]
api_key = "your-openai-api-key"
base_url = "https://api.openai.com/v1"
model = "gpt-3.5-turbo"

[ai.gemini]
api_key = "your-gemini-api-key"
base_url = "https://generativelanguage.googleapis.com"
model = "gemini-pro"

[ai.anthropic]
api_key = "your-anthropic-api-key"
base_url = "https://api.anthropic.com"
model = "claude-2"

[ai.custom]
api_key = "your-custom-api-key"
base_url = "https://api.your-provider.com"
model = "your-model"

[ui]
colorTheme = "dark" # Options: dark, light
style = "rounded"   # Options: rounded, double, thick
```

## 🔎 What RoastMe Analyzes

RoastMe looks for patterns in your command history, including:

- **Repeated commands** - Are you running the same command over and over?
- **Failed commands** - Typos and error patterns
- **Complex commands** - Extremely long one-liners or pipe chains
- **Indecision** - Excessive use of cd, ls, and other navigation commands
- **Time wasters** - Commands that access time-wasting websites
- **Skill level** - Command complexity to determine your terminal proficiency

## 🤝 Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for the beautiful TUI
- [LangChain Go](https://github.com/tmc/langchaingo) for AI integrations
- [Viper](https://github.com/spf13/viper) for configuration management
- [Cobra](https://github.com/spf13/cobra) for CLI commands

<div align="center">
<p>Made with ❤️ and a sense of humor</p>
</div>
