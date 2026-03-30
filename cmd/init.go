package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/spf13/cobra"
)


var initCmd = &cobra.Command{
	Use:   "init [repo-url]",
	Short: "Init a bare repository with worktree support",
	Long:  "Initialise a bare repository with worktree support.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoURL := args[0]
		
		fmt.Println("📦 Cloning bare repo into .bare...")
		exec.Command("git", "clone", "--bare", repoURL, ".bare").Run()

		fmt.Println("🔗 Creating .git file pointer...")
		os.WriteFile(".git", []byte("gitdir: ./.bare"), 0644)

		fmt.Println("⚙️ Configuring fetch refspecs...")
		exec.Command("git", "config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*").Run()
		
		fmt.Println("🌿 Adding main worktree...")
		exec.Command("git", "worktree", "add", "main").Run()
	},
}

func init() {
	// Centralized command registration
	rootCmd.AddCommand(initCmd)
}
