package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "git-back",
	Short: "AI-powered Git workflow manager",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	// Centralized command registration
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(wtCmd)
	rootCmd.AddCommand(commitCmd)
}
