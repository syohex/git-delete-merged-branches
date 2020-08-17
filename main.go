package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	os.Exit(_main())
}

func isIgnoredBranch(name string) bool {
	if name == "master" || name == "develop" {
		return true
	}

	// Ignore release branches
	if strings.HasPrefix(name, "release/") {
		return true
	}

	return false
}

func gitMergedBranches() ([]string, error) {
	var output bytes.Buffer
	cmd := exec.Command("git", "branch", "--merged")
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("could not get merged branches: %w", err)
	}

	var ret []string
	scanner := bufio.NewScanner(&output)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") || isIgnoredBranch(line) {
			// ignore current branch
			continue
		}

		ret = append(ret, line)
	}

	return ret, nil
}

func deleteGitBranches(branches []string) {
	for _, branch := range branches {
		cmd := exec.Command("git", "branch", "-d", branch)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			continue
		}
	}
}

func _main() int {
	branches, err := gitMergedBranches()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if len(branches) == 0 {
		fmt.Println("There is no merged branch")
		return 0
	}

	deleteGitBranches(branches)
	return 0
}
