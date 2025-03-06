package main

import (
	"image"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/saenuma/pickers/internal"
)

type Ctx struct {
	WindowWidth  int
	WindowHeight int
	ggCtx        *gg.Context
	ObjCoords    *map[int]g143.Rect
}

func New2dCtx(wWidth, wHeight int, objCoords *map[int]g143.Rect) Ctx {
	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := internal.GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, FontSize)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: wWidth, WindowHeight: wHeight, ggCtx: ggCtx,
		ObjCoords: objCoords}
	return ctx
}

func Continue2dCtx(img image.Image, objCoords *map[int]g143.Rect) Ctx {
	ggCtx := gg.NewContextForImage(img)

	// load font
	fontPath := internal.GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, FontSize)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: img.Bounds().Dx(), WindowHeight: img.Bounds().Dy(), ggCtx: ggCtx,
		ObjCoords: objCoords}
	return ctx
}

func (ctx *Ctx) drawTextInput(inputId, originX, originY, inputWidth, height int, values string) g143.Rect {
	ctx.ggCtx.SetHexColor("#ddd")
	ctx.ggCtx.DrawRoundedRectangle(float64(originX), float64(originY)+float64(height), float64(inputWidth), 2, 2)
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX)-2, float64(originY)-2, float64(inputWidth)+4, float64(height)+4)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: inputWidth, Height: height, OriginX: originX, OriginY: originY}
	(*ctx.ObjCoords)[inputId] = entryRect

	if len(strings.TrimSpace(values)) != 0 {
		strs := strings.Split(values, "\n")
		currentY := originY
		for _, str := range strs {
			ctx.ggCtx.SetHexColor("#444")
			ctx.ggCtx.DrawString(str, float64(originX), float64(currentY)+FontSize)
			currentY += FontSize + LineSpacing
		}
	}
	return entryRect
}

func (ctx *Ctx) drawTextInputWithErrors(inputId, originX, originY, inputWidth, height int, values string) g143.Rect {
	ctx.ggCtx.SetHexColor("#ddd")
	ctx.ggCtx.DrawRoundedRectangle(float64(originX), float64(originY)+float64(height), float64(inputWidth), 2, 2)
	ctx.ggCtx.Fill()

	ctx.ggCtx.SetHexColor("#fff")
	ctx.ggCtx.DrawRectangle(float64(originX)-2, float64(originY)-2, float64(inputWidth)+4, float64(height)+4)
	ctx.ggCtx.Fill()

	entryRect := g143.Rect{Width: inputWidth, Height: height, OriginX: originX, OriginY: originY}
	(*ctx.ObjCoords)[inputId] = entryRect

	if len(strings.TrimSpace(values)) != 0 {
		strs := strings.Split(values, "\n")
		currentY := originY
		for _, str := range strs {
			ctx.ggCtx.SetHexColor("#444")
			ctx.ggCtx.DrawString(str, float64(originX), float64(currentY)+FontSize)
			currentY += FontSize + LineSpacing
		}

		currentY = originY
		ctx.ggCtx.SetDash(FontSize/5, FontSize/5)
		ctx.ggCtx.SetLineWidth(3)

		for _, str := range strs {

			strParts := strings.Fields(str)
			spellcheckResults := findWordsNotInDict(str)
			for i, aResult := range spellcheckResults {
				if !aResult.Passed {
					tmpStr := strings.Join(strParts[:i], " ")
					if i != 0 {
						tmpStr += " "
					}
					tmpStrW, _ := ctx.ggCtx.MeasureString(tmpStr)
					wordW, _ := ctx.ggCtx.MeasureString(aResult.Word)
					x1 := originX + int(tmpStrW)
					y1 := currentY + FontSize + 5
					x2 := x1 + int(wordW)
					if aResult.Minor {
						ctx.ggCtx.SetHexColor("#5F8F9A")
					} else {
						ctx.ggCtx.SetHexColor("#A74747")
					}
					ctx.ggCtx.DrawLine(float64(x1), float64(y1), float64(x2), float64(y1))
					ctx.ggCtx.Stroke()
				}
			}
			currentY += FontSize + LineSpacing

		}
	}
	return entryRect
}

func (ctx *Ctx) windowRect() g143.Rect {
	return g143.NewRect(0, 0, ctx.WindowWidth, ctx.WindowHeight)
}
