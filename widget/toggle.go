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

func (m *mark) Measure(ui.Vec2) ui.Vec2 { return ui.Vec2{X: markSize, Y: markSize} }

// turn is not called set: TinyGo 0.42 crashes compiling a type that embeds
// ui.Box and declares a method named like one of Box's unexported ones.
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
	h := int16(markSize / 2)
	fill := ui.Round(clampf(m.on.Get(), 0, 1.3) * 2.6)
	if m.radio {
		c.StrokeCircle(h, h, h, col)
		if fill > 0 {
			c.FillCircle(h, h, fill-1, col)
		}
		return
	}
	c.StrokeRoundRect(gfx.Rect{W: markSize, H: markSize}, 2, col)
	if fill > 0 {
		c.FillRect(gfx.Rect{X: h - fill + 1, Y: h - fill + 1, W: 2*fill - 1, H: 2*fill - 1}, col)
	}
}

// CheckBox is a list row that toggles a box on release.
type CheckBox struct {
	ListItem
	mark     mark
	checked  bool
	onChange func(bool)
}

func NewCheckBox(title string, checked bool, opts ...ui.Option) *CheckBox {
	cb := &CheckBox{checked: checked}
	cb.init(title)
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

func (cb *CheckBox) SetChecked(v bool) {
	cb.checked = v
	cb.mark.turn(v)
}

func (cb *CheckBox) SetOnChange(fn func(bool)) { cb.onChange = fn }

func (cb *CheckBox) toggle() {
	cb.SetChecked(!cb.checked)
	if cb.onChange != nil {
		cb.onChange(cb.checked)
	}
}

// RadioGroup keeps one of its Radios selected.
type RadioGroup struct {
	radios   []*Radio
	selected int
	onChange func(int)
}

// NewRadioGroup starts with selected chosen; use -1 for none.
func NewRadioGroup(selected int, onChange func(int)) *RadioGroup {
	return &RadioGroup{selected: selected, onChange: onChange}
}

func (g *RadioGroup) Selected() int { return g.selected }

func (g *RadioGroup) Radios() []*Radio { return g.radios }

func (g *RadioGroup) Select(i int) {
	if i == g.selected {
		return
	}
	g.selected = i
	for j, r := range g.radios {
		r.mark.turn(j == i)
	}
	if g.onChange != nil {
		g.onChange(i)
	}
}

// Radio is a list row belonging to a RadioGroup.
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

func (r *Radio) choose() { r.group.Select(r.index) }
