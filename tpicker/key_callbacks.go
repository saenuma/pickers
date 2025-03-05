package main

import (
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func mCharCallback(window *glfw.Window, char rune) {
	wWidth, wHeight := window.GetSize()
	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

	getCaretYAtLineNumber := func() int {
		enteredTxtParts := strings.Split(enteredTxt, "\n")
		for i := range len(enteredTxtParts) {
			if caretY >= (i*(FontSize+Margin)) && caretY <= ((i+1)*(FontSize+Margin)) {
				return i
			}
		}
		return len(enteredTxtParts) - 1
	}

	getCaretXAtSubText := func(theCtx Ctx, lineNo int) string {
		enteredTxtParts := strings.Split(enteredTxt, "\n")
		lineNoText := enteredTxtParts[lineNo]
		for i := range len(lineNoText) {
			partTxtW, _ := theCtx.ggCtx.MeasureString(lineNoText[:i])
			if caretX == int(partTxtW)+Margin {
				return lineNoText[:i]
			}
		}

		return lineNoText
	}

	maxWidth := wWidth - (2 * Margin)

	enteredTxtParts := strings.Split(enteredTxt, "\n")
	caretYLineNo := getCaretYAtLineNumber()
	currentLine := enteredTxtParts[caretYLineNo]
	currentLineW, _ := theCtx.ggCtx.MeasureString(currentLine)
	caretXAtSubText := getCaretXAtSubText(theCtx, caretYLineNo)
	caretXAtSubTextW, _ := theCtx.ggCtx.MeasureString(caretXAtSubText)

	if caretYLineNo == len(enteredTxtParts)-1 && int(caretXAtSubTextW) == int(currentLineW) {
		// the caret is at the end of the text
		tmpEnteredLine := currentLine + string(char)
		tmpEnteredLineW, _ := theCtx.ggCtx.MeasureString(tmpEnteredLine)
		if int(tmpEnteredLineW) > maxWidth {
			currentLine = currentLine + "\n" + string(char)
			enteredTxtParts[caretYLineNo] = currentLine
			enteredTxt = strings.Join(enteredTxtParts, "\n")
			caretY += FontSize + LineSpacing
			charW, _ := theCtx.ggCtx.MeasureString(string(char))
			caretX = Margin + int(charW)
		} else {
			enteredTxtParts[caretYLineNo] = tmpEnteredLine
			enteredTxt = strings.Join(enteredTxtParts, "\n")
			caretX = Margin + int(tmpEnteredLineW)
		}

	} else {
		// the caret is not at the end of the input.
		subTextLen := len(caretXAtSubText)
		tmpEnteredLine := currentLine[:subTextLen] + string(char) + currentLine[subTextLen:]
		tmpEnteredLineW, _ := theCtx.ggCtx.MeasureString(tmpEnteredLine)
		if int(tmpEnteredLineW) > maxWidth {
			tmpText := currentLine[:subTextLen] + string(char)
			tmpTextW, _ := theCtx.ggCtx.MeasureString(tmpText)
			if int(tmpTextW) > maxWidth {
				enteredLine := currentLine[:subTextLen] + "\n" + string(char) + currentLine[subTextLen:]
				enteredTxtParts[caretYLineNo] = enteredLine
				enteredTxt = strings.Join(enteredTxtParts, "\n")
				charW, _ := theCtx.ggCtx.MeasureString(string(char))
				caretX = Margin + int(charW)
				caretY += FontSize + Margin
			} else {
				enteredLine := currentLine[:subTextLen] + string(char) + "\n" + currentLine[subTextLen:]
				enteredTxtParts[caretYLineNo] = enteredLine
				enteredTxt = strings.Join(enteredTxtParts, "\n")
				caretX = int(tmpTextW) + Margin
			}
		} else {
			enteredTxtParts[caretYLineNo] = tmpEnteredLine
			enteredTxt = strings.Join(enteredTxtParts, "\n")
			subTextW, _ := theCtx.ggCtx.MeasureString(currentLine[:subTextLen] + string(char))
			caretX = Margin + int(subTextW)
		}
	}

	// if caretYLineNo == len(enteredTxtParts)-1 {
	// 	// we are on the last line
	// 	caretXAtSubText := getCaretXAtSubText(theCtx, caretYLineNo)
	// 	lastLine := enteredTxtParts[len(enteredTxtParts)-1]
	// 	lastLineW, _ := theCtx.ggCtx.MeasureString(lastLine)
	// 	caretXAtSubTextW, _ := theCtx.ggCtx.MeasureString(caretXAtSubText)

	// 	if lastLineW == caretXAtSubTextW {
	// 		// enter char at end of text
	// 		tmp := enteredTxt + string(char)
	// 		tmpParts := strings.Split(tmp, "\n")
	// 		var lastLine string
	// 		if len(tmpParts) == 1 {
	// 			lastLine = tmp
	// 		} else {
	// 			lastLine = tmpParts[len(tmpParts)-1]
	// 		}
	// 		lastLineW, _ := theCtx.ggCtx.MeasureString(lastLine)
	// 		if int(lastLineW) >  {
	// 			enteredTxt = enteredTxt + "\n" + string(char)
	// 			caretY += FontSize + LineSpacing
	// 			caretX = Margin
	// 		} else {
	// 			enteredTxt = tmp
	// 			charDisplayW, _ := theCtx.ggCtx.MeasureString(enteredTxt)
	// 			caretX = int(charDisplayW) + Margin
	// 		}

	// 	} else {
	// 		// enter char at subtext
	// 		subTextLen := len(caretXAtSubText)
	// 		tmp := lastLine[:subTextLen] + string(char) + lastLine[subTextLen:]
	// 		tmpW, _ := theCtx.ggCtx.MeasureString(tmp)
	// 		if int(tmpW) > wWidth-(2*Margin) {
	// 			tmpText := lastLine[:subTextLen]
	// 			tmpTextW, _ := theCtx.ggCtx.MeasureString(tmpText)
	// 			if int(tmpTextW) > wWidth-(2*Margin) {
	// 				enteredTxtParts[len(enteredTxtParts)-1] = tmpText + "\n" + string(char) + lastLine[subTextLen:]
	// 				enteredTxt = strings.Join(enteredTxtParts, "\n")
	// 				caretY += FontSize + LineSpacing
	// 				caretX = Margin
	// 			} else {
	// 				enteredTxtParts[len(enteredTxtParts)-1] = tmpText + string(char) + "\n" + lastLine[subTextLen:]
	// 				enteredTxt = strings.Join(enteredTxtParts, "\n")
	// 				caretX = Margin + int(tmpTextW)
	// 			}
	// 		} else {
	// 			enteredTxtParts[len(enteredTxtParts)-1] = tmp
	// 			enteredTxt = strings.Join(enteredTxtParts, "\n")
	// 			charDisplayW, _ := theCtx.ggCtx.MeasureString(lastLine[:subTextLen] + string(char))
	// 			caretX = int(charDisplayW) + Margin
	// 		}
	// 	}

	// } else {

	// 	caretXAtSubText := getCaretXAtSubText(theCtx, caretYLineNo)
	// 	caretXAtSubTextLen := len(caretXAtSubText)
	// 	caretYSubText := enteredTxtParts[caretYLineNo]
	// 	// caretXSubTextW, _ := theCtx.ggCtx.MeasureString(caretXAtSubText)

	// 	tmp := caretYSubText[:caretXAtSubTextLen] + string(char) + caretYSubText[caretXAtSubTextLen:]
	// 	tmpW, _ := theCtx.ggCtx.MeasureString(tmp)
	// 	if int(tmpW) > wWidth-(2*Margin) {
	// 		subTextThen := caretYSubText + "\n" + string(char)
	// 		enteredTxtParts[caretYLineNo] = subTextThen
	// 		enteredTxt = strings.Join(enteredTxtParts, "\n")
	// 		caretY += FontSize + LineSpacing
	// 		caretX = Margin
	// 	} else {
	// 		enteredTxtParts[caretYLineNo] = tmp
	// 		enteredTxt = strings.Join(enteredTxtParts, "\n")
	// 		charDisplayW, _ := theCtx.ggCtx.MeasureString(caretYSubText[:caretXAtSubTextLen] + string(char))
	// 		caretX = int(charDisplayW) + Margin
	// 	}

	// }

	sIRect := objCoords[MajorTextInput]
	theCtx.drawTextInput(MajorTextInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func mKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	wWidth, wHeight := window.GetSize()
	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

	val := enteredTxt
	if key == glfw.KeyBackspace && len(enteredTxt) != 0 {
		enteredTxt = val[:len(val)-1]
		lastChar := string(val[len(val)-1])
		if lastChar == "\n" {
			caretY -= FontSize + LineSpacing
			enteredTxtParts := strings.Split(enteredTxt, "\n")
			textDisplayW, _ := theCtx.ggCtx.MeasureString(enteredTxtParts[len(enteredTxtParts)-1])
			caretX = int(textDisplayW) + Margin
		} else {

			tmpParts := strings.Split(enteredTxt, "\n")
			var lastLine string
			if len(tmpParts) == 1 {
				lastLine = enteredTxt
			} else {
				lastLine = tmpParts[len(tmpParts)-1]
			}

			lastLineW, _ := theCtx.ggCtx.MeasureString(lastLine)
			caretX = Margin + int(lastLineW)
		}
	}

	if key == glfw.KeyEnter {
		enteredTxt = val + "\n"
		caretY += FontSize + LineSpacing
		caretX = Margin
	}

	sIRect := objCoords[MajorTextInput]
	theCtx.drawTextInput(MajorTextInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}
