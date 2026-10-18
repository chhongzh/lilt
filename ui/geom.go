package ui

import "github.com/chhongzh/lilt/gfx"

// Vec2 is a point or size in layout space. Layout works in float32 so
// animated positions stay smooth; drawing rounds to pixels.
type Vec2 struct{ X, Y float32 }

// Alignment points for Anchor and Origin.
var (
	TopLeft      = Vec2{0, 0}
	TopCenter    = Vec2{0.5, 0}
	TopRight     = Vec2{1, 0}
	CenterLeft   = Vec2{0, 0.5}
	Center       = Vec2{0.5, 0.5}
	CenterRight  = Vec2{1, 0.5}
	BottomLeft   = Vec2{0, 1}
	BottomCenter = Vec2{0.5, 1}
	BottomRight  = Vec2{1, 1}
)

// Rect is a layout rectangle.
type Rect struct{ X, Y, W, H float32 }

func (r Rect) Right() float32  { return r.X + r.W }
func (r Rect) Bottom() float32 { return r.Y + r.H }
func (r Rect) Center() Vec2    { return Vec2{r.X + r.W/2, r.Y + r.H/2} }

// Pixels rounds each edge, so adjacent rects stay adjacent after rounding.
func (r Rect) Pixels() gfx.Rect {
	x, y := Round(r.X), Round(r.Y)
	return gfx.Rect{X: x, Y: y, W: Round(r.Right()) - x, H: Round(r.Bottom()) - y}
}

// Insets are margins or padding.
type Insets struct{ Left, Top, Right, Bottom float32 }

func All(v float32) Insets           { return Insets{v, v, v, v} }
func XY(x, y float32) Insets         { return Insets{x, y, x, y} }
func LTRB(l, t, r, b float32) Insets { return Insets{l, t, r, b} }
func (i Insets) Horizontal() float32 { return i.Left + i.Right }
func (i Insets) Vertical() float32   { return i.Top + i.Bottom }

type sizeMode uint8

const (
	autoSize sizeMode = iota
	fixedSize
	relativeSize
)

// Size is one axis of a node's size: fixed pixels, a fraction of the
// parent's content area, or Auto to fit the content.
type Size struct {
	mode sizeMode
	v    float32
}

var Auto = Size{}

func Px(v float32) Size  { return Size{fixedSize, v} }
func Rel(f float32) Size { return Size{relativeSize, f} }

// Unbounded is the available space along a scrolling axis. Relative sizes
// fall back to Auto there, since a fraction of infinity means nothing.
const Unbounded float32 = 1 << 24

// Round rounds halves away from zero.
func Round(x float32) int16 {
	if x < 0 {
		return -int16(-x + 0.5)
	}
	return int16(x + 0.5)
}

func clamp(x, lo, hi float32) float32 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func absf(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
