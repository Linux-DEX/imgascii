package main

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

func TestRenderWhiteAndBlack(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 2, 1))
	src.SetRGBA(0, 0, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	src.SetRGBA(1, 0, color.RGBA{A: 255})

	var buf bytes.Buffer
	if err := render(&buf, src, 2); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	want := " $\n"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
