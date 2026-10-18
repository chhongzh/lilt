package ui

import (
	"osu-mc/anim"
	"osu-mc/gfx"
)

// CursorShaper lets the focused node change the cursor's rect, given in
// screen pixels. Hold-to-confirm buttons use it to shrink the cursor.
type CursorShaper interface{ CursorShape(r Rect) Rect }

const (
	pressInset = 1.5
	// Pointer lean: how far the cursor follows a hovering pointer, and its cap.
	leanFactor = 0.12
	leanMax    = 3
)

// Cursor is the single highlight that marks focus. Each edge springs on its
// own so the highlight stretches toward where focus went and then catches up.
type Cursor struct {
	l, t, r, b anim.Value
	press      anim.Value
	leanX      anim.Value
	leanY      anim.Value

	target     Rect
	has        bool
	hasPointer bool
	px, py     int16
}

func (c *Cursor) init() {
	c.press.SetSpring(anim.Interactive)
	c.press.SetPrecision(0.01)
	c.leanX.SetSpring(anim.Interactive)
	c.leanY.SetSpring(anim.Interactive)
}

func (c *Cursor) pressed(p Phase) {
	switch p {
	case Press:
		c.press.Set(1)
	case Release:
		c.press.Set(0)
	}
}

func (c *Cursor) pointerAt(x, y int16) {
	c.hasPointer, c.px, c.py = true, x, y
}

func (c *Cursor) update() {
	f := app.Focus.current
	if f == nil || !shown(f) {
		if c.has {
			c.has = false
			m := c.target.Center()
			c.moveTo(Rect{m.X, m.Y, 0, 0})
		}
		c.leanX.Set(0)
		c.leanY.Set(0)
		return
	}
	r, _ := screenRect(f)
	if s, is := f.(CursorShaper); is {
		r = s.CursorShape(r)
	}
	if !c.has {
		// Grow out of the target's centre rather than flying in from the old spot.
		m := r.Center()
		c.l.Jump(m.X)
		c.r.Jump(m.X)
		c.t.Jump(m.Y)
		c.b.Jump(m.Y)
		c.has = true
		Invalidate()
	}
	if r != c.target {
		c.moveTo(r)
	}
	var lx, ly float32
	if c.hasPointer {
		m := r.Center()
		lx = clamp((float32(c.px)-m.X)*leanFactor, -leanMax, leanMax)
		ly = clamp((float32(c.py)-m.Y)*leanFactor, -leanMax, leanMax)
	}
	c.leanX.Set(lx)
	c.leanY.Set(ly)
}

func (c *Cursor) moveTo(r Rect) {
	old, m := c.target.Center(), r.Center()
	edgeSpring(&c.l, &c.r, m.X-old.X)
	edgeSpring(&c.t, &c.b, m.Y-old.Y)
	c.target = r
	c.l.Set(r.X)
	c.t.Set(r.Y)
	c.r.Set(r.Right())
	c.b.Set(r.Bottom())
}

// edgeSpring gives the edge facing the direction of travel the faster spring.
func edgeSpring(lo, hi *anim.Value, d float32) {
	switch {
	case d > 0:
		lo.SetSpring(CursorTrail)
		hi.SetSpring(CursorLead)
	case d < 0:
		lo.SetSpring(CursorLead)
		hi.SetSpring(CursorTrail)
	default:
		lo.SetSpring(CursorLead)
		hi.SetSpring(CursorLead)
	}
}

func (c *Cursor) draw(cv gfx.Canvas) {
	in := c.press.Get() * pressInset
	dx, dy := c.leanX.Get(), c.leanY.Get()
	r := Rect{c.l.Get() + in + dx, c.t.Get() + in + dy, 0, 0}
	r.W = c.r.Get() - in + dx - r.X
	r.H = c.b.Get() - in + dy - r.Y
	if r.W < 1 || r.H < 1 {
		return
	}
	p := r.Pixels()
	rad := min(CurrentTheme.Radius, p.H/2, p.W/2)
	if cv.CanInvert() {
		cv.InvertRoundRect(p, rad)
	} else {
		cv.StrokeRoundRect(p, rad, CurrentTheme.Accent)
	}
}

// screenRect returns d's rect in screen pixels exactly as paint will place
// it, clipped by its clipping ancestors. When nothing is visible the unclipped
// rect is returned so the cursor does not fly to the corner.
func screenRect(d Drawable) (Rect, bool) {
	x, y := pixelOrigin(d)
	p := pixelRect(d.box())
	own := gfx.Rect{X: x, Y: y, W: p.W, H: p.H}
	vis := own.Intersect(gfx.Rect{W: app.w, H: app.h})
	for n := d.box().parent; n != nil; n = n.box().parent {
		b := n.box()
		if b.flags&flagClip == 0 {
			continue
		}
		cx, cy := pixelOrigin(n)
		c := pixelRect(b)
		vis = vis.Intersect(gfx.Rect{X: cx, Y: cy, W: c.W, H: c.H})
	}
	if vis.Empty() {
		vis = own
	}
	return Rect{float32(vis.X), float32(vis.Y), float32(vis.W), float32(vis.H)}, !vis.Empty()
}

// pixelRect is b's rect plus its offset, rounded the way paint rounds it.
func pixelRect(b *Box) gfx.Rect {
	return Rect{b.rect.X + b.offsetX.Get(), b.rect.Y + b.offsetY.Get(), b.rect.W, b.rect.H}.Pixels()
}

// pixelOrigin is d's top-left corner on screen as paint computes it.
func pixelOrigin(d Drawable) (x, y int16) {
	for n := d; n != nil; n = n.box().parent {
		p := pixelRect(n.box())
		x += p.X
		y += p.Y
	}
	return
}
