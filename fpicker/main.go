package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	fps      = 10
	fontSize = 20

	BackBtn = 9801
)

var objCoords map[int]g143.Rect

var rootPath string
var basePath string
var exts string

var tmpPickerFrame image.Image

func main() {
	if len(os.Args) != 3 {
		panic("expecting a rootpath and a extension")
	}
	rootPath = os.Args[1]
	exts = os.Args[2]
	basePath = rootPath

	runtime.LockOSThread()

	objCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(1200, 800, "sae.ng file picker", false)
	allDraws(window)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetCursorPosCallback(cursorCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func getObjects(rootPath, mergedExts string) []string {
	dirEs, err := os.ReadDir(rootPath)
	if err != nil {
		log.Println(err)
		return []string{}
	}

	extsParts := strings.Split(mergedExts, "|")

	allFolders := make([]string, 0)
	allFiles := make([]string, 0)
	for _, dirE := range dirEs {
		if !strings.HasPrefix(dirE.Name(), ".") && dirE.IsDir() {
			allFolders = append(allFolders, dirE.Name()+"/")
			continue
		}

		for _, ext := range extsParts {
			if !strings.HasPrefix(dirE.Name(), ".") && strings.HasSuffix(dirE.Name(), ext) {
				allFiles = append(allFiles, dirE.Name())
			}
		}

	}

	sort.Strings(allFolders)
	sort.Strings(allFiles)
	return append(allFolders, allFiles...)
}

func shortenObject(filename string) string {
	var tmp string
	if len(filename) > 30 {
		tmp = filename[0:20] + "..." + filename[len(filename)-6:]
	} else {
		tmp = filename
	}

	return tmp
}

func allDraws(window *glfw.Window) {
	toPickFrom := getObjects(basePath, exts)

	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	ggCtx.LoadFontFace(fontPath, 20)

	ggCtx.SetHexColor("#444")
	ggCtx.DrawRectangle(20, 5, 30, 30)
	ggCtx.Fill()
	objCoords[BackBtn] = g143.NewRect(20, 5, 30, 30)

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString("<", 30, 5+fontSize)

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(basePath, 60, 5+20)

	// draw divider
	ggCtx.SetHexColor("#bbb")
	ggCtx.DrawRectangle(5, 40, float64(wWidth)-10, 2)
	ggCtx.Fill()

	currentX := 20
	currentY := 50

	for i, aFile := range toPickFrom {
		shortAFile := shortenObject(aFile)
		aFileStrW, _ := ggCtx.MeasureString(shortAFile)

		if strings.HasSuffix(aFile, "/") {
			ggCtx.SetHexColor("#9e9770")
			ggCtx.DrawString(shortAFile, float64(currentX), float64(currentY)+fontSize)
		} else {
			ggCtx.SetHexColor("#444")
			ggCtx.DrawString(shortAFile, float64(currentX), float64(currentY)+fontSize)
		}

		aFileRS := g143.Rect{OriginX: currentX, OriginY: currentY, Width: int(aFileStrW), Height: fontSize}
		objCoords[i+1] = aFileRS

		newX := currentX + int(aFileStrW) + 20
		if newX > (wWidth - int(aFileStrW)) {
			currentY += 40
			currentX = 20
		} else {
			currentX += int(aFileStrW) + 20 + 10
		}
	}

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	tmpPickerFrame = ggCtx.Image()
}

func getDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "font.ttf")
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

	if widgetCode == BackBtn {
		tmp := filepath.Dir(basePath)
		var rootPathTmp string
		if strings.HasSuffix(rootPath, "/") {
			rootPathTmp = rootPath[:len(rootPath)-1]
		}

		if strings.Count(rootPathTmp, "/") < strings.Count(tmp, "/")+1 {
			objCoords = make(map[int]g143.Rect)
			basePath = tmp
			allDraws(window)
		}

		return
	}

	toPickFrom := getObjects(basePath, exts)
	foundObject := toPickFrom[widgetCode-1]

	if strings.HasSuffix(foundObject, "/") {
		objCoords = make(map[int]g143.Rect)
		basePath = filepath.Join(basePath, foundObject)
		allDraws(window)
	} else {
		fmt.Println(filepath.Join(basePath, foundObject))
		window.SetShouldClose(true)

	}
}

var cursorEventsCount = 0

func cursorCallback(window *glfw.Window, xpos, ypos float64) {
	if runtime.GOOS == "linux" {
		// linux fires too many events
		cursorEventsCount += 1
		if cursorEventsCount != 10 {
			return
		} else {
			cursorEventsCount = 0
		}
	}

	wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	xPosInt := int(xpos)
	yPosInt := int(ypos)
	for code, RS := range objCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		allDraws(window)
		return
	}
	if widgetCode == BackBtn {
		return
	}

	toPickFrom := getObjects(basePath, exts)
	foundObject := toPickFrom[widgetCode-1]

	if strings.HasSuffix(foundObject, "/") {
		allDraws(window)
		return
	}

	foundObjectPath := filepath.Join(basePath, foundObject)

	isPicFormat := false
	picFormats := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, format := range picFormats {
		if strings.HasSuffix(foundObject, format) {
			isPicFormat = true
		}
	}
	if !isPicFormat {
		allDraws(window)
		return
	}

	foundImg, err := imaging.Open(foundObjectPath)
	if err != nil {
		allDraws(window)
		return
	}
	previewImgBoxSize := 300

	foundImg = imaging.Fit(foundImg, previewImgBoxSize, previewImgBoxSize, imaging.Lanczos)

	objCoords = make(map[int]g143.Rect)
	allDraws(window)

	previewX := xPosInt + 10
	if previewX+previewImgBoxSize > wWidth {
		previewX = xPosInt - previewImgBoxSize - 10
	}

	previewY := yPosInt + 10
	if previewY+previewImgBoxSize > wHeight {
		previewY = yPosInt - previewImgBoxSize - 10
	}

	ggCtx := gg.NewContextForImage(tmpPickerFrame)
	ggCtx.DrawImage(foundImg, previewX, previewY)
	ggCtx.Fill()

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

}
