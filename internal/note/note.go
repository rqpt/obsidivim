package note

import (
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

func CreateNew(cfg config.Config, templateSelection string, templatePath string) {
	inputText, err := editor.CaptureInput()
	if err != nil {
		log.Fatalf("Failed capturing input in editor: %v", err)
	}
	if inputText == "" {
		log.Fatal("No text saved in editor. Note creation canceled.")
	}

	destPath := resolvePath(cfg.VaultDir, inputText)
	if exists(destPath) {
		log.Fatalf("File already exists: %s", destPath)
	}

	if templateSelection == "Fleeting" {
		if err := template.ProcessAndSave(templatePath, destPath, inputText); err != nil {
			log.Fatalf("Failed to create fleeting note: %v", err)
		}

		return
	}

	if err := copy(templatePath, destPath); err != nil {
		log.Fatalf("Failed to create note: %v", err)
	}

	if err := editor.Open(destPath); err != nil {
		log.Fatalf("Failed to open editor: %v", err)
	}
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
