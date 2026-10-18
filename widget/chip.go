package widget

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

const (
	chipBarH = 2
	// chipGrow is how far past the label the bar may overshoot.
	chipGrow = 0.25
	// chipStagger slows the wave less than in a list: a row of chips is
	// short, and a long tail would lag behind the cursor.
	chipStagger = 0.15
)

// ChipSpring grows the bar under a chip out of its middle, past the label's
// width, so switching a chip on lands with a snap.
var ChipSpring = &anim.Spring{Response: 0.262, Damping: 0.5}

// Chip is one option in a Choice or Toggles: a word that is on or off, with
// a bar under it while on. Releasing OK over it runs its OnClick; the row it
// belongs to decides what being on means.
type Chip struct {
	ui.Box
	label   *Label
	on      bool
	fill    anim.Value
	onClick func()
}

func NewChip(text string, opts ...ui.Option) *Chip {
	c := &Chip{label: NewLabel(text, MarqueeMode(MarqueeOff))}
	ui.Attach(c, c.label)
	c.SetFocusable(true)
	c.SetPadding(ui.LTRB(3, 1, 3, 1+chipBarH))
	c.fill.SetSpring(ChipSpring)
	c.fill.SetPrecision(0.01)
	ui.Apply(c, opts)
	return c
}

func (c *Chip) Label() *Label { return c.label }

func (c *Chip) SetFont(f *gfx.Font) { c.label.SetFont(f) }

func (c *Chip) SetOnClick(fn func()) { c.onClick = fn }

func (c *Chip) On() bool { return c.on }

// SetOn switches the chip, the bar growing or shrinking into place.
func (c *Chip) SetOn(on bool) {
	if on == c.on {
		return
	}
	c.on = on
	if on {
		c.fill.Set(1)
	} else {
		c.fill.Set(0)
	}
	ui.Invalidate()
}

func (c *Chip) HandleKey(e *ui.KeyEvent) bool {
	if e.Key != ui.KeyOK {
		return false
	}
	if e.Phase == ui.Release && !e.Canceled && c.onClick != nil {
		c.onClick()
	}
	return true
}

func (c *Chip) Draw(cv gfx.Canvas) {
	r := c.Rect()
	full := ui.BoxOf(c.label).Rect().W
	w := min(c.fill.Get()*full, full*(1+chipGrow), r.W)
	if px := ui.Round(w); px > 0 {
		y := ui.Round(r.H) - chipBarH - 1
		cv.FillRect(gfx.Rect{X: ui.Round((r.W - w) / 2), Y: y, W: px, H: chipBarH},
			ui.CurrentTheme.Foreground)
	}
}

// chipRow fills s with a row of chips, one per label, each calling click
// with its index when activated.
func chipRow(s *Scroll, labels []string, click func(int)) []*Chip {
	row := ui.NewRow(ui.Gap(2))
	s.init(row, ScrollHorizontal)
	s.SetStagger(chipStagger)
	// A thumb under the row would read as another chip's bar.
	s.SetScrollBar(false)
	s.SetWidth(ui.Rel(1))
	chips := make([]*Chip, len(labels))
	for i, l := range labels {
		chips[i] = NewChip(l, OnClick(func() { click(i) }))
		ui.Add(row, chips[i])
	}
	return chips
}

func fontAll(chips []*Chip, f *gfx.Font) {
	for _, c := range chips {
		c.SetFont(f)
	}
}

// Choice is a horizontal row of chips of which exactly one is on: left and
// right move over the options, and releasing OK switches to the one under
// the cursor, the bar shrinking under the old option and growing under the
// new one. It scrolls when the options are wider than the screen.
type Choice struct {
	Scroll
	chips    []*Chip
	chosen   int
	onChange func(int)
}

func NewChoice(labels []string, chosen int, opts ...ui.Option) *Choice {
	c := &Choice{chosen: -1}
	c.chips = chipRow(&c.Scroll, labels, c.choose)
	c.Choose(chosen)
	for _, ch := range c.chips {
		ch.fill.Jump(ch.fill.Target())
	}
	ui.Apply(c, opts)
	return c
}

func (c *Choice) Chosen() int { return c.chosen }

func (c *Choice) Chips() []*Chip { return c.chips }

func (c *Choice) SetOnChange(fn func(int)) { c.onChange = fn }

func (c *Choice) SetFont(f *gfx.Font) { fontAll(c.chips, f) }

// Choose switches to option i without calling OnChange; -1 turns all off.
func (c *Choice) Choose(i int) {
	for j, ch := range c.chips {
		ch.SetOn(j == i)
	}
	c.chosen = i
}

func (c *Choice) choose(i int) {
	if i == c.chosen {
		return
	}
	c.Choose(i)
	if c.onChange != nil {
		c.onChange(i)
	}
}

// Toggles is a horizontal row of chips, any number of which are on;
// releasing OK flips the one under the cursor.
type Toggles struct {
	Scroll
	chips    []*Chip
	onChange func(int)
}

func NewToggles(labels []string, opts ...ui.Option) *Toggles {
	t := &Toggles{}
	t.chips = chipRow(&t.Scroll, labels, t.toggle)
	ui.Apply(t, opts)
	return t
}

func (t *Toggles) Chips() []*Chip { return t.chips }

func (t *Toggles) On(i int) bool { return t.chips[i].On() }

// SetOn switches chip i without calling OnChange.
func (t *Toggles) SetOn(i int, on bool) { t.chips[i].SetOn(on) }

func (t *Toggles) SetFont(f *gfx.Font) { fontAll(t.chips, f) }

// SetOnChange registers fn, called with the index of a chip after it flips.
func (t *Toggles) SetOnChange(fn func(int)) { t.onChange = fn }

func (t *Toggles) toggle(i int) {
	t.SetOn(i, !t.On(i))
	if t.onChange != nil {
		t.onChange(i)
	}
}
