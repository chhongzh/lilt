package ui

// Track is a grid column width: fixed pixels or a share of what remains.
type Track struct{ px, fr float32 }

func Fixed(px float32) Track { return Track{px: px} }
func Fr(n float32) Track     { return Track{fr: n} }

// Grid places children row by row into columns. Each row is as tall as its
// tallest child, and children with Auto sizes fill their cell.
type Grid struct {
	Box
	cols   []Track
	gap    float32
	widths []float32
}

func NewGrid(opts ...Option) *Grid {
	g := &Grid{}
	Apply(g, opts)
	if len(g.cols) == 0 {
		g.SetColumns(Fr(1))
	}
	return g
}

func (g *Grid) SetColumns(cols ...Track) {
	g.cols = cols
	g.widths = make([]float32, len(cols))
	Relayout()
}

func (g *Grid) SetGap(gap float32) {
	if g.gap != gap {
		g.gap = gap
		Relayout()
	}
}

func (g *Grid) columnWidths(avail float32) {
	var fixed, fr float32
	for _, t := range g.cols {
		fixed += t.px
		fr += t.fr
	}
	rest := max32(avail-fixed-g.gap*float32(len(g.cols)-1), 0)
	for i, t := range g.cols {
		g.widths[i] = t.px
		if fr > 0 {
			g.widths[i] += rest * t.fr / fr
		}
	}
}

// rows walks visible children a row at a time, measuring each and calling
// place once the row's height is known.
func (g *Grid) rows(maxH float32, place func(row []Drawable, y, h float32)) float32 {
	var y float32
	start, col := 0, 0
	var h float32
	flush := func(end int) {
		if place != nil {
			place(g.children[start:end], y, h)
		}
		y += h + g.gap
		h = 0
		col = 0
	}
	for i, c := range g.children {
		cb := c.box()
		if !cb.visible() {
			continue
		}
		cb.measured = Measure(c, shrink(Vec2{g.widths[col], maxH}, cb.margin))
		h = max32(h, cb.measured.Y+cb.margin.Vertical())
		col++
		if col == len(g.cols) {
			flush(i + 1)
			start = i + 1
		}
	}
	if col > 0 {
		flush(len(g.children))
	}
	return max32(y-g.gap, 0)
}

func (g *Grid) MeasureChildren(max Vec2) Vec2 {
	g.columnWidths(max.X)
	w := g.gap * float32(len(g.cols)-1)
	for _, cw := range g.widths {
		w += cw
	}
	return Vec2{w, g.rows(max.Y, nil)}
}

func (g *Grid) ArrangeChildren(content Rect) {
	g.columnWidths(content.W)
	g.rows(content.H, func(row []Drawable, y, h float32) {
		x := content.X
		col := 0
		for _, c := range row {
			cb := c.box()
			if !cb.visible() {
				continue
			}
			w := g.widths[col]
			s := cb.measured
			if cb.width.mode == autoSize {
				s.X = w - cb.margin.Horizontal()
			}
			if cb.height.mode == autoSize {
				s.Y = h - cb.margin.Vertical()
			}
			Arrange(c, Rect{x + cb.margin.Left, content.Y + y + cb.margin.Top, s.X, s.Y})
			x += w + g.gap
			col++
		}
	})
}
