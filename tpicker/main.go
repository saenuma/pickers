package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/essentialkaos/ek/v13/spellcheck"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	startingText := ""
	if len(os.Args) == 2 {
		startingText = os.Args[1]
	}

	// prepare words.txt for search
	allWords := strings.ReplaceAll(string(AllWords), "\r", "")
	words := strings.Split(allWords, "\n")
	wordsByFirstCharMap = make(map[string][]string)
	for _, word := range words {
		arr, ok := wordsByFirstCharMap[string(word[0])]
		if ok {
			wordsByFirstCharMap[string(word[0])] = append(arr, word)
		} else {
			wordsByFirstCharMap[string(word[0])] = []string{word}
		}
	}

	spellcheckModel = spellcheck.Train(words)
	runtime.LockOSThread()

	objCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(700, 300, "sae.ng text picker", false)
	if len(os.Args) == 2 {
		enteredTxt = startingText
	}
	drawTextView(window)

	window.SetMouseButtonCallback(mouseBtnCallback)
	// window.SetCursorPosCallback(cursorCallback)
	window.SetKeyCallback(mKeyCallback)
	window.SetCharCallback(mCharCallback)
	window.SetCloseCallback(func(w *glfw.Window) {
		fmt.Println(enteredTxt)
	})

	go func() {
		for {
			if suggestionsDialogShown {
				continue
			}

			time.Sleep(time.Second)
			windowFrameWithErrors = getDisplayWithErrors()
			frameUpdated = true
		}
	}()

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		displayCaret(window)

		if !suggestionsDialogShown {
			if frameUpdated {
				wWidth, wHeight := window.GetSize()
				theCtx := Continue2dCtx(currentWindowFrame, &objCoords)
				// send the frame to glfw window
				g143.DrawImage(wWidth, wHeight, windowFrameWithErrors, theCtx.windowRect())
				window.SwapBuffers()

				currentWindowFrame = windowFrameWithErrors
				frameUpdated = false
			}
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}
}

func drawTextView(window *glfw.Window) {
	objCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight, &objCoords)
	theCtx.drawTextInput(MajorTextInput, 10, 10, wWidth-20, wHeight-20, enteredTxt)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = theCtx.ggCtx.Image()
}

func displayCaret(window *glfw.Window) {
	if suggestionsDialogShown {
		return
	}
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
