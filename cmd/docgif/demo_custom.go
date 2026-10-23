package main

import (
	"strconv"
	"time"

	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	// --8<-- [start:drivers-docgif]
	register(demo{
		name:  "custom-slider",
		build: customSlider,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Taps(ui.KeyRight, 2, 0.4)
			s.Shot("two steps right")
			// Hold sends no repeats; a held key is Press, Repeats, Release.
			s.Emit(ui.KeyRight, ui.Press)
			s.Wait(0.2)
			for range 3 {
				s.Emit(ui.KeyRight, ui.Repeat)
				s.Wait(0.12)
			}
			s.Emit(ui.KeyRight, ui.Release)
			s.Wait(0.5)
			s.Shot("held")
			s.Emit(ui.KeyLeft, ui.Press)
			s.Wait(0.3)
			s.Shot("leaning left")
			s.Emit(ui.KeyLeft, ui.Release)
			s.Wait(0.5)
			s.Tap(ui.KeyDown, 0.6)
			s.Shot("bass focused")
			s.Tap(ui.KeyLeft, 0.6)
		},
	})
	// --8<-- [end:drivers-docgif]
	register(demo{
		name:  "custom-battery",
		build: customBattery,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Shot("30 percent")
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("reading")
			s.Wait(1)
			s.Shot("charged")
			s.Wait(0.5)
		},
	})
}

// --8<-- [start:custom-slider-type]
// Slider picks one of steps+1 positions with left and right.
type Slider struct {
	ui.Box
	value, steps int
	knob         anim.Value // in steps; springs toward value
	onChange     func(from, to int)
}

func NewSlider(value, steps int, opts ...ui.Option) *Slider {
	s := &Slider{value: value, steps: steps}
	s.knob.Jump(float32(value))
	s.knob.SetSpring(anim.Snappy)
	s.knob.SetPrecision(0.01) // the value is in steps, not pixels
	s.SetFocusable(true)
	s.SetPadding(ui.CurrentTheme.RowPadding)
	ui.Apply(s, opts) // last, so options override the defaults
	return s
}

// OnSlide sets the handler run with the old and new step after the user
// moves the knob.
func OnSlide(fn func(from, to int)) ui.Option {
	return func(d ui.Node) { d.(*Slider).onChange = fn }
}

func (s *Slider) Value() int { return s.value }

// SetValue moves the knob without calling the handler.
func (s *Slider) SetValue(v int) {
	s.value = min(max(v, 0), s.steps)
	s.knob.Set(float32(s.value))
}

// --8<-- [end:custom-slider-type]

// --8<-- [start:custom-slider-draw]
// knobRadius is 3 units in pixels, so the knob grows with the theme font.
func knobRadius() int16 { return ui.Round(3 * ui.Unit()) }

func (s *Slider) Measure(max ui.Vec2) ui.Vec2 {
	return ui.Vec2{X: max.X, Y: float32(2*knobRadius() + 3)}
}

// knobX is the knob's centre in the node's own pixels, where it is now.
func (s *Slider) knobX() float32 {
	c := s.ContentRect()
	r := float32(knobRadius())
	return c.X + r + (c.W-2*r)*s.knob.Get()/float32(s.steps)
}

func (s *Slider) Draw(c gfx.Canvas) {
	t := ui.CurrentTheme
	area := s.ContentRect().Pixels()
	cy := area.Y + area.H/2
	track := gfx.Rect{X: area.X, Y: cy - 2, W: area.W, H: 5}
	c.StrokeRoundRect(track, 2, t.Foreground)
	kx := ui.Round(s.knobX())
	track.W = kx - track.X
	c.FillRoundRect(track, 2, t.Foreground)
	r := knobRadius()
	c.FillCircle(kx, cy, r+1, t.Background)
	c.FillCircle(kx, cy, r, t.Foreground)
}

// --8<-- [end:custom-slider-draw]

// --8<-- [start:custom-slider-keys]
func (s *Slider) HandleKey(e *ui.KeyEvent) bool {
	dir := 0
	switch e.Key {
	case ui.KeyLeft:
		dir = -1
	case ui.KeyRight:
		dir = 1
	default:
		return false // up and down still move focus
	}
	switch {
	case e.Phase == ui.Press:
		// Feedback only: lean toward the key.
		s.knob.Set(float32(s.value) + 0.3*float32(dir))
	case e.Fires():
		s.move(dir)
	default:
		// Canceled, or the Release after repeats that already moved.
		s.knob.Set(float32(s.value))
	}
	return true // claim the key, so Repeat and Release come here too
}

func (s *Slider) move(dir int) {
	from := s.value
	s.SetValue(from + dir)
	if s.value != from && s.onChange != nil {
		go s.onChange(from, s.value) // never on the UI goroutine
	}
}

// --8<-- [end:custom-slider-keys]

// --8<-- [start:custom-slider-cursor]
// CursorShape wraps the cursor round the knob instead of the whole row.
func (s *Slider) CursorShape(r ui.Rect) ui.Rect {
	d := float32(2*knobRadius() + 5)
	return ui.Rect{X: r.X + s.knobX() - d/2, Y: r.Y + (r.H-d)/2, W: d, H: d}
}

// --8<-- [end:custom-slider-cursor]

func customSlider() ui.Node {
	// --8<-- [start:custom-slider]
	small := widget.Font(ui.CurrentTheme.SmallFont)
	level := widget.NewLabel("Volume 40%", small)
	volume := NewSlider(4, 10, OnSlide(func(from, to int) {
		ui.Do(func() { level.SetText("Volume " + strconv.Itoa(to*10) + "%") })
	}))
	bass := NewSlider(5, 10)
	// --8<-- [end:custom-slider]
	return page("Slider", padded(level, volume, widget.NewLabel("Bass", small), bass))
}

// --8<-- [start:custom-battery-type]
// Battery draws a charge level from 0 to 1. It takes no keys: the row it
// sits in handles activation.
type Battery struct {
	ui.Box
	level anim.Value
}

func NewBattery(level float32, opts ...ui.Option) *Battery {
	b := &Battery{}
	b.level.Jump(level)
	b.level.SetPrecision(0.002)
	ui.Apply(b, opts)
	return b
}

func (b *Battery) SetLevel(l float32) { b.level.Set(l) }

func (b *Battery) Leaf() {} // refuses children

func (b *Battery) Measure(max ui.Vec2) ui.Vec2 {
	return ui.Vec2{X: 18 * ui.Unit(), Y: 9 * ui.Unit()}
}

func (b *Battery) Draw(c gfx.Canvas) {
	fg := ui.CurrentTheme.Foreground
	r := b.ContentRect().Pixels()
	body := gfx.Rect{X: r.X, Y: r.Y, W: r.W - 2, H: r.H}
	c.StrokeRoundRect(body, 2, fg)
	c.FillRect(gfx.Rect{X: body.Right(), Y: r.Y + r.H/4, W: 2, H: r.H / 2}, fg)
	fill := body.Inset(2)
	// The spring overshoots a little; keep the fill inside the body.
	fill.W = ui.Round(float32(fill.W) * min(max(b.level.Get(), 0), 1))
	c.FillRect(fill, fg)
}

// --8<-- [end:custom-battery-type]

func customBattery() ui.Node {
	// --8<-- [start:custom-battery]
	bat := NewBattery(0.3)
	row := widget.NewListItem("Battery", widget.Trailing(bat), widget.OnClick(func() {
		time.Sleep(time.Second) // read the fuel gauge
		ui.Do(func() { bat.SetLevel(0.85) })
	}))
	// --8<-- [end:custom-battery]
	list := widget.NewList()
	ui.Add(list, row, widget.NewListItem("Storage"), widget.NewListItem("About"))
	return page("Status", list)
}
