package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/pickers/internal"
)

func main() {
	if len(os.Args) != 3 {
		panic("expecting a rootpath and a extension")
	}
	rootPath = os.Args[1]
	exts = os.Args[2]
	basePath = rootPath

	runtime.LockOSThread()

	FontSize = internal.GetFontSize()
	windowW, windowH := getWindowSize()
	window := g143.NewWindow(int(windowW), int(windowH), "sae.ng file picker", false)
	drawObjects(window, 1)

	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetCursorPosCallback(cursorCallback)
	window.SetScrollCallback(scrollBtnCB)
	window.SetFramebufferSizeCallback(frameBufferSizeCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		tmpFontSize := internal.GetFontSize()
		if internal.NotEqual(FontSize, tmpFontSize) {
			windowW, windowH := getWindowSize()
			window.SetSize(int(windowW), int(windowH))
			FontSize = tmpFontSize
		}

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func frameBufferSizeCallback(window *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	drawObjects(window, CurrentPage)
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetCursorPosCallback(cursorCallback)
	window.SetScrollCallback(scrollBtnCB)
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

func getWindowSize() (float64, float64) {
	textScale := internal.GetTextScale()
	return DefaultWindowW * textScale, DefaultWindowH * textScale
}

func shortenObject(filename string) string {
	var tmp string
	if len(filename) > 20 {
		tmp = filename[0:10] + "..." + filename[len(filename)-6:]
	} else {
		tmp = filename
	}

	strings.TrimPrefix(tmp, "/")

	return tmp
}

func TotalPages() int {
	objs := getObjects(basePath, exts)
	return int(math.Ceil(float64(len(objs)) / float64(PageSize)))
}

func GetPageObjects(page int) []string {
	beginIndex := (page - 1) * PageSize
	endIndex := beginIndex + PageSize

	toPickFrom := getObjects(basePath, exts)
	var retObjects []string
	if len(toPickFrom) <= PageSize {
		retObjects = toPickFrom
	} else if page == 1 {
		retObjects = toPickFrom[:PageSize]
	} else if endIndex > len(toPickFrom) {
		retObjects = toPickFrom[beginIndex:]
	} else {
		retObjects = toPickFrom[beginIndex:endIndex]
	}
	return retObjects
}

// func drawObjects(window *glfw.Window, page int) {
// 	CurrentPage = page
// 	wWidth, wHeight := window.GetSize()

// 	// frame buffer
// 	ggCtx := gg.NewContext(wWidth, wHeight)
// 	textScale := internal.GetTextScale()

// 	// background rectangle
// 	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
// 	ggCtx.SetHexColor("#ffffff")
// 	ggCtx.Fill()

// 	// load font
// 	fontPath := internal.GetDefaultFontPath()
// 	ggCtx.LoadFontFace(fontPath, FontSize)

// 	ggCtx.SetHexColor("#444")
// 	ggCtx.DrawRectangle(20, 5, 30*textScale, 30*textScale)
// 	ggCtx.Fill()
// 	ObjCoords[BackBtn] = g143.NewRect(20, 5, 30, 30)

// 	ggCtx.SetHexColor("#fff")
// 	ggCtx.DrawString("<", 30, 5+FontSize)

// 	ggCtx.SetHexColor("#444")
// 	ggCtx.DrawString(basePath, 60, 5+20)

// 	currentX := 20
// 	currentY := 50

// 	toPickFrom := GetPageObjects(page)
// 	for i, aFile := range toPickFrom {
// 		shortAFile := shortenObject(aFile)
// 		// aFileStrW, _ := ggCtx.MeasureString(shortAFile)
// 		objW := 270

// 		if strings.HasSuffix(aFile, "/") {
// 			ggCtx.SetHexColor("#444")
// 			ggCtx.DrawRoundedRectangle(float64(currentX), float64(currentY), float64(objW)+10, FontSize+20, 4)
// 			ggCtx.Fill()
// 			ggCtx.SetHexColor("#fff")
// 			ggCtx.DrawRoundedRectangle(float64(currentX)+1, float64(currentY)+1, float64(objW)+10-2, FontSize+20-2, 2)
// 			ggCtx.Fill()
// 			ggCtx.SetHexColor("#444")
// 			ggCtx.DrawString(shortAFile, float64(currentX+5), float64(currentY)+FontSize+10)
// 		} else {
// 			ggCtx.SetHexColor("#444")
// 			ggCtx.DrawString(shortAFile, float64(currentX+5), float64(currentY)+FontSize+10)
// 		}

// 		aFileRS := g143.Rect{OriginX: currentX, OriginY: currentY, Width: int(objW) + 10, Height: int(FontSize + 20)}
// 		ObjCoords[i+1] = aFileRS

// 		newX := currentX + int(objW) + 10
// 		if newX > (wWidth - int(objW)) {
// 			currentY += int(FontSize) + 20 + 10
// 			currentX = 20
// 		} else {
// 			currentX += int(objW) + 20
// 		}
// 	}

// 	pageLabel := fmt.Sprintf("Page %d / %d", CurrentPage, TotalPages())
// 	pageLabelW, _ := ggCtx.MeasureString(pageLabel)
// 	pageLabelX := (wWidth - int(pageLabelW)) / 2
// 	ggCtx.DrawString(pageLabel, float64(pageLabelX), float64(wHeight)-20)

// 	// send the frame to glfw window
// 	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
// 	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
// 	window.SwapBuffers()

// 	tmpPickerFrame = ggCtx.Image()
// }

func drawObjects(window *glfw.Window, page int) {
	CurrentPage = page
	wWidth, wHeight := window.GetSize()

	ObjCoords = make(map[int]g143.Rect)
	theCtx := New2dCtx(wWidth, wHeight)
	textScale := internal.GetTextScale()

	bBRS := theCtx.drawButtonA(BackBtn, 10, 5, "<", "#fff", "#444")
	bPSX := nextX(bBRS, 10)
	theCtx.ggCtx.SetHexColor("#444")
	theCtx.ggCtx.DrawString(basePath, float64(bPSX), 10+FontSize)

	currentX := 20
	currentY := nextY(bBRS, 10)

	toPickFrom := GetPageObjects(page)
	for i, aFile := range toPickFrom {
		shortAFile := shortenObject(aFile)
		// aFileStrW, _ := ggCtx.MeasureString(shortAFile)
		var oRS g143.Rect
		if strings.HasSuffix(aFile, "/") {
			oRS = theCtx.drawButtonA(i+1, currentX, currentY, shortAFile, "#444", "#ccc")
		} else {
			oRS = theCtx.drawButtonA(i+1, currentX, currentY, shortAFile, "#444", "#fff")
		}

		newX := currentX + int(oRS.Width) + int(10*textScale)
		if newX > (wWidth - int(oRS.Width)) {
			currentY += int(FontSize) + int(30*textScale)
			currentX = 20
		} else {
			currentX += int(oRS.Width) + int(20*textScale)
		}
	}

	pageLabel := fmt.Sprintf("Page %d / %d", CurrentPage, TotalPages())
	pageLabelW, _ := theCtx.ggCtx.MeasureString(pageLabel)
	pageLabelX := (wWidth - int(pageLabelW)) / 2
	theCtx.ggCtx.DrawString(pageLabel, float64(pageLabelX), float64(wHeight)-FontSize)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	tmpPickerFrame = theCtx.ggCtx.Image()
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

	for code, RS := range ObjCoords {
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
		rootPathTmp := rootPath
		if strings.HasSuffix(rootPath, "/") {
			rootPathTmp = rootPath[:len(rootPath)-1]
		}

		if strings.Count(rootPathTmp, "/") < strings.Count(tmp, "/")+1 {
			ObjCoords = make(map[int]g143.Rect)
			basePath = tmp
			drawObjects(window, CurrentPage)
		}

		return
	}

	toPickFrom := GetPageObjects(CurrentPage)
	foundObject := toPickFrom[widgetCode-1]

	if strings.HasSuffix(foundObject, "/") {
		ObjCoords = make(map[int]g143.Rect)
		basePath = filepath.Join(basePath, foundObject)
		drawObjects(window, CurrentPage)
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

	var widgetRS g143.Rect
	var widgetCode int
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}

	xPosInt := int(xpos)
	yPosInt := int(ypos)
	for code, RS := range ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		drawObjects(window, CurrentPage)
		return
	}
	if widgetCode == BackBtn {

		rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
			widgetRS.OriginX+widgetRS.Width,
			widgetRS.OriginY+widgetRS.Height)

		pieceOfCurrentFrame := imaging.Crop(tmpPickerFrame, rectA)
		invertedPiece := imaging.AdjustBrightness(pieceOfCurrentFrame, -20)

		ggCtx := gg.NewContextForImage(tmpPickerFrame)
		ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

		return
	}

	toPickFrom := GetPageObjects(CurrentPage)
	foundObject := toPickFrom[widgetCode-1]

	if strings.HasSuffix(foundObject, "/") {
		rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
			widgetRS.OriginX+widgetRS.Width,
			widgetRS.OriginY+widgetRS.Height)

		pieceOfCurrentFrame := imaging.Crop(tmpPickerFrame, rectA)
		invertedPiece := imaging.AdjustBrightness(pieceOfCurrentFrame, -20)

		ggCtx := gg.NewContextForImage(tmpPickerFrame)
		ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

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
		rectA := image.Rect(widgetRS.OriginX, widgetRS.OriginY,
			widgetRS.OriginX+widgetRS.Width,
			widgetRS.OriginY+widgetRS.Height)

		pieceOfCurrentFrame := imaging.Crop(tmpPickerFrame, rectA)
		invertedPiece := imaging.AdjustBrightness(pieceOfCurrentFrame, -20)

		ggCtx := gg.NewContextForImage(tmpPickerFrame)
		ggCtx.DrawImage(invertedPiece, widgetRS.OriginX, widgetRS.OriginY)

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
		window.SwapBuffers()

		return
	}

	foundImg, err := imaging.Open(foundObjectPath, imaging.AutoOrientation(true))
	if err != nil {
		drawObjects(window, CurrentPage)
		return
	}
	textScale := internal.GetTextScale()
	previewImgBoxSize := int(150 * textScale)

	foundImg = imaging.Fit(foundImg, previewImgBoxSize, previewImgBoxSize, imaging.Lanczos)

	ObjCoords = make(map[int]g143.Rect)
	drawObjects(window, CurrentPage)

	previewX := xPosInt + int(10*textScale)
	if previewX+previewImgBoxSize > wWidth {
		previewX = xPosInt - previewImgBoxSize - int(10*textScale)
	}

	previewY := yPosInt + int(10*textScale)
	if previewY+previewImgBoxSize > wHeight {
		previewY = yPosInt - previewImgBoxSize + int(10*textScale)
	}

	ggCtx := gg.NewContextForImage(tmpPickerFrame)
	ggCtx.DrawImage(foundImg, previewX, previewY)
	ggCtx.Fill()

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

}

func scrollBtnCB(window *glfw.Window, xoff, yoff float64) {

	if scrollEventCount != 5 {
		scrollEventCount += 1
		return
	}

	scrollEventCount = 0

	if xoff == 0 && yoff == -1 && CurrentPage != TotalPages() {
		ObjCoords = make(map[int]g143.Rect)
		drawObjects(window, CurrentPage+1)
		window.SetMouseButtonCallback(mouseBtnCallback)
		window.SetCursorPosCallback(cursorCallback)
	} else if xoff == 0 && yoff == 1 && CurrentPage != 1 {
		ObjCoords = make(map[int]g143.Rect)
		drawObjects(window, CurrentPage-1)
		window.SetMouseButtonCallback(mouseBtnCallback)
		window.SetCursorPosCallback(cursorCallback)
	}

}
