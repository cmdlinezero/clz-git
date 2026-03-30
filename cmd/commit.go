package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"os/exec"
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

		// fmt.Println("🤖 Analyzing changes and drafting message...")
		// subject, body := askOllamaForCommit(string(diff))

		// if subject == "" {
		// 	fmt.Println("❌ Failed to generate a subject.")
		// 	return
		// }

		fmt.Println("🤖 Requesting AI commit message...")
		subject, body, err := callOllama(string(diff))

    if err != nil {
			fmt.Printf("❌ AI Error: %v\n", err)
			fmt.Println("👉 Please start Ollama ('ollama serve') or perform a manual commit:")
			fmt.Println("   git commit -m \"your message\"")
			return // EXIT HERE:
		}
		
		fmt.Printf("\n📝 AI Suggestion:\nsubject: %s\nbody: %s\n", subject, body)

    // Execute the actual git commit
		commitExec := exec.Command("git", "commit", "-m", subject, "-m", body)
		commitExec.Stdout = os.Stdout
		commitExec.Stderr = os.Stderr

    if err := commitExec.Run(); err != nil {
			fmt.Printf("❌ Git commit failed: %v\n", err)
			return
		}

		fmt.Println("✅ Committed successfully!")
	},
}

func init() {
	// Centralized command registration
	rootCmd.AddCommand(commitCmd)
}

func callOllama(diff string) (string, string, error) {
  apiEndpoint := AppConfig.Ollama.URL + "/api/generate"

  prompt := `Analyze these git changes and write a two-part commit message.
  1. SUBJECT: A one-line summary (max 50 chars) using Conventional Commits (e.g., "feat: add yaml config").
  2. BODY: A concise paragraph explaining "why" the change was made and what was affected.
  
  Output your response in this format:
  SUBJECT: <subject line>
  BODY: <body paragraph>
  
  DIFF:
  ` + diff

//	prompt := "Write a concise Conventional Commit message for this diff. No conversational filler:\n\n" + diff
	
	payload, _ := json.Marshal(map[string]interface{}{
		"model": AppConfig.Ollama.Model, "prompt": prompt, "stream": false,
	})


  fmt.Printf("URL: %s\n", apiEndpoint)
  fmt.Printf("MODEL: %s\n", AppConfig.Ollama.Model)

	resp, err := http.Post(apiEndpoint, "application/json", bytes.NewBuffer(payload))

	if err != nil {
    // Check if the error is specifically a connection failure
    fmt.Printf("⚠️  Ollama appears to be offline at %s\n", AppConfig.Ollama.Model)
    fmt.Println("👉 Check if Ollama is running: 'ollama serve'")
		return "Ollama Error", "❌ fix: Manual update required (Ollama Offline)", err
	}
	
	defer resp.Body.Close()

	var res struct{ Response string `json:"response"` }
  if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		// return "", fmt.Errorf("failed to decode AI response")
    return "", "", err

	}

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

 	return subject, body, nil
}
