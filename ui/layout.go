package ui

// The whole tree is laid out again whenever any property changes. Trees on a
// microcontroller screen are small enough that tracking dirty subtrees would
// cost more code than it saves time.
var needsLayout, needsDraw bool

// Relayout schedules a layout pass and a redraw.
func Relayout() {
	needsLayout = true
	needsDraw = true
}

// Invalidate schedules a redraw.
func Invalidate() { needsDraw = true }

func resolve(s Size, avail float32) (v float32, auto bool) {
	switch s.mode {
	case fixedSize:
		return s.v, false
	case relativeSize:
		if avail >= Unbounded {
			return 0, true
		}
		return s.v * avail, false
	}
	return 0, true
}

// Measure returns the size d wants inside max, the space left after its margin.
func Measure(d Node, max Vec2) Vec2 {
	b := d.box()
	w, autoW := resolve(b.width, max.X)
	h, autoH := resolve(b.height, max.Y)
	if autoW || autoH {
		p := padPx(b)
		inner := Vec2{max.X, max.Y}
		if !autoW {
			inner.X = w
		}
		if !autoH {
			inner.Y = h
		}
		inner.X -= p.Horizontal()
		inner.Y -= p.Vertical()
		c := contentSize(d, inner)
		if autoW {
			w = c.X + p.Horizontal()
		}
		if autoH {
			h = c.Y + p.Vertical()
		}
	}
	return Vec2{w, h}
}

func contentSize(d Node, max Vec2) Vec2 {
	if m, ok := d.(Measurer); ok {
		return m.Measure(max)
	}
	if l, ok := d.(Layouter); ok {
		return l.MeasureChildren(max)
	}
	var s Vec2
	for _, c := range d.box().children {
		cb := c.box()
		if !cb.visible() {
			continue
		}
		m := marginPx(cb)
		cs := Measure(c, shrink(max, m))
		s.X = max32(s.X, cs.X+m.Horizontal()+cb.pos.X)
		s.Y = max32(s.Y, cs.Y+m.Vertical()+cb.pos.Y)
	}
	return s
}

// Arrange gives d its final rect, relative to its parent, and lays out its children.
func Arrange(d Node, r Rect) {
	b := d.box()
	if b.flags&(flagAnimateLayout|flagLaidOut) == flagAnimateLayout|flagLaidOut {
		b.offsetX.Shift(b.rect.X - r.X)
		b.offsetY.Shift(b.rect.Y - r.Y)
	}
	b.rect = r
	b.flags |= flagLaidOut
	content := b.ContentRect()
	if l, ok := d.(Layouter); ok {
		l.ArrangeChildren(content)
		return
	}
	for _, c := range b.children {
		ArrangeFree(c, content)
	}
}

// ArrangeFree places c inside area by its anchor, origin and position. It is
// the default for nodes that are not Layouters.
func ArrangeFree(c Node, area Rect) {
	cb := c.box()
	if !cb.visible() {
		return
	}
	m := marginPx(cb)
	avail := shrink(Vec2{area.W, area.H}, m)
	s := Measure(c, avail)
	x := area.X + m.Left + cb.anchor.X*avail.X - cb.origin.X*s.X + cb.pos.X
	y := area.Y + m.Top + cb.anchor.Y*avail.Y - cb.origin.Y*s.Y + cb.pos.Y
	Arrange(c, Rect{x, y, s.X, s.Y})
}

func padPx(b *Box) Insets    { return b.padding.Scale(Unit()) }
func marginPx(b *Box) Insets { return b.margin.Scale(Unit()) }

// radiusPx scales a corner radius given in units.
func radiusPx(r int16) int16 { return Round(float32(r) * Unit()) }

func shrink(v Vec2, m Insets) Vec2 {
	return Vec2{max32(v.X-m.Horizontal(), 0), max32(v.Y-m.Vertical(), 0)}
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
