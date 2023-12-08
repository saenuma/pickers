package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	fps      = 10
	fontSize = 20
)

var objCoords map[int]g143.RectSpecs

var allColors []string

func main() {
	runtime.LockOSThread()

	objCoords = make(map[int]g143.RectSpecs)

	colors := make([]map[string]string, 0)
	json.Unmarshal(ColorJson, &colors)
	hexColors := make([]string, 0)
	for _, obj := range colors {
		hexColors = append(hexColors, obj["hex"])
	}
	sort.Strings(hexColors)
	allColors = hexColors
	window := g143.NewWindow(1300, 800, "sae.ng color picker", false)
	allDraws(window, hexColors)

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

	currentX := 5
	currentY := 5

	boxDimension := 30
	for i, aColor := range colors {
		ggCtx.SetHexColor(aColor)
		ggCtx.DrawRectangle(float64(currentX), float64(currentY), float64(boxDimension), float64(boxDimension))
		ggCtx.Fill()
		aColorRS := g143.RectSpecs{OriginX: currentX, OriginY: currentY, Width: boxDimension, Height: boxDimension}
		objCoords[i+1] = aColorRS

		newX := currentX + boxDimension + 5
		if newX > (wWidth - boxDimension) {
			currentY += boxDimension + 5
			currentX = 5
		} else {
			currentX += boxDimension + 5
		}

	}

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
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

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
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
