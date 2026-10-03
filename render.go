package main

import (
	"fmt"
	"image"
	"io"
	"strings"

	"golang.org/x/image/draw"
)

const (
	alphaMin = 16
	// Light to dark. The short ramp is the default; it stays readable when pasted.
	rampShort = " .:-=+*#%@"
	rampLong  = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
)

func render(w io.Writer, src image.Image, cols int, invert, long bool) error {
	sw, sh := src.Bounds().Dx(), src.Bounds().Dy()
	if sw == 0 || sh == 0 {
		return fmt.Errorf("empty image")
	}
	if cols < 1 {
		return fmt.Errorf("width must be positive")
	}
	// A terminal cell is about twice as tall as it is wide.
	rows := int(int64(cols) * int64(sh) / int64(sw) / 2)
	if rows < 1 {
		rows = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, cols, rows))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	chars := rampShort
	if long {
		chars = rampLong
	}
	minLum, maxLum := 256, -1
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			lum, ok := pixelLum(dst, x, y)
			if !ok {
				continue
			}
			if lum < minLum {
				minLum = lum
			}
			if lum > maxLum {
				maxLum = lum
			}
		}
	}

	var line strings.Builder
	for y := 0; y < rows; y++ {
		line.Reset()
		line.Grow(cols + 1)
		for x := 0; x < cols; x++ {
			line.WriteByte(asciiChar(dst, x, y, minLum, maxLum, invert, chars))
		}
		line.WriteByte('\n')
		if _, err := io.WriteString(w, line.String()); err != nil {
			return err
		}
	}
	return nil
}

func pixelLum(img *image.RGBA, x, y int) (int, bool) {
	i := img.PixOffset(x, y)
	if img.Pix[i+3] < alphaMin {
		return 0, false
	}
	lum := (299*int(img.Pix[i]) + 587*int(img.Pix[i+1]) + 114*int(img.Pix[i+2])) / 1000
	return lum, true
}

func asciiChar(img *image.RGBA, x, y, minLum, maxLum int, invert bool, chars string) byte {
	lum, ok := pixelLum(img, x, y)
	if !ok {
		return ' '
	}
	if maxLum > minLum {
		lum = (lum - minLum) * 255 / (maxLum - minLum)
	}
	if invert {
		lum = 255 - lum
	}
	return chars[(255-lum)*(len(chars)-1)/255]
}
