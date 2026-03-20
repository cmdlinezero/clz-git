package cmd

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)

var wtAddCmd = &cobra.Command{
	Use:   "add [branch-name]",
	Short: "Create a new feature worktree folder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branch := args[0]
		
		// In a bare setup, we want the folder name to match the branch name
		// This command creates the branch AND the folder simultaneously
		fmt.Printf("🏗 Adding worktree for branch '%s'...\n", branch)
		
		// git worktree add <path> <branch>
		// By using 'branch' as both, we get ./my-feature/ containing the my-feature branch
		c := exec.Command("git", "worktree", "add", "-b", branch, branch)
		
		out, err := c.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Git Error: %s\n", string(out))
			return
		}
		
		fmt.Printf("✅ Worktree created! Run: cd %s\n", branch)
	},
}

func init() {
	rootCmd.AddCommand(wtAddCmd)
}
