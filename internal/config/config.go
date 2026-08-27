package config

import (
	"log"
	"os"
)

type Config struct {
	VaultDir     string
	TemplatesDir string
}

func Load() Config {
	return Config{
		VaultDir:     getRequiredEnv("OBSIDIAN_VAULT_DIR"),
		TemplatesDir: getRequiredEnv("OBSIDIAN_TEMPLATES_DIR"),
	}
}

func getRequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Environment variable $%s is not set.", key)
	}
	return val
}
