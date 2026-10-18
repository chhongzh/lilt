package gfx

// Canvas is a translated, clipped view of a Screen. It is a small value:
// Translate and Clip return new canvases, so nested drawing never has to
// restore state.
type Canvas struct {
	s      Screen
	inv    Inverter
	clip   Rect // screen coordinates
	ox, oy int16
}

func NewCanvas(s Screen) Canvas {
	w, h := s.Size()
	inv, _ := s.(Inverter)
	return Canvas{s: s, inv: inv, clip: Rect{0, 0, w, h}}
}

// Translate moves the origin by dx, dy.
func (c Canvas) Translate(dx, dy int16) Canvas {
	c.ox += dx
	c.oy += dy
	return c
}

// Clip restricts drawing to r, given in local coordinates.
func (c Canvas) Clip(r Rect) Canvas {
	c.clip = c.clip.Intersect(r.Offset(c.ox, c.oy))
	return c
}

// Absolute returns the canvas with its origin at the screen corner, keeping the clip.
func (c Canvas) Absolute() Canvas {
	c.ox, c.oy = 0, 0
	return c
}

// Origin is the local origin in screen coordinates.
func (c Canvas) Origin() (x, y int16) { return c.ox, c.oy }

// Bounds is the clip in local coordinates.
func (c Canvas) Bounds() Rect { return c.clip.Offset(-c.ox, -c.oy) }

// Visible reports whether any part of r, in local coordinates, survives the clip.
func (c Canvas) Visible(r Rect) bool { return c.clip.Overlaps(r.Offset(c.ox, c.oy)) }

// CanInvert reports whether InvertRect and InvertRoundRect have any effect.
func (c Canvas) CanInvert() bool { return c.inv != nil }

func (c Canvas) Pixel(x, y int16, col Color) {
	x += c.ox
	y += c.oy
	if c.clip.Contains(x, y) {
		c.s.SetPixel(x, y, col)
	}
}

func (c Canvas) FillRect(r Rect, col Color) {
	r = c.clip.Intersect(r.Offset(c.ox, c.oy))
	if !r.Empty() {
		c.s.FillRect(r.X, r.Y, r.W, r.H, col)
	}
}

func (c Canvas) InvertRect(r Rect) {
	if c.inv == nil {
		return
	}
	r = c.clip.Intersect(r.Offset(c.ox, c.oy))
	if !r.Empty() {
		c.inv.InvertRect(r.X, r.Y, r.W, r.H)
	}
}

func (c Canvas) HLine(x, y, w int16, col Color) { c.FillRect(Rect{x, y, w, 1}, col) }
func (c Canvas) VLine(x, y, h int16, col Color) { c.FillRect(Rect{x, y, 1, h}, col) }

func (c Canvas) StrokeRect(r Rect, col Color) {
	if r.Empty() {
		return
	}
	c.HLine(r.X, r.Y, r.W, col)
	c.HLine(r.X, r.Bottom()-1, r.W, col)
	c.VLine(r.X, r.Y+1, r.H-2, col)
	c.VLine(r.Right()-1, r.Y+1, r.H-2, col)
}

// FillRoundRect fills row spans so fills and inversions cost one call per row.
func (c Canvas) FillRoundRect(r Rect, radius int16, col Color) {
	c.roundSpans(r, radius, col, false)
}

func (c Canvas) InvertRoundRect(r Rect, radius int16) {
	if c.inv == nil {
		return
	}
	c.roundSpans(r, radius, Color{}, true)
}

// roundSpans takes a flag instead of a callback so the per-frame cursor draw
// cannot allocate a closure.
func (c Canvas) roundSpans(r Rect, radius int16, col Color, invert bool) {
	if r.Empty() {
		return
	}
	radius = min(radius, r.W/2, r.H/2)
	if radius <= 0 {
		c.span(r, col, invert)
		return
	}
	body := Rect{r.X, r.Y + radius, r.W, r.H - 2*radius}
	if !body.Empty() {
		c.span(body, col, invert)
	}
	for i := int16(0); i < radius; i++ {
		in := cornerInset(radius, i)
		w := r.W - 2*in
		c.span(Rect{r.X + in, r.Y + i, w, 1}, col, invert)
		c.span(Rect{r.X + in, r.Bottom() - 1 - i, w, 1}, col, invert)
	}
}

func (c Canvas) span(r Rect, col Color, invert bool) {
	if invert {
		c.InvertRect(r)
	} else {
		c.FillRect(r, col)
	}
}

func (c Canvas) StrokeRoundRect(r Rect, radius int16, col Color) {
	if r.Empty() {
		return
	}
	radius = min(radius, r.W/2, r.H/2)
	if radius <= 0 {
		c.StrokeRect(r, col)
		return
	}
	c.HLine(r.X+radius, r.Y, r.W-2*radius, col)
	c.HLine(r.X+radius, r.Bottom()-1, r.W-2*radius, col)
	c.VLine(r.X, r.Y+radius, r.H-2*radius, col)
	c.VLine(r.Right()-1, r.Y+radius, r.H-2*radius, col)
	prev := radius
	for i := int16(0); i < radius; i++ {
		in := cornerInset(radius, i)
		// Join to the previous row's inset so steep parts of the arc stay connected.
		w := max(prev-in, 1)
		if i == 0 {
			w = radius - in
		}
		yt, yb := r.Y+i, r.Bottom()-1-i
		c.HLine(r.X+in, yt, w, col)
		c.HLine(r.Right()-in-w, yt, w, col)
		c.HLine(r.X+in, yb, w, col)
		c.HLine(r.Right()-in-w, yb, w, col)
		prev = in
	}
}

// cornerInset is how far row i of a corner with radius r is indented. The
// +r bias rounds to the nearest pixel, which keeps small radii from looking square.
func cornerInset(r, i int16) int16 {
	dy := int32(r - i)
	rr := int32(r)
	dx := isqrt(rr*rr - dy*dy + rr)
	return r - int16(dx)
}

func isqrt(n int32) int32 {
	if n <= 0 {
		return 0
	}
	x := n
	y := (x + 1) / 2
	for y < x {
		x = y
		y = (x + n/x) / 2
	}
	return x
}

func (c Canvas) FillCircle(cx, cy, r int16, col Color) {
	c.FillRoundRect(Rect{cx - r, cy - r, 2*r + 1, 2*r + 1}, r+1, col)
}

func (c Canvas) StrokeCircle(cx, cy, r int16, col Color) {
	c.StrokeRoundRect(Rect{cx - r, cy - r, 2*r + 1, 2*r + 1}, r+1, col)
}

// Line uses Bresenham; good enough for icons such as check marks.
func (c Canvas) Line(x0, y0, x1, y1 int16, col Color) {
	dx, dy := x1-x0, y1-y0
	sx, sy := int16(1), int16(1)
	if dx < 0 {
		dx, sx = -dx, -1
	}
	if dy < 0 {
		dy, sy = -dy, -1
	}
	err := dx - dy
	for {
		c.Pixel(x0, y0, col)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
