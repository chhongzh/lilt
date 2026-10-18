package ui

import (
	"osu-mc/anim"
	"osu-mc/gfx"
)

// Transition is how a screen enters and leaves.
type Transition uint8

const (
	// SlideLeft pushes the new screen in from the right while the old one
	// drifts a third of the way left, as in a navigation stack.
	SlideLeft Transition = iota
	// Sheet rises from the bottom over the previous screen, sized to its content.
	Sheet
)

// screenHost wraps each pushed screen so transitions can move it as a unit.
type screenHost struct {
	Box
	content    Drawable
	transition Transition
	focus      Drawable // restored when the screen above is popped
	exiting    bool
}

func (h *screenHost) MeasureChildren(max Vec2) Vec2 { return max }

func (h *screenHost) ArrangeChildren(r Rect) {
	if h.transition != Sheet {
		Arrange(h.content, r)
		return
	}
	s := Measure(h.content, Vec2{r.W, r.H})
	Arrange(h.content, Rect{r.X, r.Bottom() - s.Y, r.W, s.Y})
}

func (h *screenHost) Draw(c gfx.Canvas) {
	t := CurrentTheme
	if h.transition != Sheet {
		c.FillRect(gfx.Rect{W: Round(h.rect.W), H: Round(h.rect.H)}, t.Background)
		return
	}
	// The sheet's lower corners hang off screen so only the top is rounded.
	r := h.content.box().rect.Pixels()
	r.H += t.Radius + 1
	c.FillRoundRect(r, t.Radius*2, t.Background)
	c.StrokeRoundRect(r, t.Radius*2, t.Foreground)
}

func (h *screenHost) opaque() bool { return h.transition != Sheet }

func (h *screenHost) settled() bool { return h.offsetX.Settled() && h.offsetY.Settled() }

// ScreenStack is the navigation stack. Only the top screen receives focus
// and pointer input; screens fully covered by an opaque screen are hidden.
type ScreenStack struct {
	hosts        []*screenHost
	focusPending bool
}

// Push shows d above the current screen.
func (s *ScreenStack) Push(d Drawable, t Transition) {
	w, h := float32(app.w), float32(app.h)
	host := &screenHost{content: d, transition: t}
	host.width, host.height = Rel(1), Rel(1)
	host.offsetX.SetSpring(anim.Smooth)
	host.offsetY.SetSpring(anim.Smooth)
	Attach(host, d)
	if top := s.top(); top != nil {
		top.focus = app.Focus.current
		if t == SlideLeft {
			top.offsetX.Set(-w / 3)
			host.offsetX.Jump(w)
		} else {
			host.offsetY.Jump(h)
		}
		host.offsetX.Set(0)
		host.offsetY.Set(0)
	}
	s.hosts = append(s.hosts, host)
	Attach(app.root, host)
	app.Focus.Set(nil)
	s.focusPending = true
}

// Pop dismisses the top screen. The root screen cannot be popped.
func (s *ScreenStack) Pop() bool {
	top := s.top()
	prev := s.below(top)
	if prev == nil {
		return false
	}
	top.exiting = true
	if top.transition == SlideLeft {
		top.offsetX.Set(float32(app.w))
		prev.offsetX.Set(0)
	} else {
		top.offsetY.Set(top.content.box().rect.H + 2)
	}
	prev.SetVisible(true)
	app.Focus.Set(nil)
	app.Focus.Set(prev.focus)
	s.focusPending = true
	return true
}

// Top returns the content of the top screen, or nil.
func (s *ScreenStack) Top() Drawable {
	if t := s.top(); t != nil {
		return t.content
	}
	return nil
}

// Len counts screens, not including ones still animating out.
func (s *ScreenStack) Len() int {
	n := 0
	for _, h := range s.hosts {
		if !h.exiting {
			n++
		}
	}
	return n
}

func (s *ScreenStack) top() *screenHost {
	for i := len(s.hosts) - 1; i >= 0; i-- {
		if !s.hosts[i].exiting {
			return s.hosts[i]
		}
	}
	return nil
}

func (s *ScreenStack) below(h *screenHost) *screenHost {
	if h == nil {
		return nil
	}
	found := false
	for i := len(s.hosts) - 1; i >= 0; i-- {
		switch {
		case s.hosts[i] == h:
			found = true
		case found && !s.hosts[i].exiting:
			return s.hosts[i]
		}
	}
	return nil
}

// scope is the subtree that focus and pointer input are confined to.
func (s *ScreenStack) scope() Drawable {
	if t := s.top(); t != nil {
		return t
	}
	return nil
}

func (s *ScreenStack) update() {
	for i := 0; i < len(s.hosts); {
		h := s.hosts[i]
		if h.exiting && h.settled() {
			copy(s.hosts[i:], s.hosts[i+1:])
			s.hosts[len(s.hosts)-1] = nil
			s.hosts = s.hosts[:len(s.hosts)-1]
			Remove(h)
			continue
		}
		i++
	}
	cover := -1
	for i := len(s.hosts) - 1; i >= 0; i-- {
		h := s.hosts[i]
		if !h.exiting && h.opaque() && h.settled() {
			cover = i
			break
		}
	}
	for i, h := range s.hosts {
		h.SetVisible(i >= cover)
	}
}

// focusFirst runs after the first layout of a new screen, so revealing the
// focused node sees real rects. A node focused right after Push is kept and
// revealed again.
func (s *ScreenStack) focusFirst() {
	if !s.focusPending {
		return
	}
	s.focusPending = false
	if cur := app.Focus.current; cur != nil {
		reveal(cur)
		return
	}
	if sc := s.scope(); sc != nil {
		app.Focus.Set(FirstFocusable(sc))
	}
}
