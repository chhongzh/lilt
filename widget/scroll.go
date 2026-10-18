package widget

import (
	"osu-mc/anim"
	"osu-mc/gfx"
	"osu-mc/ui"
)

// Axes selects the directions a Scroll moves in.
type Axes uint8

const (
	ScrollVertical Axes = 1 << iota
	ScrollHorizontal
	ScrollBoth = ScrollVertical | ScrollHorizontal
)

const (
	scrollStep   = 16
	revealMargin = 2
	// flingFactor turns the last drag step into extra travel on release.
	flingFactor = 8
)

// Scroll clips its content and moves it by animating the content's offset.
// Arrow keys move focus through the content and scroll to reveal it; past
// the last focusable node they scroll in steps, and at an edge the content
// rubber-bands and springs back on release.
type Scroll struct {
	ui.Box
	content ui.Drawable
	axes    Axes

	pos   [2]float32 // scroll distance, 0 at the start
	pull  [2]float32 // overscroll from keys, before rubber-banding
	limit [2]float32
	view  [2]float32

	dragging  bool
	dragStart [2]int16
	dragFrom  [2]float32
	last      [2]int16
	step      [2]float32
}

func NewScroll(content ui.Drawable, axes Axes, opts ...ui.Option) *Scroll {
	s := &Scroll{}
	s.init(content, axes)
	ui.Apply(s, opts)
	return s
}

func (s *Scroll) init(content ui.Drawable, axes Axes) {
	s.content, s.axes = content, axes
	s.SetClip(true)
	b := ui.BoxOf(content)
	b.OffsetX().SetSpring(anim.Snappy)
	b.OffsetY().SetSpring(anim.Snappy)
	ui.Attach(s, content)
}

// Content is where Add puts children.
func (s *Scroll) Content() ui.Drawable { return s.content }

// ScrollPosition is the distance scrolled along x and y.
func (s *Scroll) ScrollPosition() ui.Vec2 { return ui.Vec2{X: s.pos[0], Y: s.pos[1]} }

func (s *Scroll) on(i int) bool { return s.axes&(ScrollHorizontal>>i) != 0 }

func (s *Scroll) offset(i int) *anim.Value {
	if i == 0 {
		return ui.BoxOf(s.content).OffsetX()
	}
	return ui.BoxOf(s.content).OffsetY()
}

func (s *Scroll) measureContent(max ui.Vec2) ui.Vec2 {
	if s.on(0) {
		max.X = ui.Unbounded
	}
	if s.on(1) {
		max.Y = ui.Unbounded
	}
	return ui.Measure(s.content, max)
}

func (s *Scroll) MeasureChildren(max ui.Vec2) ui.Vec2 {
	c := s.measureContent(max)
	return ui.Vec2{X: min(c.X, max.X), Y: min(c.Y, max.Y)}
}

func (s *Scroll) ArrangeChildren(r ui.Rect) {
	c := s.measureContent(ui.Vec2{X: r.W, Y: r.H})
	w, h := r.W, r.H
	if s.on(0) {
		w = max(c.X, r.W)
	}
	if s.on(1) {
		h = max(c.Y, r.H)
	}
	ui.Arrange(s.content, ui.Rect{X: r.X, Y: r.Y, W: w, H: h})
	s.view = [2]float32{r.W, r.H}
	s.limit = [2]float32{w - r.W, h - r.H}
	for i := range s.pos {
		s.pos[i] = clampf(s.pos[i], 0, s.limit[i])
	}
	s.apply()
}

func (s *Scroll) apply() {
	for i := range s.pos {
		if s.on(i) {
			s.offset(i).Set(-s.pos[i] + anim.RubberBand(s.pull[i], s.view[i]))
		}
	}
	ui.Invalidate()
}

// restRect is d's resting rect in the content's coordinates.
func (s *Scroll) restRect(d ui.Drawable) ui.Rect {
	b := ui.BoxOf(d)
	r := b.Rect()
	for n := b.Parent(); n != nil && n != s.content; n = ui.BoxOf(n).Parent() {
		nb := ui.BoxOf(n)
		r.X += nb.Rect().X + nb.OffsetX().Target()
		r.Y += nb.Rect().Y + nb.OffsetY().Target()
	}
	r.X += b.OffsetX().Target()
	r.Y += b.OffsetY().Target()
	return r
}

func span(r ui.Rect, i int) (lo, size float32) {
	if i == 0 {
		return r.X, r.W
	}
	return r.Y, r.H
}

func (s *Scroll) Reveal(d ui.Drawable) {
	if !ui.IsAncestor(s.content, d) {
		return
	}
	r := s.restRect(d)
	for i := range s.pos {
		if !s.on(i) {
			continue
		}
		lo, size := span(r, i)
		if lo < s.pos[i] {
			s.pos[i] = lo - revealMargin
		} else if lo+size > s.pos[i]+s.view[i] {
			s.pos[i] = lo + size - s.view[i] + revealMargin
		}
		s.pos[i] = clampf(s.pos[i], 0, s.limit[i])
		s.pull[i] = 0
	}
	s.apply()
}

func arrowAxis(k ui.Key) (i int, sign float32) {
	switch k {
	case ui.KeyLeft:
		return 0, -1
	case ui.KeyRight:
		return 0, 1
	case ui.KeyUp:
		return 1, -1
	}
	return 1, 1
}

func (s *Scroll) ArrowScroll(dir ui.Key, next ui.Drawable, phase ui.Phase) bool {
	i, sign := arrowAxis(dir)
	if !s.on(i) {
		return false
	}
	if phase == ui.Release {
		if s.pull[i] != 0 {
			s.pull[i] = 0
			s.apply()
		}
		return false
	}
	inside := next != nil && ui.IsAncestor(s.content, next)
	// A target more than half a screen away would skip unfocusable content
	// such as a long paragraph, so step through it instead.
	if inside {
		lo, size := span(s.restRect(next), i)
		near := lo+size > s.pos[i]-s.view[i]/2 && lo < s.pos[i]+s.view[i]*1.5
		if near {
			s.pull[i] = 0
			return false
		}
	}
	room := s.pos[i]
	if sign > 0 {
		room = s.limit[i] - s.pos[i]
	}
	if room > 0.5 {
		s.pos[i] += sign * min(scrollStep, room)
		s.apply()
		return true
	}
	if next != nil {
		return false
	}
	s.pull[i] -= sign * scrollStep
	s.apply()
	return true
}

func (s *Scroll) HandlePointer(e *ui.PointerEvent) bool {
	p := [2]int16{e.X, e.Y}
	switch e.Phase {
	case ui.PointerDown:
		s.dragging, s.dragStart, s.last = false, p, p
		s.dragFrom, s.step = s.pos, [2]float32{}
		return false
	case ui.PointerMove:
		if !s.dragging {
			dx, dy := absf(float32(p[0]-s.dragStart[0])), absf(float32(p[1]-s.dragStart[1]))
			if (dy >= dx && !s.on(1)) || (dx > dy && !s.on(0)) {
				return false
			}
			s.dragging = true
		}
		for i := range s.pos {
			if !s.on(i) {
				continue
			}
			raw := s.dragFrom[i] + float32(s.dragStart[i]-p[i])
			s.pos[i] = clampf(raw, 0, s.limit[i])
			s.offset(i).Jump(-s.pos[i] - anim.RubberBand(raw-s.pos[i], s.view[i]))
			s.step[i] = float32(s.last[i] - p[i])
		}
		s.last = p
		ui.Invalidate()
		return true
	case ui.PointerUp:
		if !s.dragging {
			return false
		}
		s.dragging = false
		for i := range s.pos {
			s.pos[i] = clampf(s.pos[i]+s.step[i]*flingFactor, 0, s.limit[i])
			s.pull[i] = 0
		}
		s.apply()
		return true
	}
	return false
}

// DrawOver paints thin scroll thumbs that shrink while overscrolled.
func (s *Scroll) DrawOver(c gfx.Canvas) {
	r := s.Rect()
	pad := s.Padding()
	col := ui.CurrentTheme.Foreground
	if s.on(1) && s.limit[1] > 0 {
		lo, size := thumb(-s.offset(1).Get(), s.limit[1], s.view[1])
		c.VLine(ui.Round(r.W-1), ui.Round(pad.Top+lo), ui.Round(size), col)
	}
	if s.on(0) && s.limit[0] > 0 {
		lo, size := thumb(-s.offset(0).Get(), s.limit[0], s.view[0])
		c.HLine(ui.Round(pad.Left+lo), ui.Round(r.H-1), ui.Round(size), col)
	}
}

func thumb(p, limit, view float32) (lo, size float32) {
	size = max(view*view/(view+limit), 4)
	var over float32
	if p < 0 {
		over = -p
	} else if p > limit {
		over = p - limit
	}
	size = max(size-over, 2)
	return clampf(p/limit, 0, 1) * (view - size), size
}

func clampf(x, lo, hi float32) float32 {
	return max(lo, min(x, hi))
}

func absf(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
