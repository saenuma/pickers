package main

import (
	"image"
	"slices"
	"strings"
	"sync"
	"unicode"
)

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func clearQuotes(str string) string {
	return strings.ReplaceAll(str, "'", "")
}

func isWordInDict(toFindWord string) bool {
	if toFindWord == "a" || toFindWord == "I" || toFindWord == "A" {
		return true
	}

	if IsUpper(toFindWord) {
		return true
	}
	if IsUpper(string(toFindWord[0])) {
		return true
	}
	toFindWord = clearQuotes(toFindWord)
	if strings.HasSuffix(toFindWord, ".") || strings.HasSuffix(toFindWord, ",") {
		toFindWord = toFindWord[:len(toFindWord)-1]
	}

	// return spellcheckTrie.SearchDirect(toFindWord)
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
		go func(word string, i int, Words []string) {
			defer wg.Done()

			retBool := isWordInDict(word)
			if retBool {
				sentenceBegin := false
				if i == 0 {
					sentenceBegin = true
				}
				if i > 0 {
					if strings.HasSuffix(Words[i-1], ".") {
						sentenceBegin = true
					}
				}
				if sentenceBegin && IsLower(string(word[0])) {
					ret[i] = SpellCheckState{word, false, true}
				}
			} else {
				ret[i] = SpellCheckState{word, retBool, false}
			}
		}(word, i, inTextWords)
	}
	wg.Wait()

	return ret
}

func getDisplayWithErrors() image.Image {
	// wWidth, wHeight := window.GetSize()
	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

	sIRect := objCoords[MajorTextInput]
	theCtx.drawTextInputWithErrors(MajorTextInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

	return theCtx.ggCtx.Image()
}
