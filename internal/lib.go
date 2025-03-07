package internal

import (
	"os"
	"path/filepath"
)

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "picker_font.ttf")
	if !DoesPathExists(fontPath) {
		os.WriteFile(fontPath, DefaultFont, 0777)
	}
	return fontPath
}
