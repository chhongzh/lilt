package main

import (
	"sync/atomic"

	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Inputs on the board, indexing both the keys table and the indicator.
const (
	padUp = iota
	padDown
	padLeft
	padRight
	padMid
	padSet
	padRst
	padCount
)

// keypad draws the joystick and the SET and RST keys, each lit while its pin
// reads high. It reads the raw levels the main loop samples, so it also
// lights keys the page itself has no use for.
type keypad struct {
	ui.Box
	held  *[padCount]atomic.Bool
	gpio  *[padCount]string
	shown [padCount]bool
	fill  [padCount]anim.Value
}

func newKeypad(held *[padCount]atomic.Bool, gpio *[padCount]string) *keypad {
	k := &keypad{held: held, gpio: gpio}
	for i := range k.fill {
		k.fill[i].SetSpring(anim.Snappy)
		k.fill[i].SetPrecision(0.01)
	}
	return k
}

func (k *keypad) Leaf() {}

func (k *keypad) Measure(max ui.Vec2) ui.Vec2 { return max }

func (k *keypad) Update(float32) {
	for i := range k.held {
		down := k.held[i].Load()
		if down == k.shown[i] {
			continue
		}
		k.shown[i] = down
		if down {
			k.fill[i].Set(1)
		} else {
			k.fill[i].Set(0)
		}
	}
}

func (k *keypad) Draw(c gfx.Canvas) {
	t := ui.CurrentTheme
	b := k.Rect().Pixels()
	r := min(b.H/2-1, 25)
	cx, cy := r+2, b.H/2
	off, pr := r*16/25, r*7/25

	c.StrokeCircle(cx, cy, r, t.Muted)
	centres := [...][2]int16{
		padUp:    {cx, cy - off},
		padDown:  {cx, cy + off},
		padLeft:  {cx - off, cy},
		padRight: {cx + off, cy},
		padMid:   {cx, cy},
	}
	for i, p := range centres {
		f := k.level(i)
		c.StrokeCircle(p[0], p[1], pr, t.Foreground)
		if fr := min(ui.Round(f*float32(pr)), pr); fr > 0 {
			c.FillCircle(p[0], p[1], fr, t.Foreground)
		}
		s := k.gpio[i]
		text(c, t.SmallFont, p[0]-gfx.TextWidth(t.SmallFont, s)/2, p[1], s, ink(f))
	}

	const keyH = 15
	x := cx + r + 6
	side := [...]struct {
		pad  int
		name string
		y    int16
	}{
		{padSet, "SET", cy - keyH - 2},
		{padRst, "RST", cy + 2},
	}
	for _, s := range side {
		box := gfx.Rect{X: x, Y: s.y, W: b.W - x - 2, H: keyH}
		f := k.level(s.pad)
		c.StrokeRoundRect(box, t.Radius, t.Foreground)
		// The fill opens from the middle line outwards, like a key going down.
		fh := min(ui.Round(f*keyH), keyH)
		fill := gfx.Rect{X: box.X, Y: box.Y + (keyH-fh)/2, W: box.W, H: fh}
		if fh > 0 {
			c.FillRoundRect(fill, min(t.Radius, fh/2), t.Foreground)
		}
		caption(c, box, s.name, k.gpio[s.pad], t.Foreground)
		if fh > 0 {
			// Drawn again inside the fill so the text splits along its edge.
			caption(c.Clip(fill), box, s.name, k.gpio[s.pad], t.Background)
		}
	}
}

// caption puts a key's name on the left of box and its GPIO on the right.
func caption(c gfx.Canvas, box gfx.Rect, name, gpio string, col gfx.Color) {
	t := ui.CurrentTheme
	mid := box.Y + box.H/2
	text(c, t.Font, box.X+4, mid, name, col)
	text(c, t.SmallFont, box.X+box.W-4-gfx.TextWidth(t.SmallFont, gpio), mid, gpio, col)
}

func (k *keypad) level(i int) float32 { return max(k.fill[i].Get(), 0) }

// ink flips a pad's number to the background colour once the fill covers it.
func ink(fill float32) gfx.Color {
	if fill > 0.5 {
		return ui.CurrentTheme.Background
	}
	return ui.CurrentTheme.Foreground
}

// text draws s with its line box centred on y.
func text(c gfx.Canvas, f *gfx.Font, x, y int16, s string, col gfx.Color) {
	c.Text(f, x, y-gfx.LineHeight(f)/2+gfx.Ascent(f), s, col)
}
