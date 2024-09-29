package main

import (
	"image"
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

func allColorImg() *image.RGBA {
	width := 30
	height := 360

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			h := float64(y)
			// l := math.Abs(float64(y)/100.0 - 1.0)
			r, g, b := colorful.Hsl(h, 1, 0.5).RGB255()
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

func allColorImg2(hue int) *image.RGBA {
	width := 360
	height := 360

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			h := float64(hue)
			l := math.Abs(float64(x)/float64(width) - 1.0)
			s := math.Abs(float64(y)/float64(height) - 1.0)
			r, g, b := colorful.Hsl(h, s, l).RGB255()
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}
