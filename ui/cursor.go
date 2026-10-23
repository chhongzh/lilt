package ui

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
)

// CursorShaper lets the focused node change the cursor's rect, given in
// screen pixels. Hold-to-confirm buttons use it to shrink the cursor.
type CursorShaper interface{ CursorShape(r Rect) Rect }

const (
	// Press squeeze per side, horizontally; vertically it is half.
	pressDepth = 3
	// Leap: how high the cursor lifts, and how much each side grows there.
	leapHeight = 3
	leapGrow   = 2
	// Pointer lean: how far the cursor follows a hovering pointer, and its cap.
	leanFactor = 0.12
	leanMax    = 3
)

// Cursor is the single highlight that marks focus. Each edge moves on its
// own so the highlight stretches toward where focus went and then catches up.
//
// It lives in screen space: when focus moves to another screen it does not
// vanish with the old one, it leaps across. Pressing squeezes it; on release
// it springs back with a wobble.
type Cursor struct {
	l, t, r, b anim.Value
	press      anim.Value
	// Displacement on top of the focus rect: lift for Leap (up is
	// positive), shiftX and shiftY for Lean and Nudge.
	lift           anim.Value
	shiftX, shiftY anim.Value
	leanX, leanY   float32

	target     Rect
	has        bool
	hasPointer bool
	px, py     int16
}

func (c *Cursor) init() {
	c.press.SetPrecision(0.01)
	c.lift.SetSpring(CursorLeap)
	c.shiftX.SetSpring(CursorLean)
	c.shiftY.SetSpring(CursorLean)
}

// pressed squeezes quickly and without bounce, saving the wobble for the release.
func (c *Cursor) pressed(p Phase) {
	switch p {
	case Press:
		c.press.SetSpring(anim.Interactive)
		c.press.Set(1)
	case Release:
		c.press.SetSpring(CursorRelease)
		c.press.Set(0)
	}
}

// Lean moves the cursor by dx, dy and holds it there: anticipation while a
// key is down, before the release commits the move. Lean(0, 0) lets go.
func (c *Cursor) Lean(dx, dy float32) { c.leanX, c.leanY = dx, dy }

// Leap makes the cursor jump: it lifts off, grows as if coming closer and
// lands with a squash. Screen transitions call it.
func (c *Cursor) Leap() { c.lift.Bump(leapHeight) }

// Nudge throws the cursor about px sideways and lets it spring back.
func (c *Cursor) Nudge(px float32) { c.shiftX.Bump(px) }

// clear forgets the cursor entirely, so after a reset the next focus grows
// it afresh instead of matching, or keeping on screen, a rect from before.
func (c *Cursor) clear() {
	for _, v := range [...]*anim.Value{&c.l, &c.t, &c.r, &c.b, &c.press, &c.lift, &c.shiftX, &c.shiftY} {
		v.Jump(0)
	}
	c.target, c.has, c.hasPointer = Rect{}, false, false
	c.leanX, c.leanY = 0, 0
}

func (c *Cursor) pointerAt(x, y int16) {
	c.hasPointer, c.px, c.py = true, x, y
}

func (c *Cursor) update() {
	f := focus.current
	if f == nil || !shown(f) {
		if c.has {
			c.has = false
			m := c.target.Center()
			c.moveTo(Rect{m.X, m.Y, 0, 0})
		}
		c.shiftX.Set(0)
		c.shiftY.Set(0)
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
	lx, ly := c.leanX, c.leanY
	if c.hasPointer {
		m := r.Center()
		lx += clamp((float32(c.px)-m.X)*leanFactor, -leanMax, leanMax)
		ly += clamp((float32(c.py)-m.Y)*leanFactor, -leanMax, leanMax)
	}
	c.shiftX.Set(lx)
	c.shiftY.Set(ly)
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
	lift := c.lift.Get()
	in := c.press.Get()*pressDepth - lift*leapGrow/leapHeight
	dx, dy := c.shiftX.Get(), c.shiftY.Get()-lift
	r := Rect{c.l.Get() + in + dx, c.t.Get() + in/2 + dy, 0, 0}
	r.W = c.r.Get() - in + dx - r.X
	r.H = c.b.Get() - in/2 + dy - r.Y
	if r.W < 1 || r.H < 1 {
		return
	}
	p := r.Pixels()
	rad := min(radiusPx(CurrentTheme.Radius), p.H/2, p.W/2)
	if cv.CanInvert() {
		cv.InvertRoundRect(p, rad)
	} else {
		cv.StrokeRoundRect(p, rad, CurrentTheme.Accent)
	}
}

// screenRect returns d's rect in screen pixels exactly as paint will place
// it, clipped by its clipping ancestors. When nothing is visible the unclipped
// rect is returned so the cursor does not fly to the corner.
func screenRect(d Node) (Rect, bool) {
	x, y := pixelOrigin(d)
	p := pixelRect(d.box())
	own := gfx.Rect{X: x, Y: y, W: p.W, H: p.H}
	vis := own.Intersect(gfx.Rect{W: scrW, H: scrH})
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
func pixelOrigin(d Node) (x, y int16) {
	for n := d; n != nil; n = n.box().parent {
		p := pixelRect(n.box())
		x += p.X
		y += p.Y
	}
	return
}
