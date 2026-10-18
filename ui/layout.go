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
		if avail >= unbounded {
			return 0, true
		}
		return s.v * avail, false
	}
	return 0, true
}

// Measure returns the size d wants inside max, the space left after its margin.
func Measure(d Drawable, max Vec2) Vec2 {
	b := d.box()
	w, autoW := resolve(b.width, max.X)
	h, autoH := resolve(b.height, max.Y)
	if autoW || autoH {
		inner := Vec2{max.X, max.Y}
		if !autoW {
			inner.X = w
		}
		if !autoH {
			inner.Y = h
		}
		inner.X -= b.padding.Horizontal()
		inner.Y -= b.padding.Vertical()
		c := contentSize(d, inner)
		if autoW {
			w = c.X + b.padding.Horizontal()
		}
		if autoH {
			h = c.Y + b.padding.Vertical()
		}
	}
	return Vec2{w, h}
}

func contentSize(d Drawable, max Vec2) Vec2 {
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
		cs := Measure(c, shrink(max, cb.margin))
		s.X = max32(s.X, cs.X+cb.margin.Horizontal()+cb.pos.X)
		s.Y = max32(s.Y, cs.Y+cb.margin.Vertical()+cb.pos.Y)
	}
	return s
}

// Arrange gives d its final rect, relative to its parent, and lays out its children.
func Arrange(d Drawable, r Rect) {
	b := d.box()
	if b.flags&(flagAnimateLayout|flagLaidOut) == flagAnimateLayout|flagLaidOut {
		b.offsetX.Shift(b.rect.X - r.X)
		b.offsetY.Shift(b.rect.Y - r.Y)
	}
	b.rect = r
	b.flags |= flagLaidOut
	content := Rect{b.padding.Left, b.padding.Top, r.W - b.padding.Horizontal(), r.H - b.padding.Vertical()}
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
func ArrangeFree(c Drawable, area Rect) {
	cb := c.box()
	if !cb.visible() {
		return
	}
	avail := shrink(Vec2{area.W, area.H}, cb.margin)
	s := Measure(c, avail)
	x := area.X + cb.margin.Left + cb.anchor.X*avail.X - cb.origin.X*s.X + cb.pos.X
	y := area.Y + cb.margin.Top + cb.anchor.Y*avail.Y - cb.origin.Y*s.Y + cb.pos.Y
	Arrange(c, Rect{x, y, s.X, s.Y})
}

func shrink(v Vec2, m Insets) Vec2 {
	return Vec2{max32(v.X-m.Horizontal(), 0), max32(v.Y-m.Vertical(), 0)}
}

func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
