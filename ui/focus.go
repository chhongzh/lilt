package ui

// ArrowScroller is offered arrow keys by the focused node's ancestors before
// focus moves. next is the node focus would move to, or nil at an edge or on
// Release. Returning true consumes the key, which is how a scroll view
// scrolls through content with nothing focusable or pushes back at its end.
type ArrowScroller interface {
	ArrowScroll(dir Key, next Drawable, phase Phase) bool
}

// Revealer is told when a descendant gains focus so it can scroll it into view.
type Revealer interface{ Reveal(d Drawable) }

// FocusManager tracks the single focused node.
type FocusManager struct {
	current Drawable
}

func (f *FocusManager) Current() Drawable { return f.current }

// Set focuses d, or clears focus for nil. Hidden or unfocusable nodes are ignored.
func (f *FocusManager) Set(d Drawable) {
	if d != nil && (!d.box().Focusable() || !shown(d)) {
		return
	}
	if f.current == d {
		return
	}
	old := f.current
	f.current = d
	if fo, ok := old.(Focuser); ok {
		fo.FocusChanged(false)
	}
	if d != nil {
		if fo, ok := d.(Focuser); ok {
			fo.FocusChanged(true)
		}
		reveal(d)
	}
	Invalidate()
}

func reveal(d Drawable) {
	for n := d.box().parent; n != nil; n = n.box().parent {
		if r, ok := n.(Revealer); ok {
			r.Reveal(d)
		}
	}
}

// leaving moves focus out of b before b is removed or hidden.
func (f *FocusManager) leaving(b *Box) {
	if f.current == nil || !IsAncestor(b, f.current) {
		return
	}
	var next Drawable
	if scope := app.Screens.scope(); scope != nil {
		fd := finder{dir: KeyNone, src: restRect(f.current), from: f.current.box(), exclude: b}
		fd.walk(scope)
		next = fd.best
	}
	f.Set(nil)
	f.Set(next)
}

func shown(d Drawable) bool {
	for n := d; n != nil; n = n.box().parent {
		if !n.box().visible() {
			return false
		}
	}
	return true
}

// FindFocus picks the node focus should move to from `from` in direction dir,
// following Android's FocusFinder: candidates overlapping the source's beam
// win, then the lowest 13·major² + minor² distance.
func FindFocus(scope, from Drawable, dir Key) Drawable {
	if from == nil {
		return FirstFocusable(scope)
	}
	f := finder{dir: dir, src: restRect(from), from: from.box()}
	f.walk(scope)
	return f.best
}

// FirstFocusable returns the first focusable node in tree order.
func FirstFocusable(d Drawable) Drawable {
	b := d.box()
	if !b.visible() {
		return nil
	}
	if b.Focusable() {
		return d
	}
	for _, c := range b.children {
		if f := FirstFocusable(c); f != nil {
			return f
		}
	}
	return nil
}

// restOrigin sums layout positions and offset targets, so focus decisions
// use where things are going rather than where an animation has them now.
func restOrigin(d Drawable) (x, y float32) {
	for n := d.box().parent; n != nil; n = n.box().parent {
		b := n.box()
		x += b.rect.X + b.offsetX.Target()
		y += b.rect.Y + b.offsetY.Target()
	}
	return
}

func restRect(d Drawable) Rect {
	x, y := restOrigin(d)
	b := d.box()
	return Rect{x + b.rect.X + b.offsetX.Target(), y + b.rect.Y + b.offsetY.Target(), b.rect.W, b.rect.H}
}

type finder struct {
	dir      Key
	src      Rect
	from     *Box
	exclude  *Box
	best     Drawable
	bestRect Rect
}

func (f *finder) walk(scope Drawable) {
	x, y := restOrigin(scope)
	f.visit(scope, x, y)
}

func (f *finder) visit(d Drawable, ox, oy float32) {
	b := d.box()
	if !b.visible() || b == f.exclude {
		return
	}
	x, y := ox+b.rect.X+b.offsetX.Target(), oy+b.rect.Y+b.offsetY.Target()
	if b.Focusable() && b != f.from {
		r := Rect{x, y, b.rect.W, b.rect.H}
		if f.best == nil && f.candidate(r) || f.best != nil && f.better(r, f.bestRect) {
			f.best, f.bestRect = d, r
		}
	}
	for _, c := range b.children {
		f.visit(c, x, y)
	}
}

func (f *finder) candidate(r Rect) bool {
	s := f.src
	switch f.dir {
	case KeyRight:
		return (s.X < r.X || s.Right() <= r.X) && s.Right() < r.Right()
	case KeyLeft:
		return (s.Right() > r.Right() || s.X >= r.Right()) && s.X > r.X
	case KeyDown:
		return (s.Y < r.Y || s.Bottom() <= r.Y) && s.Bottom() < r.Bottom()
	case KeyUp:
		return (s.Bottom() > r.Bottom() || s.Y >= r.Bottom()) && s.Y > r.Y
	}
	return true
}

func (f *finder) better(a, b Rect) bool {
	if !f.candidate(a) {
		return false
	}
	if f.dir == KeyNone {
		return centerDist(f.src, a) < centerDist(f.src, b)
	}
	if f.beamBeats(a, b) {
		return true
	}
	if f.beamBeats(b, a) {
		return false
	}
	return f.weight(a) < f.weight(b)
}

func (f *finder) inBeam(r Rect) bool {
	s := f.src
	if f.dir == KeyLeft || f.dir == KeyRight {
		return r.Bottom() > s.Y && r.Y < s.Bottom()
	}
	return r.Right() > s.X && r.X < s.Right()
}

func (f *finder) beamBeats(a, b Rect) bool {
	if !f.inBeam(a) || f.inBeam(b) {
		return false
	}
	if !f.beyond(b) || f.dir == KeyLeft || f.dir == KeyRight {
		return true
	}
	return f.major(a) < f.majorFar(b)
}

// beyond reports whether r lies entirely past the source in the direction.
func (f *finder) beyond(r Rect) bool {
	s := f.src
	switch f.dir {
	case KeyRight:
		return s.Right() <= r.X
	case KeyLeft:
		return s.X >= r.Right()
	case KeyDown:
		return s.Bottom() <= r.Y
	}
	return s.Y >= r.Bottom()
}

func (f *finder) major(r Rect) float32 {
	s := f.src
	var d float32
	switch f.dir {
	case KeyRight:
		d = r.X - s.Right()
	case KeyLeft:
		d = s.X - r.Right()
	case KeyDown:
		d = r.Y - s.Bottom()
	case KeyUp:
		d = s.Y - r.Bottom()
	}
	return max32(d, 0)
}

func (f *finder) majorFar(r Rect) float32 {
	s := f.src
	var d float32
	switch f.dir {
	case KeyRight:
		d = r.Right() - s.Right()
	case KeyLeft:
		d = s.X - r.X
	case KeyDown:
		d = r.Bottom() - s.Bottom()
	case KeyUp:
		d = s.Y - r.Y
	}
	return max32(d, 1)
}

func (f *finder) weight(r Rect) float32 {
	major := f.major(r)
	sc, rc := f.src.Center(), r.Center()
	minor := absf(sc.X - rc.X)
	if f.dir == KeyLeft || f.dir == KeyRight {
		minor = absf(sc.Y - rc.Y)
	}
	return 13*major*major + minor*minor
}

func centerDist(a, b Rect) float32 {
	ac, bc := a.Center(), b.Center()
	dx, dy := ac.X-bc.X, ac.Y-bc.Y
	return dx*dx + dy*dy
}

func (a *App) navigate(e KeyEvent) {
	cur := a.Focus.current
	if e.Phase == Release {
		if cur == nil {
			return
		}
		for n := cur.box().parent; n != nil; n = n.box().parent {
			if s, ok := n.(ArrowScroller); ok {
				s.ArrowScroll(e.Key, nil, Release)
			}
		}
		return
	}
	scope := a.Screens.scope()
	if scope == nil {
		return
	}
	if cur == nil {
		a.Focus.Set(FirstFocusable(scope))
		return
	}
	next := FindFocus(scope, cur, e.Key)
	for n := cur.box().parent; n != nil; n = n.box().parent {
		if s, ok := n.(ArrowScroller); ok && s.ArrowScroll(e.Key, next, e.Phase) {
			return
		}
	}
	if next != nil {
		a.Focus.Set(next)
	}
}
