package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

const (
  OLLAMA_URL   = "http://localhost:11434"
  OLLAMA_MODEL = "llama3"
)

// Config holds our global settings
type Config struct {
	Ollama struct {
		URL   string `yaml:"url"`
		Model string `yaml:"model"`
	} `yaml:"ollama"`
}

// AppConfig is the global instance accessed by other commands
var AppConfig Config

// Inside cmd/config.go
var configInitCmd = &cobra.Command{
	Use:   "config",
	Short: "Initialize a default configuration file in home directory",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		path := filepath.Join(home, ".git-back.yaml")

		// Don't overwrite if it already exists
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("⚠️  Config file already exists at %s\n", path)
			return
		}

		// Create a default config object
		defaultConfig := Config{}

		// defaultConfig.Ollama.URL = "http://localhost:11434"
		// defaultConfig.Ollama.Model = "gemma3:latest"

    defaultConfig.Ollama.URL = AppConfig.Ollama.URL
    defaultConfig.Ollama.Model = AppConfig.Ollama.Model

		// Marshal to YAML
		data, _ := yaml.Marshal(&defaultConfig)
		
		// Write to disk
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			fmt.Printf("❌ Error creating config: %v\n", err)
			return
		}

		fmt.Printf("✨ Created default configuration at %s\n", path)
	},
}


func init() {
	// Centralized command registration
	rootCmd.AddCommand(configInitCmd)
}

func LoadConfig() {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".git-back.yaml")

	// 1. Set "Sensible Defaults" first
	AppConfig.Ollama.URL = OLLAMA_URL 
	AppConfig.Ollama.Model = OLLAMA_MODEL 

	// 2. Try to read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return // If file doesn't exist, we just stay with defaults
	}

	// 3. Overlay file settings onto our defaults
	yaml.Unmarshal(data, &AppConfig)
}
