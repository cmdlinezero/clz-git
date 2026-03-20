package main

import (
  "fmt"
  "git-back/cmd"
)

func main() {

  // Load User configure
  if err := cmd.LoadConfig(); err != nil {
		fmt.Printf("Warning: Could not load config: %v\n", err)
	}

	cmd.Execute()
}
