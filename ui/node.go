package ui

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
)

// Node is any node in the tree. It is satisfied only by types that embed
// Box, which carries the tree links, layout properties and the animated offset.
// Behaviour comes from the optional interfaces below, found by type assertion.
type Node interface {
	box() *Box
}

// Painter paints the node's own content, before its children.
type Painter interface{ Draw(c gfx.Canvas) }

// OverPainter paints after the node's children.
type OverPainter interface{ DrawOver(c gfx.Canvas) }

// Updater runs once per frame while the node is attached and visible.
type Updater interface{ Update(dt float32) }

// Measurer reports the content size of a leaf, excluding padding.
type Measurer interface{ Measure(max Vec2) Vec2 }

// Layouter positions children itself, like Flex and Grid. content is the
// padded area in the node's own coordinates.
type Layouter interface {
	MeasureChildren(max Vec2) Vec2
	ArrangeChildren(content Rect)
}

// ContentHost redirects Add and RemoveChildren to an inner node, so a List
// can take items while its scroll machinery stays private.
type ContentHost interface{ Content() Node }

// Leaf marks nodes that refuse children.
type Leaf interface{ Leaf() }

// ChildRemover is told when a child detaches, including a child removing itself.
type ChildRemover interface{ ChildRemoved(child Node) }

// Focuser is told when the node gains or loses focus.
type Focuser interface{ FocusChanged(focused bool) }

type flags uint16

const (
	flagHidden flags = 1 << iota
	flagFocusable
	flagClip
	flagAnimateLayout
	flagLaidOut
	flagBackground
	flagBorder
)

// Box is the base every node embeds. Properties go through setters so a
// change can schedule a relayout.
type Box struct {
	parent   Node
	children []Node

	width, height       Size
	anchor, origin, pos Vec2
	margin, padding     Insets
	grow                float32
	flags               flags
	bg, border          gfx.Color
	radius              int16

	rect     Rect // relative to the parent's top-left corner
	measured Vec2 // scratch space for Layouters

	offsetX, offsetY anim.Value
}

func (b *Box) box() *Box { return b }

// BoxOf returns the Box embedded in d.
func BoxOf(d Node) *Box { return d.box() }

func (b *Box) Parent() Node { return b.parent }

// Children must not be modified; use Add, Insert and Remove.
func (b *Box) Children() []Node { return b.children }

// Rect is the laid-out rectangle relative to the parent, without the offset.
func (b *Box) Rect() Rect { return b.rect }

// ContentRect is the area inside the padding, in pixels and in the node's
// own coordinates.
func (b *Box) ContentRect() Rect {
	p := padPx(b)
	return Rect{p.Left, p.Top, b.rect.W - p.Horizontal(), b.rect.H - p.Vertical()}
}

// OffsetX is an animated translation applied on top of layout. Children
// inherit it because they are drawn relative to their parent.
func (b *Box) OffsetX() *anim.Value { return &b.offsetX }
func (b *Box) OffsetY() *anim.Value { return &b.offsetY }

func (b *Box) Width() Size  { return b.width }
func (b *Box) Height() Size { return b.height }

func (b *Box) SetWidth(s Size) {
	if b.width != s {
		b.width = s
		Relayout()
	}
}

func (b *Box) SetHeight(s Size) {
	if b.height != s {
		b.height = s
		Relayout()
	}
}

func (b *Box) Anchor() Vec2 { return b.anchor }

// SetAnchor picks the point in the parent the node is attached to.
func (b *Box) SetAnchor(v Vec2) {
	if b.anchor != v {
		b.anchor = v
		Relayout()
	}
}

func (b *Box) Origin() Vec2 { return b.origin }

// SetOrigin picks the point of the node that sits on the anchor.
func (b *Box) SetOrigin(v Vec2) {
	if b.origin != v {
		b.origin = v
		Relayout()
	}
}

func (b *Box) Position() Vec2 { return b.pos }

func (b *Box) SetPosition(v Vec2) {
	if b.pos != v {
		b.pos = v
		Relayout()
	}
}

func (b *Box) Margin() Insets { return b.margin }

func (b *Box) SetMargin(i Insets) {
	if b.margin != i {
		b.margin = i
		Relayout()
	}
}

func (b *Box) Padding() Insets { return b.padding }

func (b *Box) SetPadding(i Insets) {
	if b.padding != i {
		b.padding = i
		Relayout()
	}
}

func (b *Box) Grow() float32 { return b.grow }

// SetGrow sets the share of a Flex's free space this node takes.
func (b *Box) SetGrow(g float32) {
	if b.grow != g {
		b.grow = g
		Relayout()
	}
}

func (b *Box) Visible() bool { return b.flags&flagHidden == 0 }

// SetVisible hides or shows the node. Hidden nodes take no space, are not
// drawn and cannot hold focus.
func (b *Box) SetVisible(v bool) {
	if b.Visible() == v {
		return
	}
	b.set(flagHidden, !v)
	Relayout()
	if !v {
		focus.leaving(b)
	}
}

func (b *Box) Focusable() bool { return b.flags&flagFocusable != 0 }

func (b *Box) SetFocusable(f bool) {
	b.set(flagFocusable, f)
	if !f {
		focus.leaving(b)
	}
}

// SetClip confines drawing of the node and its children to its rect.
func (b *Box) SetClip(c bool) {
	b.set(flagClip, c)
	Invalidate()
}

// SetAnimateLayout makes position changes from layout glide instead of jump.
func (b *Box) SetAnimateLayout(a bool) { b.set(flagAnimateLayout, a) }

// SetBackground fills the node's rect with c, rounded by the corner radius.
func (b *Box) SetBackground(c gfx.Color) {
	b.bg = c
	b.set(flagBackground, true)
	Invalidate()
}

// SetBorder strokes a one-pixel outline in c.
func (b *Box) SetBorder(c gfx.Color) {
	b.border = c
	b.set(flagBorder, true)
	Invalidate()
}

// SetRadius rounds the background and border, in units.
func (b *Box) SetRadius(r int16) {
	b.radius = r
	Invalidate()
}

func (b *Box) set(f flags, on bool) {
	if on {
		b.flags |= f
	} else {
		b.flags &^= f
	}
}

func (b *Box) visible() bool { return b.flags&flagHidden == 0 }
