// Package silent is a display driver that draws into memory and shows
// nothing. It is a one-bit frame buffer like the OLED it stands in for, with
// helpers to turn the buffer into images and animated GIFs, so a UI can be
// rendered headless, in tests or for the README.
package silent

import (
	"image"
	"image/color"

	"github.com/chhongzh/lilt/gfx"
)

// Palette is the off and on colour of a one-bit screen, in that order.
var Palette = color.Palette{
	color.RGBA{8, 10, 14, 255},
	color.RGBA{200, 230, 255, 255},
}

// Screen is a gfx.Screen and gfx.Inverter backed by a slice of bools. It is
// not safe for concurrent use: read it from the goroutine that draws, or
// after the frame loop has stopped.
type Screen struct {
	w, h int
	px   []bool
}

func New(w, h int) *Screen { return &Screen{w: w, h: h, px: make([]bool, w*h)} }

func (s *Screen) Size() (w, h int16) { return int16(s.w), int16(s.h) }

func (s *Screen) SetPixel(x, y int16, c gfx.Color) {
	if int(x) >= 0 && int(y) >= 0 && int(x) < s.w && int(y) < s.h {
		s.px[int(y)*s.w+int(x)] = gfx.IsLit(c)
	}
}

func (s *Screen) FillRect(x, y, w, h int16, c gfx.Color) {
	on := gfx.IsLit(c)
	s.each(x, y, w, h, func(i int) { s.px[i] = on })
}

func (s *Screen) InvertRect(x, y, w, h int16) {
	s.each(x, y, w, h, func(i int) { s.px[i] = !s.px[i] })
}

// Flush does nothing: the buffer is always up to date.
func (s *Screen) Flush() error { return nil }

func (s *Screen) each(x, y, w, h int16, fn func(i int)) {
	x0, y0 := max(int(x), 0), max(int(y), 0)
	x1, y1 := min(int(x)+int(w), s.w), min(int(y)+int(h), s.h)
	for j := y0; j < y1; j++ {
		for i := x0; i < x1; i++ {
			fn(j*s.w + i)
		}
	}
}

// Pixels is the live buffer, row by row, true for lit.
func (s *Screen) Pixels() []bool { return s.px }

// Image renders the whole screen, each pixel scale times larger.
func (s *Screen) Image(scale int) *image.Paletted {
	return s.Region(image.Rect(0, 0, s.w, s.h), scale)
}

// Region renders the screen pixels inside r, scaled. The image keeps the
// position r has in the scaled screen, so a GIF frame can carry just the part
// that changed.
func (s *Screen) Region(r image.Rectangle, scale int) *image.Paletted {
	img := image.NewPaletted(image.Rect(r.Min.X*scale, r.Min.Y*scale, r.Max.X*scale, r.Max.Y*scale), Palette)
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			if s.px[(y/scale)*s.w+x/scale] {
				img.Pix[img.PixOffset(x, y)] = 1
			}
		}
	}
	return img
}

// Diff is the bounding box, in screen pixels, of what differs between the
// buffer and prev, a copy taken earlier from Pixels.
func (s *Screen) Diff(prev []bool) image.Rectangle {
	var r image.Rectangle
	for i, p := range s.px {
		if p != prev[i] {
			r = r.Union(image.Rect(i%s.w, i/s.w, i%s.w+1, i/s.w+1))
		}
	}
	return r
}
