package main

import (
	"fmt"
	"strings"

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
		currentY := Margin
		for i := range len(txtParts) {
			if yPosInt >= (currentY+(i*FontSize)) && yPosInt <= (currentY+((i+1)*FontSize)) {
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
		fmt.Println("Right click")
	}

}
