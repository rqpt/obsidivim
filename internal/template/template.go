package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rqpt/picker"
)

var Map = map[string]string{
	"Fleeting": "fleeting_note_template.md",
	"Lecture":  "lecture_note_template.md",
	"Source":   "source_note_template.md",
}

func Select(templatesDir string) (string, string, error) {
	options := make([]string, 0, len(Map))
	for name := range Map {
		options = append(options, name)
	}

	selected, err := picker.Run(options)
	if err != nil || selected == "" {
		return "", "", err
	}

	fileName, ok := Map[selected]
	if !ok {
		return "", "", fmt.Errorf("invalid template option selected: %s", selected)
	}

	return selected, filepath.Join(templatesDir, fileName), nil
}

func ProcessAndSave(src, dst, captureText string) error {
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
