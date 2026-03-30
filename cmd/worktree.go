package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)


var addCmd = &cobra.Command{
	Use:   "add [branch-name]",
	Short: "Add a new feature worktree folder",
  Long:  "This is the command you use daily to create a new folder for a feature e.g. git checkout -b.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		
		fmt.Printf("🌿 Creating worktree for branch: %s\n", branchName)
		
		// This runs: git worktree add <branchName>
		// If the branch doesn't exist, Git will create it based on the folder name
		out, err := exec.Command("git", "worktree", "add", branchName).CombinedOutput()
		
		if err != nil {
			fmt.Printf("❌ Failed to add worktree: %s\n%s", err, string(out))
			return
		}

		fmt.Printf("✨ Worktree created in folder: ./%s\n", branchName)
		fmt.Printf("👉 Run: cd %s\n", branchName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
