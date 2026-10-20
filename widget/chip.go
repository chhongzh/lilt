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
// a bar under it while on. Releasing OK over it runs its OnClick handler;
// the row it belongs to decides what being on means. While a handler is
// slow, loading dots stand in for the bar.
type Chip struct {
	ui.Box
	Control
	label *Label
	on    bool
	fill  anim.Value
	// pick runs on the UI goroutine first, for the row the chip is in.
	pick func()
	// spin is the control whose loading the chip shows: its own, or its
	// row's while the row reports a change made on this chip.
	spin *Control
}

func NewChip(text string, opts ...ui.Option) *Chip {
	c := &Chip{label: NewLabel(text)}
	c.spin = &c.Control
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
	if e.Phase == ui.Release && !e.Canceled && !c.Busy() {
		if c.pick != nil {
			c.pick()
		}
		c.run(clickCall(c.click), func(err error) {
			if err != nil {
				shake(c)
			}
		})
	}
	return true
}

func (c *Chip) Update(dt float32) { c.tick(dt) }

func (c *Chip) Draw(cv gfx.Canvas) {
	r := c.Rect()
	y := ui.Round(r.H) - chipBarH - 1
	col := ui.CurrentTheme.Foreground
	if c.spin.Loading() {
		drawDots(cv, ui.Round(r.W/2), y+1, c.spin.elapsed, col)
		return
	}
	full := ui.BoxOf(c.label).Rect().W
	w := min(c.fill.Get()*full, full*(1+chipGrow), r.W)
	if px := ui.Round(w); px > 0 {
		cv.FillRect(gfx.Rect{X: ui.Round((r.W - w) / 2), Y: y, W: px, H: chipBarH}, col)
	}
}

// chipRow fills s with a row of chips, one per label, each calling pick
// with its index on the UI goroutine when activated.
func chipRow(s *Scroll, labels []string, pick func(int)) []*Chip {
	row := ui.NewRow(ui.Gap(2))
	s.init(row, ScrollHorizontal)
	s.SetStagger(chipStagger)
	// A thumb under the row would read as another chip's bar.
	s.SetScrollBar(false)
	s.SetWidth(ui.Rel(1))
	s.OffsetX().SetSpring(anim.Bouncy)
	chips := make([]*Chip, len(labels))
	for i, l := range labels {
		chips[i] = NewChip(l)
		chips[i].pick = func() { pick(i) }
		ui.Add(row, chips[i])
	}
	return chips
}

func fontAll(chips []*Chip, f *gfx.Font) {
	for _, c := range chips {
		c.SetFont(f)
	}
}

// report runs a row's change handler, the chip that made the change showing
// the loading, and calls undo if the handler fails.
func report(row *Control, s *Scroll, chip *Chip, call func() error, undo func()) {
	chip.spin = row
	row.run(call, func(err error) {
		chip.spin = &chip.Control
		if err != nil {
			undo()
			shake(s)
		}
	})
}

// Choice is a horizontal row of chips of which exactly one is on: left and
// right move over the options, and releasing OK switches to the one under
// the cursor, the bar shrinking under the old option and growing under the
// new one. It scrolls when the options are wider than the screen. OnChange
// gets the old and new index.
type Choice struct {
	Scroll
	Control
	chips  []*Chip
	chosen int
}

func NewChoice(labels []string, chosen int, opts ...ui.Option) *Choice {
	c := &Choice{chosen: -1}
	c.kind = 0
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

func (c *Choice) SetFont(f *gfx.Font) { fontAll(c.chips, f) }

// Choose switches to option i without calling OnChange; -1 turns all off.
func (c *Choice) Choose(i int) {
	for j, ch := range c.chips {
		ch.SetOn(j == i)
	}
	c.chosen = i
}

func (c *Choice) choose(i int) {
	from := c.chosen
	if i == from || c.Busy() {
		return
	}
	c.Choose(i)
	report(&c.Control, &c.Scroll, c.chips[i], changeCall(c.change, from, i), func() { c.Choose(from) })
}

func (c *Choice) Update(dt float32) { c.tick(dt) }

// Toggles is a horizontal row of chips, any number of which are on;
// releasing OK flips the one under the cursor. OnChange gets the old and
// new set as bit masks, bit i for chip i, so a row holds at most 32 chips.
type Toggles struct {
	Scroll
	Control
	chips []*Chip
}

func NewToggles(labels []string, opts ...ui.Option) *Toggles {
	t := &Toggles{}
	t.kind = uint32(0)
	t.chips = chipRow(&t.Scroll, labels, t.toggle)
	ui.Apply(t, opts)
	return t
}

func (t *Toggles) Chips() []*Chip { return t.chips }

func (t *Toggles) On(i int) bool { return t.chips[i].On() }

// SetOn switches chip i without calling OnChange.
func (t *Toggles) SetOn(i int, on bool) { t.chips[i].SetOn(on) }

// Mask is the set of chips that are on, bit i for chip i.
func (t *Toggles) Mask() uint32 {
	var m uint32
	for i, c := range t.chips {
		if c.On() {
			m |= 1 << i
		}
	}
	return m
}

// SetMask switches every chip to match m without calling OnChange.
func (t *Toggles) SetMask(m uint32) {
	for i, c := range t.chips {
		c.SetOn(m&(1<<i) != 0)
	}
}

func (t *Toggles) SetFont(f *gfx.Font) { fontAll(t.chips, f) }

func (t *Toggles) toggle(i int) {
	if t.Busy() {
		return
	}
	from := t.Mask()
	t.SetOn(i, !t.On(i))
	report(&t.Control, &t.Scroll, t.chips[i], changeCall(t.change, from, t.Mask()), func() { t.SetMask(from) })
}

func (t *Toggles) Update(dt float32) { t.tick(dt) }
