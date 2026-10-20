package widget

import (
	"time"

	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

type buttonState uint8

const (
	idle buttonState = iota
	loading
	succeeded
	failed
)

// resultHold is how long the success or failure mark stays before the
// label returns.
const resultHold = 1.2

// Button runs its OnClick handler when OK is released over it. A handler
// still running after a moment swaps the label for loading dots and then
// for a check, or for a cross on error; a quick one leaves the label alone
// unless it fails. With HoldToConfirm the cursor shrinks while OK is held
// and the button fires when the cursor is gone.
type Button struct {
	ui.Box
	Control
	label *Label
	state buttonState
	timer float32
	// mark is the result still on screen; result runs 0 to 1 as it rises in
	// and on to 2 as it leaves upward.
	mark buttonState
	// after runs on the UI goroutine once the handler returns.
	after func(error)

	hold     float32
	holding  bool
	fired    bool
	progress anim.Value
	result   anim.Value
}

func NewButton(text string, opts ...ui.Option) *Button {
	b := &Button{label: NewLabel(text)}
	b.label.SetAnchor(ui.Center)
	b.label.SetOrigin(ui.Center)
	b.label.OffsetY().SetSpring(anim.Snappy)
	ui.Attach(b, b.label)
	t := ui.CurrentTheme
	b.SetFocusable(true)
	b.SetClip(true)
	b.SetPadding(t.ButtonPadding)
	b.SetBorder(t.Foreground)
	b.SetRadius(t.Radius)
	b.OffsetX().SetSpring(anim.Bouncy)
	b.progress.SetSpring(anim.Snappy)
	b.progress.SetPrecision(0.002)
	b.result.SetSpring(anim.Snappy)
	b.result.SetPrecision(0.01)
	ui.Apply(b, opts)
	return b
}

func (b *Button) Label() *Label { return b.label }

func (b *Button) SetFont(f *gfx.Font) { b.label.SetFont(f) }

func (b *Button) SetHold(d time.Duration) { b.hold = float32(d.Seconds()) }

func (b *Button) HandleKey(e *ui.KeyEvent) bool {
	if e.Key != ui.KeyOK {
		return false
	}
	switch e.Phase {
	case ui.Press:
		if b.state == idle && !b.Busy() && b.hold > 0 {
			b.holding, b.fired = true, false
		}
	case ui.Release:
		switch {
		case b.hold == 0:
			if !e.Canceled && b.state == idle {
				b.fire()
			}
		case b.fired:
			// A quick handler is already done, so bring the cursor back.
			if !b.Busy() && b.state == idle {
				b.progress.Set(0)
			}
		default:
			b.progress.Set(0)
		}
		b.holding, b.fired = false, false
	}
	return true
}

func (b *Button) fire() { b.run(clickCall(b.click), b.done) }

func (b *Button) done(err error) {
	waited := b.state == loading
	b.state = idle
	switch {
	case err != nil:
		b.state = failed
		shake(b)
	case waited:
		b.state = succeeded
	}
	if b.state == idle {
		if !b.holding {
			b.progress.Set(0)
		}
	} else {
		b.timer, b.mark = 0, b.state
		b.label.OffsetY().Set(-b.Rect().H)
		b.result.Jump(0)
		b.result.Set(1)
	}
	if b.after != nil {
		b.after(err)
	}
}

func (b *Button) Update(dt float32) {
	b.tick(dt)
	if b.holding && !b.fired {
		p := b.progress.Current() + dt/b.hold
		if p >= 1 {
			p, b.fired = 1, true
			b.fire()
		}
		b.progress.Jump(p)
		ui.Invalidate()
	}
	switch b.state {
	case idle:
		if b.Loading() {
			b.state = loading
			b.label.OffsetY().Set(-b.Rect().H)
		}
	case succeeded, failed:
		b.timer += dt
		if b.timer >= resultHold {
			b.state = idle
			b.result.Set(2)
			b.label.OffsetY().Jump(b.Rect().H)
			b.label.OffsetY().Set(0)
			if !b.holding {
				b.progress.Set(0)
			}
		}
	}
}

// CursorShape eats the cursor from the left as the hold progresses, so it
// runs out at the right edge.
func (b *Button) CursorShape(r ui.Rect) ui.Rect {
	cut := r.W * clampf(b.progress.Get(), 0, 1)
	r.X += cut
	r.W -= cut
	return r
}

func (b *Button) Draw(c gfx.Canvas) {
	r := b.Rect()
	cx, cy := ui.Round(r.W/2), ui.Round(r.H/2)
	col := ui.CurrentTheme.Foreground
	if b.state == loading {
		drawDots(c, cx, cy, b.elapsed, col)
	}
	if b.mark != idle && b.result.Get() < 1.99 {
		y := cy + ui.Round((1-b.result.Get())*r.H)
		if b.mark == succeeded {
			c.Line(cx-4, y, cx-1, y+3, col)
			c.Line(cx-1, y+3, cx+4, y-3, col)
		} else {
			c.Line(cx-3, y-3, cx+3, y+3, col)
			c.Line(cx-3, y+3, cx+3, y-3, col)
		}
	}
}
