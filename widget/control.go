package widget

import (
	"math"

	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// ClickFunc is an OnClick handler. Returning an error shows the action
// failed.
type ClickFunc interface{ func() | func() error }

// ChangeFunc is an OnChange handler for one of the value types widgets
// report: bool for CheckBox, int for Choice, Select and RadioGroup, string
// for Input, and uint32 for Toggles, bit i being chip i. Returning an error
// puts the old value back.
type ChangeFunc interface {
	func(from, to bool) | func(from, to bool) error |
		func(from, to int) | func(from, to int) error |
		func(from, to string) | func(from, to string) error |
		func(from, to uint32) | func(from, to uint32) error
}

// loadingDelay is how long a handler may run before the widget shows it is
// waiting, so quick handlers never flash a spinner.
const loadingDelay = 0.15

// Control holds a widget's event handlers and the state they drive. Every
// handler runs on its own goroutine, never on the UI goroutine, so a slow
// one cannot stall a frame; meanwhile the widget ignores activation and,
// past loadingDelay, shows loading dots. Handlers change the tree only
// inside ui.Do.
//
// Interactive widgets embed Control next to their ui.Box, which is how they
// get OnClick and OnChange.
type Control struct {
	click   any
	change  any
	kind    any // zero value of the type OnChange reports, nil for none
	busy    bool
	elapsed float32
}

// OnClick sets the handler run when the widget is activated.
func (c *Control) OnClick[F ClickFunc](fn F) { c.click = fn }

// OnChange sets the handler run with the old and new value after the user
// changes it. Setting a value from code does not call it. It panics when
// the handler's type is not the one the widget reports.
func (c *Control) OnChange[F ChangeFunc](fn F) {
	if !accepts(c.kind, fn) {
		panic("widget: OnChange handler does not match the value type of this widget")
	}
	c.change = fn
}

func (c *Control) control() *Control { return c }

// Busy reports whether a handler is still running.
func (c *Control) Busy() bool { return c.busy }

// Loading reports whether a handler has run long enough to show it.
func (c *Control) Loading() bool { return c.busy && c.elapsed >= loadingDelay }

// tick advances the loading timer and keeps the dots moving.
func (c *Control) tick(dt float32) {
	if !c.busy {
		return
	}
	c.elapsed += dt
	if c.elapsed >= loadingDelay {
		ui.Invalidate()
	}
}

// run starts call on its own goroutine unless one is in flight and passes
// its error to done on the UI goroutine. A nil call finishes at once.
func (c *Control) run(call func() error, done func(error)) bool {
	if c.busy {
		return false
	}
	if call == nil {
		if done != nil {
			done(nil)
		}
		return true
	}
	c.busy, c.elapsed = true, 0
	go func() {
		err := call()
		ui.Do(func() {
			c.busy = false
			ui.Invalidate()
			if done != nil {
				done(err)
			}
		})
	}()
	return true
}

func clickCall(h any) func() error {
	switch fn := h.(type) {
	case func():
		return func() error { fn(); return nil }
	case func() error:
		return fn
	}
	return nil
}

func changeCall[T comparable](h any, from, to T) func() error {
	if from == to {
		return nil
	}
	switch fn := h.(type) {
	case func(T, T):
		return func() error { fn(from, to); return nil }
	case func(T, T) error:
		return func() error { return fn(from, to) }
	}
	return nil
}

func accepts(kind, fn any) bool {
	switch kind.(type) {
	case bool:
		switch fn.(type) {
		case func(bool, bool), func(bool, bool) error:
			return true
		}
	case int:
		switch fn.(type) {
		case func(int, int), func(int, int) error:
			return true
		}
	case string:
		switch fn.(type) {
		case func(string, string), func(string, string) error:
			return true
		}
	case uint32:
		switch fn.(type) {
		case func(uint32, uint32), func(uint32, uint32) error:
			return true
		}
	}
	return false
}

// shake throws d sideways, the shared sign that an action failed.
func shake(d ui.Node) { ui.BoxOf(d).OffsetX().Kick(-260) }

// drawWaiting clears a patch at cx, cy and bounces the dots in it, for
// rows and fields whose text may run underneath.
func drawWaiting(c gfx.Canvas, cx, cy int16, t float32) {
	th := ui.CurrentTheme
	u := ui.Unit()
	w, h := ui.Round(18*u), ui.Round(7*u)
	c.FillRect(gfx.Rect{X: cx - w/2, Y: cy - h/2 - 1, W: w, H: h}, th.Background)
	drawDots(c, cx, cy+1, t, th.Foreground)
}

// drawDots paints three dots bouncing in turn around cx, cy.
func drawDots(c gfx.Canvas, cx, cy int16, t float32, col gfx.Color) {
	gap := ui.Round(5 * ui.Unit())
	for i := int16(-1); i <= 1; i++ {
		phase := float64(t)*7 - float64(i)*0.8
		lift := ui.Round(float32(2 * math.Abs(math.Sin(phase))))
		c.FillCircle(cx+i*gap, cy-lift, 1, col)
	}
}
