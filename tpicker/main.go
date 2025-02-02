package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/gookit/color"
)

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

	window := g143.NewWindow(900, 500, "sae.ng text picker", false)
	allDraws(window)
	if changeLoc {
		window.SetPos(windowX, windowY)
	}

	// window.SetMouseButtonCallback(mouseBtnCallback)
	// window.SetCursorPosCallback(cursorCallback)
	window.SetKeyCallback(mKeyCallback)
	window.SetCharCallback(mCharCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}
}

func allDraws(window *glfw.Window) {
	objCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight, &objCoords)

	// save button
	theCtx.drawButtonA(DoneBtn, wWidth-100, wHeight-50, "Done", "#fff", "#5c8f8c")
	theCtx.drawTextInput(MajorTextInput, 10, 10, wWidth-20, wHeight-70, "")

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}
