package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultRepoURL = "https://github.com/brightsidedeveloper/go-native-template"

type Config struct {
	ProjectName string
	TargetDir   string
}

func main() {
	var (
		projectName = flag.String("name", "", "Project name (required)")
		targetDir   = flag.String("dir", "", "Target directory name (defaults to project name)")
	)
	flag.Parse()

	if *projectName == "" {
		fmt.Fprintf(os.Stderr, "Error: --name is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s --name <project-name> [--dir <target-dir>]\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	config := Config{
		ProjectName: *projectName,
		TargetDir:   *targetDir,
	}

	// Default target directory to project name if not specified
	if config.TargetDir == "" {
		config.TargetDir = config.ProjectName
	}

	fmt.Printf("Creating new project: %s\n", config.ProjectName)
	fmt.Printf("  Target Directory: %s\n", config.TargetDir)
	fmt.Println()

	// Clone repository
	if err := cloneRepository(defaultRepoURL, config.TargetDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error cloning repository: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Repository cloned to %s\n", config.TargetDir)

	// Template the project
	if err := templateProject(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Project templated successfully!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. cd %s\n", config.TargetDir)
	fmt.Println("  2. cd server && go mod tidy && make gen")
	fmt.Println("  3. cd ../mobile && npm install")
}

func templateProject(config Config) error {
	// Template server files
	if err := templateServer(config); err != nil {
		return fmt.Errorf("templating server: %w", err)
	}

	// Template mobile files
	if err := templateMobile(config); err != nil {
		return fmt.Errorf("templating mobile: %w", err)
	}

	// Template root files
	if err := templateRoot(config); err != nil {
		return fmt.Errorf("templating root: %w", err)
	}

	return nil
}

func templateServer(config Config) error {
	serverDir := filepath.Join(config.TargetDir, "server")
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		return fmt.Errorf("server directory not found: %s", serverDir)
	}

	// Use just the project name as the module path
	newModulePath := config.ProjectName
	replacements := map[string]string{
		// Go module path - replace ALL occurrences
		"github.com/brightsidedeveloper/go-native-template": newModulePath,
		// JWT issuer
		"loop-app": fmt.Sprintf("%s-app", config.ProjectName),
		// Test database name
		"loop_test": fmt.Sprintf("%s_test", config.ProjectName),
		// Email defaults
		"noreply@template.app": fmt.Sprintf("noreply@%s.app", config.ProjectName),
		"Template":             config.ProjectName,
	}

	// Process ALL files recursively in server directory to catch all occurrences
	return filepath.Walk(serverDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Skip binary files and other non-text files
		ext := strings.ToLower(filepath.Ext(path))
		skipExts := map[string]bool{
			".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
			".ico": true, ".pdf": true, ".zip": true, ".tar": true,
			".gz": true, ".exe": true, ".dll": true, ".so": true,
			".dylib": true,
		}
		if skipExts[ext] {
			return nil
		}

		// Process all files including generated ones
		// (Generated files will be regenerated later, but we replace them now
		// so go mod tidy works before running make gen)
		return replaceInFile(path, replacements)
	})
}

func templateMobile(config Config) error {
	mobileDir := filepath.Join(config.TargetDir, "mobile")
	if _, err := os.Stat(mobileDir); os.IsNotExist(err) {
		return fmt.Errorf("mobile directory not found: %s", mobileDir)
	}

	// More precise replacements for JSON files
	// Handle different indentation levels - order matters (most specific first)
	replacements := map[string]string{
		// app.json (4 spaces indentation inside expo object)
		`    "name": "template"`:   fmt.Sprintf(`    "name": "%s"`, config.ProjectName),
		`    "slug": "template"`:   fmt.Sprintf(`    "slug": "%s"`, config.ProjectName),
		`    "scheme": "template"`: fmt.Sprintf(`    "scheme": "%s"`, config.ProjectName),
		// package.json (2 spaces indentation)
		`  "name": "template"`: fmt.Sprintf(`  "name": "%s"`, config.ProjectName),
		// Fallback for any other cases
		`"name": "template"`:   fmt.Sprintf(`"name": "%s"`, config.ProjectName),
		`"slug": "template"`:   fmt.Sprintf(`"slug": "%s"`, config.ProjectName),
		`"scheme": "template"`: fmt.Sprintf(`"scheme": "%s"`, config.ProjectName),
	}

	files := []string{
		"package.json",
		"app.json",
	}

	for _, file := range files {
		filePath := filepath.Join(mobileDir, file)
		if err := replaceInFile(filePath, replacements); err != nil {
			return fmt.Errorf("processing %s: %w", file, err)
		}
	}

	return nil
}

func templateRoot(config Config) error {
	replacements := map[string]string{
		"# Template": fmt.Sprintf("# %s", config.ProjectName),
	}

	files := []string{
		"Readme.md",
	}

	for _, file := range files {
		filePath := filepath.Join(config.TargetDir, file)
		if err := replaceInFile(filePath, replacements); err != nil {
			// Root README is optional
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("processing %s: %w", file, err)
		}
	}

	return nil
}

func replaceInFile(filePath string, replacements map[string]string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	modified := string(content)

	// Process replacements in order from longest to shortest to avoid partial matches
	// Convert map to slice of pairs and sort by length
	type replacement struct {
		old string
		new string
	}
	var sorted []replacement
	for old, new := range replacements {
		sorted = append(sorted, replacement{old: old, new: new})
	}
	// Sort by length descending
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if len(sorted[i].old) < len(sorted[j].old) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Apply replacements in sorted order
	for _, r := range sorted {
		modified = strings.ReplaceAll(modified, r.old, r.new)
	}

	if modified != string(content) {
		return os.WriteFile(filePath, []byte(modified), 0644)
	}

	return nil
}

func cloneRepository(repoURL, targetDir string) error {
	// Check if target directory already exists
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("target directory already exists: %s", targetDir)
	}

	// Clone the repository
	fmt.Printf("Cloning repository from %s...\n", repoURL)
	cmd := exec.Command("git", "clone", repoURL, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Remove .git directory to start fresh
	gitDir := filepath.Join(targetDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		// Non-fatal, just warn
		fmt.Printf("Warning: could not remove .git directory: %v\n", err)
	}

	return nil
}
