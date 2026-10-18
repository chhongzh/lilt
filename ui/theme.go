package ui

import (
	"osu-mc/anim"
	"osu-mc/gfx"

	"tinygo.org/x/tinyfont/notosans"
	"tinygo.org/x/tinyfont/proggy"
)

// Theme holds the colours and fonts widgets draw with.
type Theme struct {
	Background gfx.Color
	Foreground gfx.Color
	// Accent draws the cursor on panels that cannot invert pixels.
	Accent gfx.Color
	// Muted is for hints and disabled text.
	Muted     gfx.Color
	Font      *gfx.Font
	SmallFont *gfx.Font
	Radius    int16
}

// CurrentTheme defaults to light-on-dark, which suits OLED panels.
var CurrentTheme = &Theme{
	Background: gfx.Black,
	Foreground: gfx.White,
	Accent:     gfx.White,
	Muted:      gfx.Color{R: 128, G: 128, B: 128, A: 255},
	Font:       &notosans.Notosans12pt,
	SmallFont:  &proggy.TinySZ8pt7b,
	Radius:     3,
}

// Cursor springs: the edge on the side of travel uses CursorLead and the
// opposite edge CursorTrail, so the cursor stretches toward its destination.
var (
	CursorLead  = &anim.Spring{Response: 0.2, Damping: 0.9}
	CursorTrail = &anim.Spring{Response: 0.38, Damping: 0.9}
)
