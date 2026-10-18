package ui

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"

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

// Cursor motion. The edge on the side of travel eases with CursorLead and
// the opposite edge with the slower CursorTrail, so the cursor stretches
// toward its destination and contracts as it arrives.
var (
	CursorLead  = &anim.Spring{Response: 0.22, Ease: true}
	CursorTrail = &anim.Spring{Response: 0.52, Ease: true}
	// CursorRelease brings the cursor back from a press with a wobble; the
	// press itself uses anim.Interactive.
	CursorRelease = &anim.Spring{Response: 0.57, Damping: 0.4}
	// CursorLeap lifts the cursor across screens and lands it with a squash.
	CursorLeap = &anim.Spring{Response: 0.63, Damping: 0.45}
	// CursorLean moves the cursor toward a held key and back.
	CursorLean = &anim.Spring{Response: 0.35, Damping: 0.7}
)
