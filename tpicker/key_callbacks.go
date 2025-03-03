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
	lastLineW, _ := theCtx.ggCtx.MeasureString(lastLine)
	if int(lastLineW) > wWidth-(2*Margin) {
		enteredTxt = enteredTxt + "\n" + string(char)
		caretY += FontSize + LineSpacing
		caretX = Margin
	} else {
		enteredTxt = tmp
		charDisplayW, _ := theCtx.ggCtx.MeasureString(lastLine)
		caretX = int(charDisplayW) + Margin
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
