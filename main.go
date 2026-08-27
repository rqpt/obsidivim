package main

import (
	"log"

	"github.com/rqpt/obsidivim/internal/config"
	"github.com/rqpt/obsidivim/internal/editor"
	"github.com/rqpt/obsidivim/internal/note"
	"github.com/rqpt/obsidivim/internal/template"
)

func main() {
	cfg := config.Load()

	selectedOption, templatePath, err := template.Select(cfg.TemplatesDir)
	if err != nil {
		log.Fatalf("Template selection failed: %v", err)
	}
	if templatePath == "" {
		return
	}

	inputText, err := editor.CaptureInput()
	if err != nil {
		log.Fatalf("Failed capturing input in editor: %v", err)
	}
	if inputText == "" {
		log.Fatal("No text saved in editor. Note creation canceled.")
	}

	destPath := note.ResolvePath(cfg.VaultDir, inputText)
	if note.Exists(destPath) {
		log.Fatalf("File already exists: %s", destPath)
	}

	if selectedOption == "Fleeting" {
		if err := template.ProcessAndSave(templatePath, destPath, inputText); err != nil {
			log.Fatalf("Failed to create fleeting note: %v", err)
		}
		return
	}

	if err := note.Copy(templatePath, destPath); err != nil {
		log.Fatalf("Failed to create note: %v", err)
	}

	if err := editor.Open(destPath); err != nil {
		log.Fatalf("Failed to open editor: %v", err)
	}
}
