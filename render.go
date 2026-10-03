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
	// Light to dark. Plain ASCII so the picture pastes as text.
	ramp = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
)

func render(w io.Writer, src image.Image, cols int) error {
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

	var line strings.Builder
	for y := 0; y < rows; y++ {
		line.Reset()
		line.Grow(cols + 1)
		for x := 0; x < cols; x++ {
			line.WriteByte(asciiChar(dst, x, y))
		}
		line.WriteByte('\n')
		if _, err := io.WriteString(w, line.String()); err != nil {
			return err
		}
	}
	return nil
}

func asciiChar(img *image.RGBA, x, y int) byte {
	i := img.PixOffset(x, y)
	if img.Pix[i+3] < alphaMin {
		return ' '
	}
	lum := (299*int(img.Pix[i]) + 587*int(img.Pix[i+1]) + 114*int(img.Pix[i+2])) / 1000
	return ramp[(255-lum)*(len(ramp)-1)/255]
}
