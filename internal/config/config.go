package config

import (
	"flag"
	"log"
	"os"
)

type Config struct {
	VaultDir     string
	TemplatesDir string
	NoteType     *string
}

func Load() Config {
	typePtr := flag.String(
		"type",
		"",
		"an optional type of note - Fleeting, Research, Blog, Existing",
	)
	flag.Parse()

	return Config{
		VaultDir:     getRequiredEnv("OBSIDIAN_VAULT_DIR"),
		TemplatesDir: getRequiredEnv("OBSIDIAN_TEMPLATES_DIR"),
		NoteType:     typePtr,
	}
}

func getRequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable $%s is not set.", key)
	}
	return val
}
