// Package gfx is the drawing layer: the Screen interface a display driver
// implements, and a Canvas that adds clipping, translation and shapes on top.
package gfx

import "image/color"

// Color is plain RGBA so colour panels and monochrome panels share one API;
// monochrome drivers decide what counts as lit with IsLit.
type Color = color.RGBA

var (
	Black = Color{A: 255}
	White = Color{R: 255, G: 255, B: 255, A: 255}
)

// IsLit reports whether c should light a pixel on a one-bit panel.
func IsLit(c Color) bool {
	return int(c.R)*3+int(c.G)*6+int(c.B) >= 128*10
}

// Screen is what a display driver must provide. Coordinates outside the
// screen must be ignored rather than rejected.
type Screen interface {
	Size() (w, h int16)
	SetPixel(x, y int16, c Color)
	FillRect(x, y, w, h int16, c Color)
	// Flush pushes the finished frame to the panel.
	Flush() error
}

// Inverter is optional. Panels that can invert pixels get an inverting
// cursor drawn over content; others get a highlight drawn beneath it.
type Inverter interface {
	InvertRect(x, y, w, h int16)
}

// Displayer matches tinygo.org/x/drivers.Displayer without importing it.
type Displayer interface {
	Size() (x, y int16)
	SetPixel(x, y int16, c color.RGBA)
	Display() error
}

// FromDisplayer adapts any tinygo driver. Drivers with a FillRectangle
// method fill through it; the rest fall back to one SetPixel per pixel.
func FromDisplayer(d Displayer) Screen {
	if f, ok := d.(rectFiller); ok {
		return filler{displayer{d}, f}
	}
	return displayer{d}
}

// rectFiller is the native fill most colour drivers and ssd1306 provide.
type rectFiller interface {
	FillRectangle(x, y, w, h int16, c color.RGBA) error
}

type displayer struct{ Displayer }

func (d displayer) FillRect(x, y, w, h int16, c Color) {
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			d.SetPixel(i, j, c)
		}
	}
}

type filler struct {
	displayer
	f rectFiller
}

// FillRect ignores the error: Canvas has already clipped the rectangle to
// the screen, which is all the drivers check.
func (d filler) FillRect(x, y, w, h int16, c Color) { d.f.FillRectangle(x, y, w, h, c) }

func (d displayer) Flush() error { return d.Display() }
