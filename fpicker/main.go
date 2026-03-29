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
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/pickers/internal"
)

const (
	fps      = 10
	fontSize = 20

	BackBtn  = 9801
	PageSize = 14 * 4
)

var (
	objCoords        map[int]g143.Rect
	rootPath         string
	basePath         string
	exts             string
	tmpPickerFrame   image.Image
	scrollEventCount int
	CurrentPage      int
)

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
	drawObjects(window, 1)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetCursorPosCallback(cursorCallback)
	window.SetScrollCallback(scrollBtnCB)

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
		tmp = filename[0:15] + "..." + filename[len(filename)-6:]
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

func drawObjects(window *glfw.Window, page int) {
	CurrentPage = page
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := internal.GetDefaultFontPath()
	ggCtx.LoadFontFace(fontPath, 20)

	ggCtx.SetHexColor("#444")
	ggCtx.DrawRectangle(20, 5, 30, 30)
	ggCtx.Fill()
	objCoords[BackBtn] = g143.NewRect(20, 5, 30, 30)

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString("<", 30, 5+fontSize)

	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(basePath, 60, 5+20)

	currentX := 20
	currentY := 50

	toPickFrom := GetPageObjects(page)
	for i, aFile := range toPickFrom {
		shortAFile := shortenObject(aFile)
		// aFileStrW, _ := ggCtx.MeasureString(shortAFile)
		objW := 270

		if strings.HasSuffix(aFile, "/") {
			ggCtx.SetHexColor("#444")
			ggCtx.DrawRoundedRectangle(float64(currentX), float64(currentY), float64(objW)+10, fontSize+20, 4)
			ggCtx.Fill()
			ggCtx.SetHexColor("#fff")
			ggCtx.DrawRoundedRectangle(float64(currentX)+1, float64(currentY)+1, float64(objW)+10-2, fontSize+20-2, 2)
			ggCtx.Fill()
			ggCtx.SetHexColor("#444")
			ggCtx.DrawString(shortAFile, float64(currentX+5), float64(currentY)+fontSize+10)
		} else {
			ggCtx.SetHexColor("#444")
			ggCtx.DrawString(shortAFile, float64(currentX+5), float64(currentY)+fontSize+10)
		}

		aFileRS := g143.Rect{OriginX: currentX, OriginY: currentY, Width: int(objW) + 10, Height: fontSize + 20}
		objCoords[i+1] = aFileRS

		newX := currentX + int(objW) + 10
		if newX > (wWidth - int(objW)) {
			currentY += fontSize + 20 + 10
			currentX = 20
		} else {
			currentX += int(objW) + 20
		}
	}

	pageLabel := fmt.Sprintf("Page %d / %d", CurrentPage, TotalPages())
	pageLabelW, _ := ggCtx.MeasureString(pageLabel)
	pageLabelX := (wWidth - int(pageLabelW)) / 2
	ggCtx.DrawString(pageLabel, float64(pageLabelX), float64(wHeight)-20)

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	tmpPickerFrame = ggCtx.Image()
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
		rootPathTmp := rootPath
		if strings.HasSuffix(rootPath, "/") {
			rootPathTmp = rootPath[:len(rootPath)-1]
		}

		if strings.Count(rootPathTmp, "/") < strings.Count(tmp, "/")+1 {
			objCoords = make(map[int]g143.Rect)
			basePath = tmp
			drawObjects(window, CurrentPage)
		}

		return
	}

	toPickFrom := GetPageObjects(CurrentPage)
	foundObject := toPickFrom[widgetCode-1]

	if strings.HasSuffix(foundObject, "/") {
		objCoords = make(map[int]g143.Rect)
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
	for code, RS := range objCoords {
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

	toPickFrom := getObjects(basePath, exts)
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
	previewImgBoxSize := 300

	foundImg = imaging.Fit(foundImg, previewImgBoxSize, previewImgBoxSize, imaging.Lanczos)

	objCoords = make(map[int]g143.Rect)
	drawObjects(window, CurrentPage)

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
		objCoords = make(map[int]g143.Rect)
		drawObjects(window, CurrentPage+1)
		window.SetCursorPosCallback(cursorCallback)
	} else if xoff == 0 && yoff == 1 && CurrentPage != 1 {
		objCoords = make(map[int]g143.Rect)
		drawObjects(window, CurrentPage-1)
		window.SetCursorPosCallback(cursorCallback)
	}

}
