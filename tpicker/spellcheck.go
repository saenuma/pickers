package main

import (
	"image"
	"slices"
	"strings"
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func isWordInDict(toFindWord string) bool {
	if toFindWord == "a" || toFindWord == "I" || toFindWord == "A" {
		return true
	}

	toFindWord = strings.ToLower(toFindWord)
	if strings.HasSuffix(toFindWord, ".") || strings.HasSuffix(toFindWord, ",") {
		toFindWord = toFindWord[:len(toFindWord)-1]
	}

	smallList, ok := wordsByFirstCharMap[string(toFindWord[0])]
	if !ok {
		return false
	}
	return slices.Contains(smallList, toFindWord)
}

func findWordsNotInDict(inText string) []SpellCheckState {
	inTextWords := strings.Fields(inText)

	var wg sync.WaitGroup
	ret := make([]SpellCheckState, len(inTextWords))
	for i, word := range inTextWords {
		wg.Add(1)
		go func(word string, i int) {
			defer wg.Done()

			retBool := isWordInDict(word)
			ret[i] = SpellCheckState{word, retBool}
		}(word, i)
	}
	wg.Wait()

	return ret
}

func getDisplayWithErrors(window *glfw.Window) image.Image {
	// wWidth, wHeight := window.GetSize()
	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

	sIRect := objCoords[MajorTextInput]
	theCtx.drawTextInputWithErrors(MajorTextInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

	return theCtx.ggCtx.Image()
}
