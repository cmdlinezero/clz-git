package cmd

import (
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
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
