package main

import (
	"slices"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

	getLineClicked := func() int {
		txtParts := strings.Split(enteredTxt, "\n")
		for i := range len(txtParts) {
			if yPosInt > Margin+(i*(FontSize+LineSpacing)) && yPosInt < Margin+((i+1)*(FontSize+LineSpacing)) {
				return i
			}
		}
		return len(txtParts)
	}

	getSubTextFromClick := func(theCtx Ctx, lineNo, xPosInt int) string {
		txtParts := strings.Split(enteredTxt, "\n")
		lineClicked := txtParts[lineNo]

		for i := range len(lineClicked) - 1 {
			subLineW, _ := theCtx.ggCtx.MeasureString(lineClicked[:i])
			subLineW2, _ := theCtx.ggCtx.MeasureString(lineClicked[:i+1])
			if xPosInt >= int(subLineW) && xPosInt <= int(subLineW2) {
				return lineClicked[:i]
			}
		}
		return lineClicked
	}

	getWordRightClicked := func(theCtx Ctx, lineNo, xPosInt int) string {
		txtParts := strings.Split(enteredTxt, "\n")
		if lineNo >= len(txtParts) {
			return ""
		}

		lineClicked := txtParts[lineNo]

		clickedWordIndex := 0
		for i := range len(lineClicked) {
			subLineW, _ := theCtx.ggCtx.MeasureString(lineClicked[:i])
			subLineW2, _ := theCtx.ggCtx.MeasureString(lineClicked[:i+1])
			if xPosInt > int(subLineW) && xPosInt < int(subLineW2) {
				clickedWordIndex = i
			}
		}

		lineClickedW, _ := theCtx.ggCtx.MeasureString(lineClicked)
		if clickedWordIndex == 0 && xPosInt > int(lineClickedW) {
			return ""
		}

		prevChars := make([]string, 0)
		prevCharIndex := clickedWordIndex
		for {
			currentChar := string(lineClicked[prevCharIndex])
			prevChars = append(prevChars, currentChar)
			if currentChar == " " {
				break
			}

			prevCharIndex -= 1
			if prevCharIndex < 0 {
				break
			}
		}
		slices.Reverse(prevChars)

		nextChars := make([]string, 0)
		nextCharIndex := clickedWordIndex

		for {
			currentChar := string(lineClicked[nextCharIndex])
			nextChars = append(nextChars, currentChar)
			if currentChar == " " {
				break
			}

			nextCharIndex += 1
			if nextCharIndex >= len(lineClicked) {
				break
			}
		}

		foundChars := slices.Concat(prevChars, nextChars[1:])
		return strings.TrimSpace(strings.Join(foundChars, ""))
	}

	// wWidth, wHeight := window.GetSize()
	if button == glfw.MouseButtonLeft {
		txtParts := strings.Split(enteredTxt, "\n")
		maxCaretY := Margin + ((len(txtParts) - 1) * (FontSize + Margin))
		lastLineW, _ := theCtx.ggCtx.MeasureString(txtParts[len(txtParts)-1])

		if yPosInt > maxCaretY {
			caretY = maxCaretY
			if xPosInt > int(lastLineW) {
				caretX = Margin + int(lastLineW)
			} else {
				lineClicked := len(txtParts) - 1
				tmpSubText := getSubTextFromClick(theCtx, lineClicked, xPosInt)
				tmpSubTextW, _ := theCtx.ggCtx.MeasureString(tmpSubText)
				caretX = int(tmpSubTextW) + Margin
			}
		} else {
			lineClicked := getLineClicked()
			lineClickedW, _ := theCtx.ggCtx.MeasureString(txtParts[lineClicked])

			if xPosInt > int(lineClickedW) {
				caretX = int(lineClickedW) + Margin
			} else {
				tmpSubText := getSubTextFromClick(theCtx, lineClicked, xPosInt)
				tmpSubTextW, _ := theCtx.ggCtx.MeasureString(tmpSubText)
				caretX = int(tmpSubTextW) + Margin
			}
			caretY = Margin + (lineClicked * (FontSize + Margin))
		}

	} else if button == glfw.MouseButtonRight {
		lineClicked := getLineClicked()
		rightClickedWord := getWordRightClicked(theCtx, lineClicked, xPosInt)

		if !isWordInDict(rightClickedWord) {
			suggestions := spellcheckModel.Suggest(rightClickedWord, 30)

			currentLineClicked = lineClicked
			currentRightClickedWord = rightClickedWord
			currentSuggestions = suggestions

			// bring up suggestions
			suggestionsDialogShown = true
			drawDialog(window, suggestions)
			window.SetMouseButtonCallback(sWMouseBtnCallback)
			// window.SetCursorPosCallback(cursorCallback)
			window.SetKeyCallback(nil)
			window.SetCharCallback(nil)
		}

	}

}

func sWMouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range scObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	if widgetCode == SWD_CloseBtn {
		// move to text view without applying any suggestion
		drawTextView(window)
		window.SetMouseButtonCallback(mouseBtnCallback)
		// window.SetCursorPosCallback(cursorCallback)
		window.SetKeyCallback(mKeyCallback)
		window.SetCharCallback(mCharCallback)
		window.SetCursorPosCallback(nil)

		suggestionsDialogShown = false
	}

	if widgetCode > 1000 && widgetCode < 2000 {
		num := widgetCode - 1000 - 1
		suggestedWord := currentSuggestions[num]
		enteredTxtParts := strings.Split(enteredTxt, "\n")
		lineClickedStr := enteredTxtParts[currentLineClicked]

		incorrectWordBegin := strings.Index(lineClickedStr, currentRightClickedWord)
		incorrectWordEnd := incorrectWordBegin + len(currentRightClickedWord)
		if incorrectWordBegin != -1 {
			editedLineStr := lineClickedStr[:incorrectWordBegin] + suggestedWord + lineClickedStr[incorrectWordEnd:]
			enteredTxtParts[currentLineClicked] = editedLineStr
			enteredTxt = strings.Join(enteredTxtParts, "\n")
		}

		// move to text view
		drawTextView(window)
		window.SetMouseButtonCallback(mouseBtnCallback)
		// window.SetCursorPosCallback(cursorCallback)
		window.SetKeyCallback(mKeyCallback)
		window.SetCharCallback(mCharCallback)
		window.SetCursorPosCallback(nil)

		suggestionsDialogShown = false
	}
}
