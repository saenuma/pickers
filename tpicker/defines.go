package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/essentialkaos/ek/v13/spellcheck"
)

const (
	FPS                  = 24
	FontSize             = 30
	LineSpacing          = 10
	Margin               = 10
	MaxCaretDisplayCount = FPS / 3

	MajorTextInput = 102
)

var (
	objCoords   map[int]g143.Rect
	scObjCoords map[int]g143.Rect

	currentWindowFrame    image.Image
	windowFrameWithErrors image.Image

	enteredTxt        string
	caretX            = Margin
	caretY            = Margin
	caretDisplayed    bool
	caretDisplayCount int
	frameUpdated      bool

	wordsByFirstCharMap    map[string][]string
	spellcheckModel        *spellcheck.Model
	currentSuggestions     []string
	suggestionsDialogShown bool
)

type SpellCheckState struct {
	Word   string
	Passed bool
	Minor  bool
}
