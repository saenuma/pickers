package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	fps      = 10
	fontSize = 20
)

var objCoords map[int]g143.Rect

var allColors []string

func main() {
	runtime.LockOSThread()

	objCoords = make(map[int]g143.Rect)

	colors := strings.Split(string(Colors2), "\n")

	allColors = colors
	window := g143.NewWindow(1350, 700, "sae.ng color picker", false)
	allDraws(window, colors)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func allDraws(window *glfw.Window, colors []string) {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	gutter := 10
	currentX := gutter
	currentY := gutter

	boxDimension := 55
	for i, aColor := range colors {
		ggCtx.SetHexColor(aColor)
		ggCtx.DrawRectangle(float64(currentX), float64(currentY), float64(boxDimension), float64(boxDimension))
		ggCtx.Fill()
		aColorRS := g143.Rect{OriginX: currentX, OriginY: currentY, Width: boxDimension, Height: boxDimension}
		objCoords[i+1] = aColorRS

		newX := currentX + boxDimension + gutter
		if newX > (wWidth - boxDimension) {
			currentY += boxDimension + gutter
			currentX = gutter
		} else {
			currentX += boxDimension + gutter
		}

	}

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}

func getDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "fpicker_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	fmt.Println(allColors[widgetCode-1])
	window.SetShouldClose(true)
}
