package cmd

import (
  "fmt"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Ollama struct {
		URL   string `yaml:"url"`
		Model string `yaml:"model"`
	} `yaml:"ollama"`
	GitHub struct {
		Token string `yaml:"token"`
	} `yaml:"github"`
}

var AppConfig Config

// configCmd represents the base config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage git-back configuration",
}

// initConfigCmd creates a default YAML file
var initConfigCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		path := filepath.Join(home, ".git-back.yaml")

		if _, err := os.Stat(path); err == nil {
			fmt.Printf("⚠️  Config already exists at %s. Skipping.\n", path)
			return
		}

		// Default Settings
		AppConfig.Ollama.URL = "http://localhost:11434"
		AppConfig.Ollama.Model = "gemma3:latest"

		data, _ := yaml.Marshal(&AppConfig)
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			fmt.Printf("❌ Failed to write config: %v\n", err)
			return
		}

		fmt.Printf("✨ Created default config at %s\n", path)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initConfigCmd)
}

func LoadConfig() error {
	// Default path: ~/.git-back.yaml
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".git-back.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, set sensible defaults
		AppConfig.Ollama.URL = "http://localhost:11434"
		AppConfig.Ollama.Model = "gemma3:latest"
		return nil
	}

	return yaml.Unmarshal(data, &AppConfig)
}
