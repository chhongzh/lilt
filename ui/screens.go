package ui

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
)

// Popup is where a layer opened with Open sits and which way it moves.
type Popup uint8

const (
	// DrawerBottom rises from the bottom edge, as tall as its content.
	DrawerBottom Popup = iota
	DrawerTop
	// DrawerLeft slides in from the left edge, as wide as its content.
	DrawerLeft
	DrawerRight
	// Dialog drops into the middle of the screen.
	Dialog
	// Toast rises near the bottom and takes neither focus nor keys.
	Toast
	// fullPage is a screen shown with Push.
	fullPage
)

const (
	// popupMargin keeps dialogs and toasts off the screen edges.
	popupMargin = 4
	// backPeek is how far Back held down moves the top layer; one that
	// cannot leave moves half as far.
	backPeek = 8
	// backNudge is the bounce when Back is released on the root page.
	backNudge = 6
)

// layer wraps each page or popup so transitions move it as a unit.
type layer struct {
	Box
	content Node
	kind    Popup
	focus   Node // restored when the layer above closes
	placed  bool // the entrance, which needs a laid-out content, has begun
	exiting bool
}

func (l *layer) modal() bool { return l.kind != Toast }

func (l *layer) MeasureChildren(max Vec2) Vec2 { return max }

func (l *layer) ArrangeChildren(r Rect) {
	if l.kind == fullPage {
		Arrange(l.content, r)
		return
	}
	Arrange(l.content, l.place(r))
	if !l.placed {
		l.placed = true
		dx, dy := l.away()
		l.offsetX.Jump(dx)
		l.offsetY.Jump(dy)
		l.offsetX.Set(0)
		l.offsetY.Set(0)
	}
}

// place sizes a popup's content and pins it to its edge or the middle.
func (l *layer) place(r Rect) Rect {
	switch l.kind {
	case DrawerBottom, DrawerTop:
		h := min(Measure(l.content, Vec2{r.W, r.H}).Y, r.H)
		if l.kind == DrawerTop {
			return Rect{r.X, r.Y, r.W, h}
		}
		return Rect{r.X, r.Bottom() - h, r.W, h}
	case DrawerLeft, DrawerRight:
		w := min(Measure(l.content, Vec2{r.W, r.H}).X, r.W)
		if l.kind == DrawerLeft {
			return Rect{r.X, r.Y, w, r.H}
		}
		return Rect{r.Right() - w, r.Y, w, r.H}
	}
	const m = popupMargin
	room := Vec2{r.W - 2*m, r.H - 2*m}
	s := Measure(l.content, room)
	w, h := min(s.X, room.X), min(s.Y, room.Y)
	x := r.X + (r.W-w)/2
	if l.kind == Toast {
		return Rect{x, r.Bottom() - m - h, w, h}
	}
	return Rect{x, r.Y + (r.H-h)/2, w, h}
}

// away is the offset that carries the layer just off screen, the way it
// came in and leaves.
func (l *layer) away() (dx, dy float32) {
	c := l.content.box().rect
	w, h := float32(scrW), float32(scrH)
	switch l.kind {
	case fullPage:
		return w, 0
	case DrawerBottom, Toast:
		return 0, h - c.Y + 1
	case DrawerTop:
		return 0, -c.Bottom() - 1
	case DrawerLeft:
		return -c.Right() - 1, 0
	case DrawerRight:
		return w - c.X + 1, 0
	}
	return 0, -c.Bottom() - 1
}

// peek moves the layer d toward where it would leave and returns the shift.
func (l *layer) peek(d float32) (dx, dy float32) {
	ax, ay := l.away()
	dx, dy = sign(ax)*d, sign(ay)*d
	l.offsetX.Set(dx)
	l.offsetY.Set(dy)
	return dx, dy
}

func (l *layer) settled() bool { return l.offsetX.Settled() && l.offsetY.Settled() }

func (l *layer) Draw(c gfx.Canvas) {
	t := CurrentTheme
	if l.kind == fullPage {
		c.FillRect(gfx.Rect{W: Round(l.rect.W), H: Round(l.rect.H)}, t.Background)
		return
	}
	r := l.content.box().rect.Pixels()
	rad := t.Radius
	if l.kind <= DrawerRight {
		// A drawer's frame runs past its edge, so only the inner corners show.
		rad *= 2
		ext := rad + 1
		switch l.kind {
		case DrawerBottom:
			r.H += ext
		case DrawerTop:
			r.Y, r.H = r.Y-ext, r.H+ext
		case DrawerLeft:
			r.X, r.W = r.X-ext, r.W+ext
		case DrawerRight:
			r.W += ext
		}
	}
	c.FillRoundRect(r, rad, t.Background)
	c.StrokeRoundRect(r, rad, t.Foreground)
}

func sign(x float32) float32 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	}
	return 0
}

// layerStack holds pages and popups in drawing order. Only the top modal
// layer receives focus and pointer input; pages fully covered by other
// pages are hidden.
type layerStack struct {
	list         []*layer
	focusPending bool
}

func (s *layerStack) add(n Node, kind Popup) {
	l := &layer{content: n, kind: kind}
	l.width, l.height = Rel(1), Rel(1)
	l.offsetX.SetSpring(anim.Smooth)
	l.offsetY.SetSpring(anim.Smooth)
	if kind == Dialog {
		l.offsetY.SetSpring(anim.Default)
	}
	Attach(l, n)
	if kind == fullPage {
		l.placed = true
		if p := s.pageBelow(nil); p != nil {
			p.offsetX.Set(-float32(scrW) / 3)
			l.offsetX.Jump(float32(scrW))
			l.offsetX.Set(0)
		}
	}
	if kind != Toast {
		if top := s.top(); top != nil {
			top.focus = focus.current
			cursor.Leap()
		}
		focus.Set(nil)
		s.focusPending = true
	}
	s.list = append(s.list, l)
	Attach(root, l)
}

// close sends l away. The root page cannot close.
func (s *layerStack) close(l *layer) bool {
	if l == nil || l.exiting {
		return false
	}
	below := s.pageBelow(l)
	if l.kind == fullPage && below == nil {
		return false
	}
	wasTop := l == s.top()
	l.exiting = true
	dx, dy := l.away()
	l.offsetX.Set(dx)
	l.offsetY.Set(dy)
	if l.kind == fullPage {
		below.offsetX.Set(0)
		below.SetVisible(true)
	}
	if wasTop {
		cursor.Leap()
		focus.Set(nil)
		if t := s.top(); t != nil {
			t.SetVisible(true)
			focus.Set(t.focus)
		}
		s.focusPending = true
	}
	return true
}

func (s *layerStack) find(n Node) *layer {
	for i := len(s.list) - 1; i >= 0; i-- {
		if l := s.list[i]; !l.exiting && l.content.box() == n.box() {
			return l
		}
	}
	return nil
}

// top is the topmost modal layer that is not leaving.
func (s *layerStack) top() *layer {
	for i := len(s.list) - 1; i >= 0; i-- {
		if l := s.list[i]; !l.exiting && l.modal() {
			return l
		}
	}
	return nil
}

// pageBelow is the topmost page under l, or under everything for nil.
func (s *layerStack) pageBelow(l *layer) *layer {
	found := l == nil
	for i := len(s.list) - 1; i >= 0; i-- {
		switch c := s.list[i]; {
		case c == l:
			found = true
		case found && !c.exiting && c.kind == fullPage:
			return c
		}
	}
	return nil
}

// back handles a Back key nothing else claimed. Pressing peeks the top
// layer aside as if starting to dismiss it, and the cursor leans with it;
// releasing dismisses it. The root page bounces back instead.
func (s *layerStack) back(e *KeyEvent) {
	top := s.top()
	if top == nil {
		return
	}
	d := float32(backPeek)
	if top.kind == fullPage && s.pageBelow(top) == nil {
		d /= 2
	}
	switch e.Phase {
	case Press:
		cursor.Lean(top.peek(d))
	case Release:
		cursor.Lean(0, 0)
		top.peek(0)
		if e.Canceled || s.close(top) {
			return
		}
		top.offsetX.Bump(backNudge)
		cursor.Nudge(backNudge)
	}
}

// scope is the subtree that focus and pointer input are confined to.
func (s *layerStack) scope() Node {
	if t := s.top(); t != nil {
		return t
	}
	return nil
}

func (s *layerStack) update() {
	for i := 0; i < len(s.list); {
		l := s.list[i]
		if l.exiting && l.settled() {
			copy(s.list[i:], s.list[i+1:])
			s.list[len(s.list)-1] = nil
			s.list = s.list[:len(s.list)-1]
			Remove(l)
			continue
		}
		i++
	}
	// A page held aside by Back or parked a third to the left covers only
	// part of the width, so pages are stacked until their resting spans
	// close over the whole screen.
	cover := -1
	w := float32(scrW)
	var lo, hi float32
	for i := len(s.list) - 1; i >= 0; i-- {
		l := s.list[i]
		if l.exiting || l.kind != fullPage || !l.settled() || l.offsetY.Target() != 0 {
			continue
		}
		x := l.offsetX.Target()
		switch {
		case lo == hi:
			lo, hi = x, x+w
		case x <= hi && x+w >= lo:
			lo, hi = min(lo, x), max(hi, x+w)
		}
		if lo <= 0 && hi >= w {
			cover = i
			break
		}
	}
	for i, l := range s.list {
		l.SetVisible(i >= cover)
	}
}

// focusFirst runs after the first layout of a new layer, so revealing the
// focused node sees real rects. A node focused right after opening is kept
// and revealed again.
func (s *layerStack) focusFirst() {
	if !s.focusPending {
		return
	}
	s.focusPending = false
	if cur := focus.current; cur != nil {
		reveal(cur)
		return
	}
	if sc := s.scope(); sc != nil {
		focus.Set(FirstFocusable(sc))
	}
}

// Push shows p as a new full screen, sliding in from the right.
func Push(p Node) { layers.add(p, fullPage) }

// Pop dismisses the top layer, page or popup. The root page stays.
func Pop() bool { return layers.close(layers.top()) }

// Open shows n as a popup. Drawers and dialogs hold focus and keys until
// they close; Back closes them.
func Open(n Node, p Popup) { layers.add(n, p) }

// Close dismisses the page or popup showing n.
func Close(n Node) bool { return layers.close(layers.find(n)) }
