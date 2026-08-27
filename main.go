package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rqpt/picker"
)

var templates = map[string]string{
	"Fleeting": "fleeting_note_template.md",
	"Lecture":  "lecture_note_template.md",
	"Source":   "source_note_template.md",
}

func main() {
	vaultDir := getRequiredEnv("OBSIDIAN_VAULT_DIR")
	templatesDir := getRequiredEnv("OBSIDIAN_TEMPLATES_DIR")

	selectedOption, templatePath, err := selectTemplate(templatesDir)
	if err != nil {
		log.Fatalf("Template selection failed: %v", err)
	}
	if templatePath == "" {
		return
	}

	inputText, err := captureInEditor()
	if err != nil {
		log.Fatalf("Failed capturing input in editor: %v", err)
	}
	if inputText == "" {
		log.Fatal("No text saved in editor. Note creation canceled.")
	}

	destPath := resolveDestPath(vaultDir, inputText)
	if fileExists(destPath) {
		log.Fatalf("File already exists: %s", destPath)
	}

	if selectedOption == "Fleeting" {
		if err := processAndCopyTemplate(templatePath, destPath, inputText); err != nil {
			log.Fatalf("Failed to create fleeting note: %v", err)
		}
		return // Skip opening editor again
	}

	// Non-fleeting templates: standard copy & open final note in editor
	if err := copyFile(templatePath, destPath); err != nil {
		log.Fatalf("Failed to create note: %v", err)
	}

	if err := openInEditor(destPath); err != nil {
		log.Fatalf("Failed to open editor: %v", err)
	}
}

func selectTemplate(templatesDir string) (string, string, error) {
	options := make([]string, 0, len(templates))
	for name := range templates {
		options = append(options, name)
	}

	selected, err := picker.Run(options)
	if err != nil || selected == "" {
		return "", "", err
	}

	fileName, ok := templates[selected]
	if !ok {
		return "", "", fmt.Errorf("invalid template option selected: %s", selected)
	}

	return selected, filepath.Join(templatesDir, fileName), nil
}

func captureInEditor() (string, error) {
	tmpFile, err := os.CreateTemp("", "obsidian-title-*.txt")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := openInEditor(tmpPath); err != nil {
		return "", err
	}

	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("reading captured text: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

func processAndCopyTemplate(src, dst, captureText string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading template: %w", err)
	}

	processed := strings.ReplaceAll(string(content), "{{VALUE}}", captureText)

	if err := os.WriteFile(dst, []byte(processed), 0o644); err != nil {
		return fmt.Errorf("writing processed note: %w", err)
	}

	return nil
}

func resolveDestPath(vaultDir, title string) string {
	if !strings.HasSuffix(title, ".md") {
		title += ".md"
	}
	return filepath.Join(vaultDir, title)
}

func openInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening template: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("creating note file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copying template content: %w", err)
	}
	return dstFile.Sync()
}

func getRequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable $%s is not set.", key)
	}
	return val
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
