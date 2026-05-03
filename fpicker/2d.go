package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/saenuma/pickers/internal"
)

type Ctx struct {
	WindowWidth  int
	WindowHeight int
	ggCtx        *gg.Context
	OldFrame     image.Image
}

func New2dCtx(wWidth, wHeight int) Ctx {
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

	ctx := Ctx{WindowWidth: wWidth, WindowHeight: wHeight, ggCtx: ggCtx}
	return ctx
}

func Continue2dCtx(img image.Image) Ctx {
	ggCtx := gg.NewContextForImage(img)

	// load font
	fontPath := internal.GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, FontSize)
	if err != nil {
		panic(err)
	}

	ctx := Ctx{WindowWidth: img.Bounds().Dx(), WindowHeight: img.Bounds().Dy(), ggCtx: ggCtx}
	return ctx
}

func (ctx *Ctx) drawButtonA(btnId, originX, originY int, text, textColor, bgColor string) g143.Rect {
	// textScale := internal.GetTextScale()

	// draw bounding rect
	textW, textH := ctx.ggCtx.MeasureString(text)
	width, height := textW+(FontSize), textH+(FontSize)
	ctx.ggCtx.SetHexColor(bgColor)
	ctx.ggCtx.DrawRectangle(float64(originX), float64(originY), float64(width), float64(height))
	ctx.ggCtx.Fill()

	// draw text
	ctx.ggCtx.SetHexColor(textColor)
	ctx.ggCtx.DrawString(text, float64(originX)+(FontSize/2), float64(originY)+FontSize+(FontSize/4))

	// save dimensions
	btnARect := g143.NewRect(originX, originY, int(width), int(height))
	ObjCoords[btnId] = btnARect
	return btnARect
}

func nextX(aRect g143.Rect, margin int) int {
	textScale := internal.GetTextScale()
	newMargin := float64(margin) * textScale
	return aRect.OriginX + aRect.Width + int(newMargin)
}

func nextY(aRect g143.Rect, margin int) int {
	textScale := internal.GetTextScale()
	newMargin := float64(margin) * textScale
	return aRect.OriginY + aRect.Height + int(newMargin)
}

func (ctx *Ctx) windowRect() g143.Rect {
	return g143.NewRect(0, 0, ctx.WindowWidth, ctx.WindowHeight)
}
