package main

import (
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func mCharCallback(window *glfw.Window, char rune) {
	wWidth, wHeight := window.GetSize()
	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)

	tmp := enteredTxt + string(char)
	tmpParts := strings.Split(tmp, "\n")
	var lastLine string
	if len(tmpParts) == 1 {
		lastLine = tmp
	} else {
		lastLine = tmpParts[len(tmpParts)-1]
	}
	if len(tmpParts[len(tmpParts)-1]) > 45 {
		enteredTxt = enteredTxt + "\n" + string(char)
		caretY += FontSize + LineSpacing
		caretX = 10
	} else {
		enteredTxt = tmp
		charDisplayW, _ := theCtx.ggCtx.MeasureString(lastLine)
		caretX = int(charDisplayW) + 10
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
			caretX = int(textDisplayW) + 10
		} else {

			tmpParts := strings.Split(enteredTxt, "\n")
			var lastLine string
			if len(tmpParts) == 1 {
				lastLine = enteredTxt
			} else {
				lastLine = tmpParts[len(tmpParts)-1]
			}

			charDisplayW, _ := theCtx.ggCtx.MeasureString(lastLine)
			caretX = 10 + int(charDisplayW)
		}
	}

	if key == glfw.KeyEnter {
		enteredTxt = val + "\n"
		caretY += FontSize + LineSpacing
		caretX = 10
	}

	sIRect := objCoords[MajorTextInput]
	theCtx.drawTextInput(MajorTextInput, sIRect.OriginX, sIRect.OriginY, sIRect.Width, sIRect.Height, enteredTxt)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}
