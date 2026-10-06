package note

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rqpt/editor"
	"github.com/rqpt/obsidivim/internal/config"
	"github.com/rqpt/obsidivim/internal/template"
	"github.com/rqpt/picker"
)

func SelectExisting(cfg config.Config) (string, error) {
	notes, err := picker.ListFiles(cfg.VaultDir, []string{".md"})
	if err != nil {
		log.Fatalf("Failed listing notes in vault path: %v", err)
	}

	return picker.Run(notes)
}

func EditExisting(cfg config.Config, selectedNote string) error {
	notePath := filepath.Join(cfg.VaultDir, selectedNote)

	return editor.Open(notePath)
}

func CreateNew(cfg config.Config, templateSelection string, templatePath string) error {
	inputText, err := editor.CaptureInput()
	if err != nil {
		return fmt.Errorf("Failed capturing input in editor: %w", err)
	}
	if inputText == "" {
		return errors.New("filename is empty")
	}

	destPath := resolvePath(cfg.VaultDir, inputText)
	if exists(destPath) {
		return fmt.Errorf("File already exists: %s", destPath)
	}

	if templateSelection == "Fleeting" {
		if err := template.ProcessAndSave(templatePath, destPath, inputText); err != nil {
			return fmt.Errorf("Failed to create fleeting note: %v", err)
		}

		return nil
	}

	if err := copy(templatePath, destPath); err != nil {
		return fmt.Errorf("Failed to create note: %w", err)
	}

	if err := editor.Open(destPath); err != nil {
		return fmt.Errorf("Failed to open editor: %w", err)
	}

	return nil
}

func resolvePath(vaultDir, title string) string {
	if !strings.HasSuffix(title, ".md") {
		title += ".md"
	}
	return filepath.Join(vaultDir, title)
}

func copy(src, dst string) error {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
