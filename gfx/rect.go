package gfx

// Rect is a pixel rectangle.
type Rect struct {
	X, Y, W, H int16
}

func (r Rect) Empty() bool { return r.W <= 0 || r.H <= 0 }

func (r Rect) Right() int16  { return r.X + r.W }
func (r Rect) Bottom() int16 { return r.Y + r.H }

func (r Rect) Contains(x, y int16) bool {
	return x >= r.X && y >= r.Y && x < r.X+r.W && y < r.Y+r.H
}

func (r Rect) Offset(dx, dy int16) Rect { return Rect{r.X + dx, r.Y + dy, r.W, r.H} }

// Inset shrinks r by d on every side; negative d grows it.
func (r Rect) Inset(d int16) Rect { return Rect{r.X + d, r.Y + d, r.W - 2*d, r.H - 2*d} }

func (r Rect) Intersect(o Rect) Rect {
	x0, y0 := max(r.X, o.X), max(r.Y, o.Y)
	x1, y1 := min(r.Right(), o.Right()), min(r.Bottom(), o.Bottom())
	if x1 <= x0 || y1 <= y0 {
		return Rect{}
	}
	return Rect{x0, y0, x1 - x0, y1 - y0}
}

func (r Rect) Overlaps(o Rect) bool { return !r.Intersect(o).Empty() }
