package main

import "git-back/cmd"

func main() {
  // Load the YAML (or defaults) before doing anything else
  cmd.LoadConfig()

  // Start the Cobra CLI
	cmd.Execute()
}
