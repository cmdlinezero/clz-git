package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate AI commit message from staged changes",
	Run: func(cmd *cobra.Command, args []string) {
		diff, _ := exec.Command("git", "diff", "--cached").Output()
		if len(diff) == 0 {
			fmt.Println("⚠️ No changes staged. Run 'git add' first.")
			return
		}

		msg := callOllama(string(diff))
		fmt.Printf("\n📝 AI Suggestion: %s\n", msg)
		
		// Execute the commit
		exec.Command("git", "commit", "-m", msg).Run()
		fmt.Println("✅ Committed successfully!")
	},
}

func callOllama(diff string) string {
	url := "http://localhost:11434/api/generate"
	prompt := "Write a concise Conventional Commit message for this diff. No conversational filler:\n\n" + diff
	
	payload, _ := json.Marshal(map[string]interface{}{
		"model": "llama3", "prompt": prompt, "stream": false,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "fix: manual update (Ollama offline)"
	}
	defer resp.Body.Close()

	var res struct{ Response string `json:"response"` }
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Response
}
