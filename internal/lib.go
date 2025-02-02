package internal

import (
	"os"
	"path/filepath"
)

func GetDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "picker_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}
