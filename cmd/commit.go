package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate an AI commit message from staged changes",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Get staged changes
		diff, err := exec.Command("git", "diff", "--cached").Output()
		if err != nil || len(diff) == 0 {
			fmt.Println("⚠️ No staged changes found. Run 'git add' first.")
			return
		}

		fmt.Println("🤖 Analyzing changes and drafting message...")
		msg := askOllamaForCommit(string(diff))

		if msg == "" {
			fmt.Println("❌ Failed to generate a message.")
			return
		}

		// 2. Show the user the message
		fmt.Printf("\nProposed Commit Message:\n---\n%s\n---\n", msg)

		// 3. Execute the commit
		commitExec := exec.Command("git", "commit", "-m", msg)
		commitExec.Stdout = os.Stdout
		commitExec.Stderr = os.Stderr
		
		if err := commitExec.Run(); err != nil {
			fmt.Printf("❌ Git commit failed: %v\n", err)
			return
		}
		
		fmt.Println("✨ Commit successful!")
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func askOllamaForCommit(diff string) string {
	url := fmt.Sprintf("%s/api/generate", AppConfig.Ollama.URL)

	prompt := `Write a concise, professional Git commit message for these changes.
Use the Conventional Commits format (e.g., "feat: ...", "fix: ...").
Keep the first line under 50 characters.

DIFF:
` + diff

	payload, _ := json.Marshal(map[string]interface{}{
		"model":  AppConfig.Ollama.Model,
		"prompt": prompt,
		"stream": false,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var res struct{ Response string }
	json.NewDecoder(resp.Body).Decode(&res)
	
	return strings.TrimSpace(res.Response)
}
