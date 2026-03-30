package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
  "os"
	"net/http"
	"os"
	"os/exec"
  "path/filepath"
  "strings"


	"github.com/spf13/cobra"
)


var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate AI commit message from staged changes",
	Long:  "Generate AI commit message from staged changes. Perform git commit on the active worktree/branch in preparation to sync files from staging to remote.",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Get staged changes
		diff, err := exec.Command("git", "diff", "--cached").Output()
		if err != nil || len(diff) == 0 {
			fmt.Println("⚠️ No staged changes found. Run 'git add' first.")
			return
		}

		fmt.Println("🤖 Analyzing changes and drafting message...")
		subject, body := askOllamaForCommit(string(diff))

		if subject == "" {
			fmt.Println("❌ Failed to generate a subject.")
			return
		}

		fmt.Println("🤖 Requesting AI commit message...")
		msg, err := callOllama(string(diff))

    if err != nil {
			fmt.Printf("❌ AI Error: %v\n", err)
			fmt.Println("👉 Please start Ollama ('ollama serve') or perform a manual commit:")
			fmt.Println("   git commit -m \"your message\"")
			return // EXIT HERE:
		}
		
		fmt.Printf("\n📝 AI Suggestion: %s\n", msg)

    // Execute the actual git commit
		commitExec := exec.Command("git", "commit", "-m", msg)
		commitExec.Stdout = os.Stdout
		commitExec.Stderr = os.Stderr

    if err := commitExec.Run(); err != nil {
			fmt.Printf("❌ Git commit failed: %v\n", err)
			return
		}

		fmt.Println("✅ Committed successfully!")
	},
}

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
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(configInitCmd)
}

func callOllama(diff string) (string, error) {
	prompt := "Write a concise Conventional Commit message for this diff. No conversational filler:\n\n" + diff
	
	payload, _ := json.Marshal(map[string]interface{}{
		"model": AppConfig.Ollama.Model, "prompt": prompt, "stream": false,
	})

	resp, err := http.Post(AppConfig.Ollama.Model, "application/json", bytes.NewBuffer(payload))

	if err != nil {
    // Check if the error is specifically a connection failure
    fmt.Printf("⚠️  Ollama appears to be offline at %s\n", AppConfig.Ollama.Model)
    fmt.Println("👉 Check if Ollama is running: 'ollama serve'")
		return "❌ fix: Manual update required (Ollama Offline)", err
	}
	
	defer resp.Body.Close()

	var res struct{ Response string `json:"response"` }
  if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		// return "", fmt.Errorf("failed to decode AI response")
    return "", err

	}

  return strings.TrimSpace(res.Response), nil
}
