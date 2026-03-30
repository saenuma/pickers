package main

import (
	// "fmt"

	"slices"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func getCaretYAtLineNumber() int {
	enteredTxtParts := strings.Split(enteredTxt, "\n")
	for i := range len(enteredTxtParts) {
		if caretY >= (i*(FontSize+Margin)) && caretY <= ((i+1)*(FontSize+Margin)) {
			return i
		}
	}
	return len(enteredTxtParts) - 1
}

func getCaretXAtSubText(theCtx Ctx, lineNo int) string {
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

func findLastSpace(line string) int {
	for i := len(line) - 1; i >= 0; i -= 1 {
		if string(line[i]) == " " {
			return i + 1
		}
	}
	return -1
}

func lineBreak(theCtx Ctx, tmpText string, maxWidth int) []string {
	var tmpLine string
	lines := make([]string, 0)
	for i, aChar := range tmpText {
		aCharStr := string(aChar)
		if aCharStr == "\n" {
			lines = append(lines, tmpLine)
			tmpLine = ""
			continue
		}

		tmpLine += aCharStr
		tmpLineW, _ := theCtx.ggCtx.MeasureString(tmpLine)
		if int(tmpLineW) > maxWidth {
			lastIndexTmpLine := len(tmpLine) - 1
			if string(tmpText[i]) == " " {
				lines = append(lines, tmpLine[:lastIndexTmpLine])
				tmpLine = ""
			} else {
				lastSpaceIndex := strings.LastIndex(tmpLine, " ")
				lines = append(lines, tmpLine[:lastSpaceIndex+1])
				tmpLine = tmpLine[lastSpaceIndex+1:]
			}

		}
	}
	lines = append(lines, tmpLine)
	return lines
}

func mCharCallback(window *glfw.Window, char rune) {
	capitalizeSentences := func(text string) string {
		textParts := strings.Split(text, ". ")
		for i, part := range textParts {
			if len(part) <= 1 {
				continue
			}
			textParts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
		return strings.Join(textParts, ". ")
	}

	wWidth, wHeight := window.GetSize()
	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

	maxWidth := wWidth - (2 * Margin)

	if len(enteredTxt) > 0 {
		enteredTxt = capitalizeSentences(enteredTxt)
	}
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
			if string(char) == " " {
				currentLine = currentLine + "\n" + string(char)
				charW, _ := theCtx.ggCtx.MeasureString(string(char))
				caretX = Margin + int(charW)
			} else {
				// currentLine = currentLine[:len(currentLine)-1] + "\n" + string(currentLine[len(currentLine)-1]) + string(char)
				lastSpaceIndex := findLastSpace(currentLine)
				currentLine = currentLine[:lastSpaceIndex] + "\n" + currentLine[lastSpaceIndex:] + string(char)
				charsW, _ := theCtx.ggCtx.MeasureString(currentLine[lastSpaceIndex:] + string(char))
				caretX = Margin + int(charsW)
			}
			enteredTxtParts[caretYLineNo] = currentLine
			enteredTxt = strings.Join(enteredTxtParts, "\n")
			caretY += FontSize + LineSpacing
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
				if string(char) == " " {
					enteredLine := currentLine[:subTextLen] + "\n" + string(char) + currentLine[subTextLen:]
					enteredTxtParts[caretYLineNo] = enteredLine
					enteredTxt = strings.Join(enteredTxtParts, "\n")
					charW, _ := theCtx.ggCtx.MeasureString(string(char))
					caretX = Margin + int(charW)
				} else {
					lastSpaceIndex := findLastSpace(currentLine)
					enteredLine := currentLine[:lastSpaceIndex] + "\n" + currentLine[lastSpaceIndex:] + string(char)
					enteredTxtParts[caretYLineNo] = enteredLine
					enteredTxt = strings.Join(enteredTxtParts, "\n")
					charsW, _ := theCtx.ggCtx.MeasureString(currentLine[lastSpaceIndex:] + string(char))
					caretX = Margin + int(charsW)
				}
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

	sIRect := objCoords[MajorTextInput]
	theCtx.drawTextInput(MajorTextInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func mKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press || action == glfw.Repeat {
		wWidth, wHeight := window.GetSize()
		theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

		maxWidth := wWidth - (2 * Margin)

		enteredTxtParts := strings.Split(enteredTxt, "\n")
		caretYLineNo := getCaretYAtLineNumber()
		currentLine := enteredTxtParts[caretYLineNo]
		currentLineW, _ := theCtx.ggCtx.MeasureString(currentLine)
		caretXAtSubText := getCaretXAtSubText(theCtx, caretYLineNo)
		caretXAtSubTextW, _ := theCtx.ggCtx.MeasureString(caretXAtSubText)

		if key == glfw.KeyLeft {
			if caretX != Margin {
				tmp := caretXAtSubText[:len(caretXAtSubText)-1]
				tmpW, _ := theCtx.ggCtx.MeasureString(tmp)
				caretX = Margin + int(tmpW)
			}
		}
		if key == glfw.KeyRight {
			if currentLine != caretXAtSubText {
				tmp := caretXAtSubText + currentLine[len(caretXAtSubText):len(caretXAtSubText)+1]
				tmpW, _ := theCtx.ggCtx.MeasureString(tmp)
				caretX = Margin + int(tmpW)
			}
		}
		if key == glfw.KeyUp {
			if caretYLineNo != 0 {
				caretY -= FontSize + LineSpacing
				possibleNewSubtext := enteredTxtParts[caretYLineNo-1][:len(caretXAtSubText)]
				possibleNewSubtextW, _ := theCtx.ggCtx.MeasureString(possibleNewSubtext)
				caretX = Margin + int(possibleNewSubtextW)
			}
		}
		if key == glfw.KeyDown {
			if caretYLineNo != len(enteredTxtParts)-1 {
				caretY += FontSize + LineSpacing
				possibleNewSubtext := enteredTxtParts[caretYLineNo+1][:len(caretXAtSubText)]
				possibleNewSubtextW, _ := theCtx.ggCtx.MeasureString(possibleNewSubtext)
				caretX = Margin + int(possibleNewSubtextW)
			}
		}

		if caretYLineNo == len(enteredTxtParts)-1 && int(caretXAtSubTextW) == int(currentLineW) {
			// the caret is at the end of the text

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

			// implement paste
			if window.GetKey(glfw.KeyLeftControl) == glfw.Press && key == glfw.KeyV {
				copiedText := glfw.GetClipboardString()
				if len(copiedText) > 0 {
					tmpText := enteredTxt + copiedText
					tmpTextBroken := lineBreak(theCtx, tmpText, maxWidth)
					enteredTxt = strings.Join(tmpTextBroken, "\n")
					lastLine := tmpTextBroken[len(tmpTextBroken)-1]
					lastLineW, _ := theCtx.ggCtx.MeasureString(lastLine)
					caretX = Margin + int(lastLineW)
					caretY = Margin + ((len(tmpTextBroken) - 1) * (FontSize + LineSpacing))
					if caretX == maxWidth {
						caretX = Margin
						caretY += FontSize + LineSpacing
					}
				}

			}

		} else {
			// the caret is not at the end of the input.

			subTextLen := len(caretXAtSubText)

			if key == glfw.KeyBackspace {
				if caretX == Margin && caretY == Margin {
					return
				}
				if len(caretXAtSubText) == 0 {
					lastLine := enteredTxtParts[caretYLineNo-1]
					editedLastLine := lastLine + currentLine
					editedLastLineW, _ := theCtx.ggCtx.MeasureString(editedLastLine)
					if int(editedLastLineW) > maxWidth {
						tmpLastLine := lastLine[:len(lastLine)-1]
						enteredTxtParts[caretYLineNo-1] = tmpLastLine
						enteredTxt = strings.Join(enteredTxtParts, "\n")
						caretX = maxWidth
					} else {
						enteredTxtParts[caretYLineNo-1] = editedLastLine
						enteredTxtParts = slices.Delete(enteredTxtParts, caretYLineNo, caretYLineNo+1)
						enteredTxt = strings.Join(enteredTxtParts, "\n")
						lastLineW, _ := theCtx.ggCtx.MeasureString(lastLine)
						caretX = Margin + int(lastLineW)
					}

					caretY -= FontSize + LineSpacing

				} else {
					editedSubText := caretXAtSubText[:len(caretXAtSubText)-1] + currentLine[subTextLen:]
					allFromEdit := editedSubText + " " + strings.Join(enteredTxtParts[caretYLineNo+1:], " ")
					allFromEditWrapped := theCtx.ggCtx.WordWrap(allFromEdit, float64(WindowWidth))
					enteredTxtParts = append(enteredTxtParts[:caretYLineNo], allFromEditWrapped...)
					// enteredTxtParts[caretYLineNo] = editedSubText
					enteredTxt = strings.Join(enteredTxtParts, "\n")
					subTextW, _ := theCtx.ggCtx.MeasureString(caretXAtSubText[:len(caretXAtSubText)-1])
					caretX = Margin + int(subTextW)
				}
			}

			if key == glfw.KeyEnter {
				enteredTxtParts[caretYLineNo] = caretXAtSubText
				enteredTxtParts = slices.Insert(enteredTxtParts, caretYLineNo+1, currentLine[subTextLen:])
				enteredTxt = strings.Join(enteredTxtParts, "\n")
				caretY += FontSize + LineSpacing
				caretX = Margin
			}

			// implement paste
			if window.GetKey(glfw.KeyLeftControl) == glfw.Press && key == glfw.KeyV {
				copiedText := glfw.GetClipboardString()
				if len(copiedText) > 0 {

					enteredLine := currentLine[:subTextLen] + copiedText + currentLine[subTextLen:]
					enteredTxtParts[caretYLineNo] = enteredLine

					nEnteredTxt := strings.Join(enteredTxtParts, "\n")
					tmpTextBroken := lineBreak(theCtx, nEnteredTxt, maxWidth)
					enteredTxt = strings.Join(tmpTextBroken, "\n")
					lastLine := tmpTextBroken[len(tmpTextBroken)-1]
					lastLineW, _ := theCtx.ggCtx.MeasureString(lastLine)
					caretX = Margin + int(lastLineW)
					caretY = Margin + ((len(tmpTextBroken) - 1) * (FontSize + LineSpacing))
					if caretX == maxWidth {
						caretX = Margin
						caretY += FontSize + LineSpacing
					}
				}

			}

		}

		sIRect := objCoords[MajorTextInput]
		theCtx.drawTextInput(MajorTextInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		currentWindowFrame = theCtx.ggCtx.Image()

	}
}
