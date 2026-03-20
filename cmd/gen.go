package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)

var isLongForm bool

func init() {
	// Add the flag to the changelog subcommand
	changelogCmd.Flags().BoolVarP(&isLongForm, "long", "l", false, "Use long form git log with file statuses")
	genCmd.AddCommand(changelogCmd)
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate documentation from git history",
}

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate a Markdown changelog",
	Run: func(cmd *cobra.Command, args []string) {
		var gitArgs []string

		if isLongForm {
			fmt.Println("🔍 Extracting detailed history (Long Form)...")
			gitArgs = []string{"log", "-n", "5", "--pretty=format:Commit: %h %nAuthor: %an %nDate: %ar %nMessage: %s", "--name-status"}
		} else {
			fmt.Println("📜 Extracting summary history (Short Form)...")
			gitArgs = []string{"log", "-n", "10", "--pretty=format:%s %cr <%an>"}
		}

		logs, err := exec.Command("git", gitArgs...).Output()
		if err != nil {
			fmt.Println("❌ Error fetching git logs:", err)
			return
		}

		fmt.Println("🤖 Ollama is drafting the markdown...")
		markdown := askOllamaForMarkdown(string(logs), isLongForm)

		os.WriteFile("CHANGELOG_DRAFT.md", []byte(markdown), 0644)
		fmt.Println("✨ Done! Created CHANGELOG_DRAFT.md")
	},
}

func askOllamaForMarkdown(logs string, detailed bool) string {
	url := "http://localhost:11434/api/generate"
	
	systemPrompt := "You are a technical writer. Summarize these git logs into a clean Markdown changelog."
	if detailed {
		systemPrompt += " Use the file change information to explain *what* parts of the system were affected."
	}

	prompt := fmt.Sprintf("%s\n\nLOGS:\n%s", systemPrompt, logs)

	payload, _ := json.Marshal(map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "Error: Ollama connection failed."
	}
	defer resp.Body.Close()

	var res struct{ Response string }
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Response
}
