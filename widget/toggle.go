package widget

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

const markSize = 9

// mark is the box or circle a CheckBox or Radio draws; on springs 0 to 1.
type mark struct {
	ui.Box
	radio bool
	on    anim.Value
}

func (m *mark) init(radio bool) {
	m.radio = radio
	m.on.SetSpring(anim.Bouncy)
	m.on.SetPrecision(0.01)
}

func (m *mark) Leaf() {}

func (m *mark) Measure(ui.Vec2) ui.Vec2 {
	s := float32(ui.Round(markSize * ui.Unit()))
	return ui.Vec2{X: s, Y: s}
}

// turn is not called set: TinyGo 0.42 and 0.43 crash compiling a type that
// embeds ui.Box and declares a method named like one of Box's unexported ones.
func (m *mark) turn(on bool) {
	if on {
		m.on.Set(1)
	} else {
		m.on.Set(0)
	}
	ui.Invalidate()
}

func (m *mark) Draw(c gfx.Canvas) {
	col := ui.CurrentTheme.Foreground
	size := m.Rect().Pixels().W
	h := size / 2
	fill := ui.Round(clampf(m.on.Get(), 0, 1.3) * float32(size) * 0.29)
	if m.radio {
		c.StrokeCircle(h, h, h, col)
		if fill > 0 {
			c.FillCircle(h, h, fill-1, col)
		}
		return
	}
	c.StrokeRoundRect(gfx.Rect{W: size, H: size}, 2, col)
	if fill > 0 {
		c.FillRect(gfx.Rect{X: h - fill + 1, Y: h - fill + 1, W: 2*fill - 1, H: 2*fill - 1}, col)
	}
}

// CheckBox is a list row that toggles a box on release. The box flips at
// once; an OnChange handler that fails flips it back.
type CheckBox struct {
	ListItem
	mark    mark
	checked bool
}

func NewCheckBox(title string, checked bool, opts ...ui.Option) *CheckBox {
	cb := &CheckBox{checked: checked}
	cb.init(title)
	cb.kind = false
	cb.mark.init(false)
	if checked {
		cb.mark.on.Jump(1)
	}
	cb.SetTrailing(&cb.mark)
	cb.activate = cb.toggle
	ui.Apply(cb, opts)
	return cb
}

func (cb *CheckBox) Checked() bool { return cb.checked }

// SetChecked changes the box without calling OnChange.
func (cb *CheckBox) SetChecked(v bool) {
	cb.checked = v
	cb.mark.turn(v)
}

func (cb *CheckBox) toggle() {
	from := cb.checked
	cb.SetChecked(!from)
	cb.run(changeCall(cb.change, from, !from), func(err error) {
		if err != nil {
			cb.SetChecked(from)
			shake(cb)
		}
	})
}

// RadioGroup keeps one of its Radios selected. Its OnChange handler gets
// the old and new index; one that fails puts the old selection back.
type RadioGroup struct {
	Control
	radios   []*Radio
	selected int
	// picked runs on the UI goroutine after a radio is chosen, for widgets
	// such as Select that build a group of their own.
	picked func(from, to int)
}

// NewRadioGroup starts with selected chosen; use -1 for none.
func NewRadioGroup(selected int) *RadioGroup {
	g := &RadioGroup{selected: selected}
	g.kind = 0
	return g
}

func (g *RadioGroup) Selected() int { return g.selected }

func (g *RadioGroup) Radios() []*Radio { return g.radios }

// Select chooses radio i without calling OnChange.
func (g *RadioGroup) Select(i int) {
	g.selected = i
	for j, r := range g.radios {
		r.mark.turn(j == i)
	}
}

// Radio is a list row belonging to a RadioGroup. Its OnClick handler, if
// any, runs after the group's OnChange.
type Radio struct {
	ListItem
	mark  mark
	group *RadioGroup
	index int
}

func NewRadio(g *RadioGroup, title string, opts ...ui.Option) *Radio {
	r := &Radio{group: g, index: len(g.radios)}
	g.radios = append(g.radios, r)
	r.init(title)
	r.mark.init(true)
	if r.index == g.selected {
		r.mark.on.Jump(1)
	}
	r.SetTrailing(&r.mark)
	r.activate = r.choose
	ui.Apply(r, opts)
	return r
}

// choose shows the radio's own loading while the group's handler runs, and
// holds off the rest of the group meanwhile.
func (r *Radio) choose() {
	g := r.group
	if g.busy {
		return
	}
	from, to := g.selected, r.index
	g.Select(to)
	if g.picked != nil {
		g.picked(from, to)
	}
	change, click := changeCall(g.change, from, to), clickCall(r.click)
	call := change
	if click != nil {
		call = func() error {
			if change != nil {
				if err := change(); err != nil {
					return err
				}
			}
			return click()
		}
	}
	g.busy = call != nil
	r.run(call, func(err error) {
		g.busy = false
		if err != nil {
			g.Select(from)
			shake(r)
		}
	})
}
