package note

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(vaultDir, title string) string {
	if !strings.HasSuffix(title, ".md") {
		title += ".md"
	}
	return filepath.Join(vaultDir, title)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Copy(src, dst string) error {
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
