package main

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gookit/color"
)

const (
	fps      = 10
	fontSize = 20

	AllColorBox = 31
	AColorBox   = 32
	SelectBtn   = 33
)

var objCoords map[int]g143.Rect
var tmpFrame image.Image
var pickedColor string
var currentHue int
var CursorEventsCount int

func main() {
	changeLoc := false
	var windowX, windowY int
	if len(os.Args) == 3 {
		changeLoc = true
		xCoord, err := strconv.Atoi(os.Args[1])
		if err != nil {
			color.Red.Println("First argument X was not int")
			os.Exit(1)
		}

		yCoord, err := strconv.Atoi(os.Args[2])
		if err != nil {
			color.Red.Println("Second argument Y was not int")
			os.Exit(1)
		}

		windowX, windowY = xCoord, yCoord
	}

	runtime.LockOSThread()

	objCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(450, 430, "sae.ng all colors picker", false)
	allDraws(window, 0)
	if changeLoc {
		window.SetPos(windowX, windowY)
	}

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetCursorPosCallback(CursorPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func allDraws(window *glfw.Window, hue int) {
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

	// draw color picker
	acImg := allColorImg()
	ggCtx.DrawImage(acImg, 10, 10)
	objCoords[AllColorBox] = g143.NewRect(10, 10, acImg.Rect.Dx(),
		acImg.Rect.Dy())

	acImg2 := subsetColorImg(hue)
	ggCtx.DrawImage(acImg2, 70, 10)
	objCoords[AColorBox] = g143.NewRect(70, 10, acImg2.Rect.Dx(),
		acImg2.Rect.Dy())

	// draw picked button
	pblW, _ := ggCtx.MeasureString("select")
	ggCtx.SetHexColor("#666")
	ggCtx.DrawRectangle(350, 380, pblW+20, 40)
	ggCtx.Fill()

	objCoords[SelectBtn] = g143.NewRect(350, 380, int(pblW)+20, 40)

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString("select", 360, 385+fontSize)

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	tmpFrame = ggCtx.Image()
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

	wWidth, wHeight := window.GetSize()

	var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	switch widgetCode {
	case AllColorBox:
		ptClicked := yPosInt - widgetRS.OriginY
		currentHue = ptClicked
		allDraws(window, ptClicked)

	case AColorBox:
		xPtClicked := xPosInt - widgetRS.OriginX
		yPtClicked := yPosInt - widgetRS.OriginY

		acImg2 := subsetColorImg(currentHue)
		r, g, b, _ := acImg2.At(xPtClicked, yPtClicked).RGBA()
		hexColor := fmt.Sprintf("#%02x%02x%02x", r/255, g/255, b/255)
		pickedColor = hexColor

		ggCtx := gg.NewContextForImage(tmpFrame)

		ggCtx.SetHexColor(hexColor)
		ggCtx.DrawRectangle(100, 20+360, 200, 40)
		ggCtx.Fill()

		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

		tmpFrame = ggCtx.Image()

	case SelectBtn:
		if pickedColor != "" {
			fmt.Println(pickedColor)
			window.SetShouldClose(true)
		}

	}
}

func CursorPosCB(window *glfw.Window, xpos, ypos float64) {
	if runtime.GOOS == "linux" {
		// linux fires too many events
		CursorEventsCount += 1
		if CursorEventsCount != 10 {
			return
		} else {
			CursorEventsCount = 0
		}
	}

	wWidth, wHeight := window.GetSize()

	var widgetRS g143.Rect
	var widgetCode int

	xPosInt := int(xpos)
	yPosInt := int(ypos)
	for code, RS := range objCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == SelectBtn {
		rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
			widgetRS.OriginX+widgetRS.Width,
			widgetRS.OriginY+widgetRS.Height)

		pieceOfCurrentFrame := imaging.Crop(tmpFrame, rectA)
		invertedPiece := imaging.AdjustBrightness(pieceOfCurrentFrame, -20)

		ggCtx := gg.NewContextForImage(tmpFrame)
		ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

		// send the frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

	} else {
		// send the last drawn frame to glfw window
		windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
		g143.DrawImage(wWidth, wHeight, tmpFrame, windowRS)
		window.SwapBuffers()
		return
	}

}
