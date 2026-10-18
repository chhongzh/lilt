// Package widget is the standard component set built on package ui.
package widget

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Marquee says when a label too long for its box scrolls. Every label that
// overflows scrolls; the mode only picks whether it waits for focus.
type Marquee uint8

const (
	// MarqueeFocused scrolls while the label or an ancestor has focus, so a
	// list only animates the row being read.
	MarqueeFocused Marquee = iota
	MarqueeAlways
)

const (
	marqueeSpeed = 24 // px/s
	marqueePause = 1  // s at each end
	marqueeGap   = 2  // px kept clear past the last glyph
)

type marqueePhase uint8

const (
	holdStart marqueePhase = iota
	forward
	holdEnd
	back
)

// Label is a single line of text. It is vertically centred in its box and
// placed horizontally by its alignment, 0 left to 1 right.
type Label struct {
	ui.Box
	text    string
	font    *gfx.Font
	color   *gfx.Color
	align   float32
	marquee Marquee
	textW   int16

	phase  marqueePhase
	timer  float32
	scroll anim.Value
}

func NewLabel(text string, opts ...ui.Option) *Label {
	l := &Label{text: text}
	l.measure()
	ui.Apply(l, opts)
	return l
}

func (l *Label) Leaf() {}

func (l *Label) Text() string { return l.text }

func (l *Label) SetText(s string) {
	if s == l.text {
		return
	}
	l.text = s
	l.measure()
	l.restart()
	ui.Relayout()
}

func (l *Label) Font() *gfx.Font {
	if l.font != nil {
		return l.font
	}
	return ui.CurrentTheme.Font
}

func (l *Label) SetFont(f *gfx.Font) {
	l.font = f
	l.measure()
	ui.Relayout()
}

// SetColor overrides the theme's foreground colour.
func (l *Label) SetColor(c gfx.Color) {
	l.color = &c
	ui.Invalidate()
}

func (l *Label) SetAlign(x float32) {
	l.align = x
	ui.Invalidate()
}

func (l *Label) SetMarquee(m Marquee) {
	l.marquee = m
	l.restart()
}

func (l *Label) measure() { l.textW = gfx.TextWidth(l.Font(), l.text) }

func (l *Label) restart() {
	l.phase, l.timer = holdStart, 0
	l.scroll.Set(0)
}

func (l *Label) Measure(max ui.Vec2) ui.Vec2 {
	return ui.Vec2{X: min(float32(l.textW), max.X), Y: float32(gfx.LineHeight(l.Font()))}
}

// overflow is how far the text runs past the content box.
func (l *Label) overflow() float32 {
	p := l.Padding()
	return float32(l.textW) + marqueeGap - (l.Rect().W - p.Horizontal())
}

func (l *Label) scrolling() bool {
	if l.marquee == MarqueeAlways {
		return true
	}
	f := ui.Focused()
	return f != nil && ui.IsAncestor(f, l)
}

// shift is how far overflowing text is moved left. It rests on the side it
// is aligned to, so a right-aligned label shows its tail, and the marquee
// carries it to the other side.
func (l *Label) shift(over float32) float32 {
	if l.align > 0.5 {
		return over - l.scroll.Get()
	}
	return l.scroll.Get()
}

func (l *Label) Update(dt float32) {
	over := l.overflow()
	if over <= marqueeGap || !l.scrolling() {
		if l.phase != holdStart || l.timer != 0 {
			l.restart()
		}
		return
	}
	l.timer += dt
	switch l.phase {
	case holdStart:
		if l.timer >= marqueePause {
			l.phase, l.timer = forward, 0
		}
	case forward:
		x := l.timer * marqueeSpeed
		if x >= over {
			x = over
			l.phase, l.timer = holdEnd, 0
		}
		l.scroll.Jump(x)
		ui.Invalidate()
	case holdEnd:
		if l.timer >= marqueePause {
			l.phase, l.timer = back, 0
			l.scroll.Set(0)
		}
	case back:
		if l.scroll.Settled() {
			l.phase, l.timer = holdStart, 0
		}
	}
}

func (l *Label) Draw(c gfx.Canvas) {
	r := l.Rect()
	p := l.Padding()
	box := ui.Rect{X: p.Left, Y: p.Top, W: r.W - p.Horizontal(), H: r.H - p.Vertical()}.Pixels()
	c = c.Clip(box)
	f := l.Font()
	col := ui.CurrentTheme.Foreground
	if l.color != nil {
		col = *l.color
	}
	x := box.X + ui.Round(l.align*float32(box.W-l.textW))
	if over := l.overflow(); over > marqueeGap {
		x = box.X - ui.Round(l.shift(over))
	}
	y := box.Y + (box.H-gfx.LineHeight(f))/2 + gfx.Ascent(f)
	c.Text(f, x, y, l.text, col)
}
