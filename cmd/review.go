package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
  "path/filepath"
  "strings"
  "time"

	"github.com/spf13/cobra"
)

var includeStaticSite bool

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "AI Code Review of staged changes",
	Long:  "Analyzes staged changes and generates a Markdown report with suggestions for improvements, security, and style.",
	Run: func(cmd *cobra.Command, args []string) {

    var finalOutput string

    var diff []byte
    var reviewTarget string

    // Add Intelligence to generate REVIEW 

    // 1. Try Staged Changes
    diff, _ = exec.Command("git", "diff", "--cached").Output()
    if len(diff) > 0 {
        reviewTarget = "Staged Changes"
    }

    // 2. Fallback to Unstaged Changes (WIP)
    if len(diff) == 0 {
        diff, _ = exec.Command("git", "diff").Output()
        if len(diff) > 0 {
            reviewTarget = "Unstaged Changes (WIP)"
        }
    }

    // 3. Fallback to Last Commit
    if len(diff) == 0 {
        fmt.Println("📝 No local changes found. Reviewing the last commit...")
        diff, _ = exec.Command("git", "diff", "HEAD~1", "HEAD").Output()
        reviewTarget = "Last Commit (HEAD)"
    }

    fmt.Printf("🔍 AI is reviewing: %s\n", reviewTarget)

		// 1. Get staged changes (consistent with commit command)
		diff, _ = exec.Command("git", "diff", "--cached").Output()
		if len(diff) == 0 {
			fmt.Println("⚠️  No staged changes to review. Run 'git add' first.")
			return
		}

    // Get Repo Name
    outName, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
    repoName := filepath.Base(strings.TrimSpace(string(outName)))

		// Get the Current Short SHA
		shaOut, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
		shortSHA := strings.TrimSpace(string(shaOut))

    // Get Current Branch for the Description
    branchOut, _ := exec.Command("git", "branch", "--show-current").Output()
    branchName := strings.TrimSpace(string(branchOut))

		// Prepare the dynamic fields
		// Title: clz-git (a1b2c3d)
		dynamicTitle := fmt.Sprintf("%s (%s)", repoName, shortSHA)
    description := fmt.Sprintf("Code Review for branch: %s", branchName)

		fmt.Println("🔍 AI is reviewing your code...")
		report, err := callOllamaForReview(string(diff))
		if err != nil {
			fmt.Printf("❌ Review failed: %v\n", err)
			return
		}

		// 2. Determine filename based on current SHA or "draft"
		shaByte, _ := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
		fileName := fmt.Sprintf("review-%s.md", bytes.TrimSpace(shaByte))

    // Apply Hugo enrichment
    if includeStaticSite {
        finalOutput = getHugoFrontmatter(dynamicTitle, description) + report
    } else {
        finalOutput = report
    }

		err = os.WriteFile(fileName, []byte(finalOutput), 0644)
		if err != nil {
			fmt.Printf("❌ Could not save report: %v\n", err)
			return
		}

		fmt.Printf("✨ Review complete! Report saved to: %s\n", fileName)
	},
}

func init() {
  reviewCmd.Flags().BoolVarP(&includeStaticSite, "site", "s", false, "Include Hugo Static Site frontmatter in the output")
	rootCmd.AddCommand(reviewCmd)
}

func getHugoFrontmatter(repoName string, description string) string {
	currentDate := time.Now().Format("2006-01-02")
	
	// Clean up description (ensure it's one line and quoted for YAML)
	if description == "" {
		description = "Code Review"
	}

	return fmt.Sprintf(`---
title: "%s"
description: "%s"
date: %s
category: "Developer Journal"
tags: [ "review", "security" ]
duration: 5:00
draft: true
hero_title: "git-back"
hero_image: "/images/labdemo-logo.png"
image: "/images/labdemo-logo.png"
issue: 4
volume: 1
special_edition: false
---

`, repoName, description, currentDate)
}


func callOllamaForReview(diff string) (string, error) {
	url := fmt.Sprintf("%s/api/generate", AppConfig.Ollama.URL)

	prompt := `You are an expert Senior Software Engineer. Review the following git diff.
    Provide a Markdown report with the following sections:
    1. ## Summary: Brief overview of the changes.
    2. ## Logic & Efficiency: Potential bugs or performance bottlenecks.
    3. ## Security: Any visible vulnerabilities (hardcoded secrets, injection, etc.).
    4. ## Style & Readability: Suggestions for cleaner code.

    DIFF:
    ` + diff

	payload, _ := json.Marshal(map[string]interface{}{
		"model":  AppConfig.Ollama.Model,
		"prompt": prompt,
		"stream": false,
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("Ollama offline")
	}
	defer resp.Body.Close()

	var res struct{ Response string }
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Response, nil
}
