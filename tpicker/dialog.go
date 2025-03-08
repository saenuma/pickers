package main

import (
	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
)

func drawDialog(window *glfw.Window, suggestions []string) {
	scObjCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	// background image
	img := imaging.AdjustBrightness(currentWindowFrame, -40)
	theCtx := Continue2dCtx(img, &scObjCoords)

	// dialog rectangle
	dialogWidth := wWidth - 50
	dialogHeight := wHeight - 50

	dialogOriginX := (wWidth - dialogWidth) / 2
	dialogOriginY := (wHeight - dialogHeight) / 2

	theCtx.ggCtx.SetHexColor("#fff")
	theCtx.ggCtx.DrawRoundedRectangle(float64(dialogOriginX), float64(dialogOriginY), float64(dialogWidth),
		float64(dialogHeight), 10)
	theCtx.ggCtx.Fill()

	closeBtnX := dialogOriginX + dialogWidth - 100
	theCtx.drawButtonA(SWD_CloseBtn, closeBtnX, dialogOriginY+10, "Close", "#fff", "#B75F5F")

	currentX := dialogOriginX + 20
	currentY := dialogOriginY + 50

	for i, suggestedWord := range suggestions {
		btnId := 1000 + (i + 1)
		sWRect := theCtx.drawButtonA(btnId, currentX, currentY, suggestedWord, "#444", "#fff")

		newX := currentX + sWRect.Width + 10
		maxWidth := dialogOriginX + dialogWidth
		if newX > (maxWidth - sWRect.Width) {
			currentY += 10 + FontSize

			if (currentY + FontSize) > (dialogOriginY + dialogHeight) {
				break
			}

			currentX = 40
		} else {
			currentX += sWRect.Width + 10
		}

	}

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// // save the frame
	// currentWindowFrame = theCtx.ggCtx.Image()
}
