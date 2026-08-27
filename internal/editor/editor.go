package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Open(filePath string) error {
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

func CaptureInput() (string, error) {
	tmpFile, err := os.CreateTemp("", "obsidian-title-*.txt")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := Open(tmpPath); err != nil {
		return "", err
	}

	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("reading captured text: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}
