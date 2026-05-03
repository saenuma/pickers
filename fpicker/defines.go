package main

import (
	"image"

	g143 "github.com/bankole7782/graphics143"
)

const (
	fps = 24

	BackBtn  = 9801
	PageSize = 11 * 4

	DefaultWindowW = 800
	DefaultWindowH = 540
)

var (
	ObjCoords        map[int]g143.Rect
	rootPath         string
	basePath         string
	exts             string
	tmpPickerFrame   image.Image
	scrollEventCount int
	CurrentPage      int

	FontSize float64
)
