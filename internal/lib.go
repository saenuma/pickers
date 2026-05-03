package internal

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

func GetTextScale() float64 {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "text-scaling-factor")
	textScaleStr, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return 1.0
	}
	textScale, _ := strconv.ParseFloat(strings.TrimSpace(string(textScaleStr)), 64)
	return textScale
}

func GetFontSize() float64 {
	textScale := GetTextScale()
	tmpFontSize := textScale * DefaultFontSize
	return tmpFontSize
}

func NotEqual(a, b float64) bool {
	aInt := int(math.Ceil(a))
	bInt := int(math.Ceil(b))
	return aInt != bInt
}
