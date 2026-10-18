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

// FromDisplayer adapts any tinygo driver. FillRect falls back to SetPixel,
// so drivers with a native fill deserve their own adapter.
func FromDisplayer(d Displayer) Screen { return displayer{d} }

type displayer struct{ Displayer }

func (d displayer) FillRect(x, y, w, h int16, c Color) {
	for j := y; j < y+h; j++ {
		for i := x; i < x+w; i++ {
			d.SetPixel(i, j, c)
		}
	}
}

func (d displayer) Flush() error { return d.Display() }
