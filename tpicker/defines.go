package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
)

const (
	FPS                  = 24
	FontSize             = 30
	LineSpacing          = 10
	MaxCaretDisplayCount = FPS / 2

	DoneBtn        = 101
	MajorTextInput = 102
)

var (
	objCoords          map[int]g143.Rect
	currentWindowFrame image.Image
	enteredTxt         string
	caretX             = 10
	caretY             = 10
	caretDisplayed     bool
	caretDisplayCount  int
)
