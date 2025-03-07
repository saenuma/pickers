package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
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

	// prepare words.txt for search
	words := strings.Split(string(AllWords), "\n")
	wordsByFirstCharMap = make(map[string][]string)
	for _, word := range words {
		arr, ok := wordsByFirstCharMap[string(word[0])]
		if ok {
			wordsByFirstCharMap[string(word[0])] = append(arr, word)
		} else {
			wordsByFirstCharMap[string(word[0])] = []string{word}
		}
	}

	spellcheckTrie = NewSpellcheckTrie()
	runtime.LockOSThread()

	objCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(900, 300, "sae.ng text picker", false)
	allDraws(window)
	if changeLoc {
		window.SetPos(windowX, windowY)
	}

	window.SetMouseButtonCallback(mouseBtnCallback)
	// window.SetCursorPosCallback(cursorCallback)
	window.SetKeyCallback(mKeyCallback)
	window.SetCharCallback(mCharCallback)
	window.SetCloseCallback(func(w *glfw.Window) {
		fmt.Println(enteredTxt)
	})

	go func() {
		for {
			time.Sleep(time.Second)
			windowFrameWithErrors = getDisplayWithErrors()
			frameUpdated = true
		}
	}()

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		displayCaret(window)

		if frameUpdated {
			wWidth, wHeight := window.GetSize()
			theCtx := Continue2dCtx(currentWindowFrame, &objCoords)
			// send the frame to glfw window
			g143.DrawImage(wWidth, wHeight, windowFrameWithErrors, theCtx.windowRect())
			window.SwapBuffers()

			currentWindowFrame = windowFrameWithErrors
			frameUpdated = false
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}
}

func allDraws(window *glfw.Window) {
	objCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight, &objCoords)
	theCtx.drawTextInput(MajorTextInput, 10, 10, wWidth-20, wHeight-20, "")

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func displayCaret(window *glfw.Window) {
	if caretDisplayCount != MaxCaretDisplayCount {
		caretDisplayCount += 1
		return
	}
	caretDisplayCount = 0

	wWidth, wHeight := window.GetSize()

	theCtx := Continue2dCtx(currentWindowFrame, &objCoords)
	if caretDisplayed {
		caretDisplayed = false

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, currentWindowFrame, theCtx.windowRect())
		window.SwapBuffers()
	} else {
		theCtx.ggCtx.SetHexColor("#444")
		caretDisplayed = true

		theCtx.ggCtx.DrawRectangle(float64(caretX), float64(caretY)-5, 2, FontSize+10)
		theCtx.ggCtx.Fill()

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()
	}

}
