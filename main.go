package main

import (
	"log"
	"path/filepath"

	"github.com/rqpt/obsidivim/internal/config"
	"github.com/rqpt/obsidivim/internal/mode"
	"github.com/rqpt/obsidivim/internal/note"
	"github.com/rqpt/obsidivim/internal/template"
)

type menuState int

const (
	stateSelectMode menuState = iota
	stateSelectFinal
)

func main() {
	cfg := config.Load()

	state := stateSelectMode

	var (
		modeSelection string
		err           error
	)

	if *cfg.NoteType != "" {
		templatePath, ok := template.Map[*cfg.NoteType]
		if ok {
			note.CreateNew(
				cfg,
				*cfg.NoteType,
				filepath.Join(cfg.TemplatesDir, templatePath),
			)

			return
		} else {
			state = stateSelectFinal
			modeSelection = "Existing"
		}
	}

	for {
		switch state {

		case stateSelectMode:
			modeSelection, err = mode.Select()
			if err != nil || modeSelection == "" {
				log.Fatalf("Mode selection failed: %v", err)
			}

			state = stateSelectFinal

		case stateSelectFinal:
			if modeSelection == "New" {
				templateSelection, templatePath, err := template.Select(cfg.TemplatesDir)
				if err != nil {
					log.Fatalf("Template selection failed: %v", err)
				}
				if templatePath == "" {
					state = stateSelectMode
					continue
				}

				note.CreateNew(cfg, templateSelection, templatePath)

				return
			} else if modeSelection == "Existing" {
				selectedNote, err := note.SelectExisting(cfg)
				if err != nil {
					log.Fatalf("Failed selecting a note: %v", err)
				}
				if selectedNote == "" {
					state = stateSelectMode
					continue
				}

				if err := note.EditExisting(cfg, selectedNote); err != nil {
					log.Fatalf("Failed to edit existing note: %v", err)
				}

				return
			} else {
				return
			}
		}
	}
}
