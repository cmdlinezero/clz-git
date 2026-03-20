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
		subject, body := askOllamaForCommit(string(diff))

		if subject == "" {
			fmt.Println("❌ Failed to generate a subject.")
			return
		}

		// 2. Show the user the message
    fmt.Printf("\nDrafting Commit:\nTitle: %s\nBody:  %s\n", subject, body)

    // Execute: git commit -m "Subject" -m "Body"
    commitExec := exec.Command("git", "commit", "-m", subject, "-m", body)
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

func askOllamaForCommit(diff string) (string, string) {
	url := fmt.Sprintf("%s/api/generate", AppConfig.Ollama.URL)

  prompt := `Analyze these git changes and write a two-part commit message.
  1. SUBJECT: A one-line summary (max 50 chars) using Conventional Commits (e.g., "feat: add yaml config").
  2. BODY: A concise paragraph explaining "why" the change was made and what was affected.
  
  Output your response in this format:
  SUBJECT: <subject line>
  BODY: <body paragraph>
  
  DIFF:
  ` + diff

	payload, _ := json.Marshal(map[string]interface{}{
		"model":  AppConfig.Ollama.Model,
		"prompt": prompt,
		"stream": false,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

  var res struct{ Response string }
	json.NewDecoder(resp.Body).Decode(&res)

	// Simple parser to split the AI response
	lines := strings.Split(res.Response, "\n")
	var subject, body string
	for _, line := range lines {
		if strings.HasPrefix(line, "SUBJECT:") {
			subject = strings.TrimSpace(strings.TrimPrefix(line, "SUBJECT:"))
		} else if strings.HasPrefix(line, "BODY:") {
			body = strings.TrimSpace(strings.TrimPrefix(line, "BODY:"))
		}
	}

	return subject, body

}
