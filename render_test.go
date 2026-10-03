package main

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func TestRender(t *testing.T) {
	wb := image.NewRGBA(image.Rect(0, 0, 2, 1))
	wb.SetRGBA(0, 0, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	wb.SetRGBA(1, 0, color.RGBA{A: 255})

	if got := mustRender(t, wb, false, false); got != " @\n" {
		t.Fatalf("short: got %q", got)
	}
	if got := mustRender(t, wb, true, false); got != "@ \n" {
		t.Fatalf("invert: got %q", got)
	}
	if got := mustRender(t, wb, false, true); got != " $\n" {
		t.Fatalf("long: got %q", got)
	}

	gray := image.NewRGBA(image.Rect(0, 0, 2, 1))
	gray.SetRGBA(0, 0, color.RGBA{R: 10, G: 10, B: 10, A: 255})
	gray.SetRGBA(1, 0, color.RGBA{R: 40, G: 40, B: 40, A: 255})
	if got := mustRender(t, gray, false, false); got != "@ \n" {
		t.Fatalf("contrast: got %q", got)
	}
}

func mustRender(t *testing.T, src image.Image, invert, long bool) string {
	t.Helper()
	var buf bytes.Buffer
	if err := render(&buf, src, 2, invert, long); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
