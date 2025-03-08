package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/essentialkaos/ek/v13/spellcheck"
)

const (
	FPS                  = 24
	FontSize             = 25
	LineSpacing          = 10
	Margin               = 10
	MaxCaretDisplayCount = FPS / 3

	MajorTextInput = 102

	SWD_CloseBtn = 201
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

	wordsByFirstCharMap map[string][]string
	spellcheckModel     *spellcheck.Model

	currentSuggestions      []string
	suggestionsDialogShown  bool
	currentLineClicked      int
	currentRightClickedWord string
)

type SpellCheckState struct {
	Word   string
	Passed bool
	Minor  bool
}
