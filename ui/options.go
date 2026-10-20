package ui

import "github.com/chhongzh/lilt/gfx"

// Option configures a node at construction, which lets trees be written as
// one nested expression.
type Option func(Node)

// Apply runs opts against d. Widget constructors call it last.
func Apply(d Node, opts []Option) {
	for _, o := range opts {
		o(d)
	}
}

// Children adds nodes to a container.
func Children(nodes ...Node) Option {
	return func(d Node) { Add(d, nodes...) }
}

// Ref stores the node being built in *p, for keeping a handle on a node
// buried in a nested tree.
func Ref[T Node](p *T) Option {
	return func(d Node) { *p = d.(T) }
}

func Width(s Size) Option  { return func(d Node) { d.box().SetWidth(s) } }
func Height(s Size) Option { return func(d Node) { d.box().SetHeight(s) } }

// Fill makes the node as large as its parent's content area.
func Fill() Option {
	return func(d Node) {
		d.box().SetWidth(Rel(1))
		d.box().SetHeight(Rel(1))
	}
}

func FillWidth() Option  { return Width(Rel(1)) }
func FillHeight() Option { return Height(Rel(1)) }

func Anchor(v Vec2) Option { return func(d Node) { d.box().SetAnchor(v) } }
func Origin(v Vec2) Option { return func(d Node) { d.box().SetOrigin(v) } }

// Align sets anchor and origin together, the usual way to pin a node to a
// side or corner of its parent.
func Align(v Vec2) Option {
	return func(d Node) {
		d.box().SetAnchor(v)
		d.box().SetOrigin(v)
	}
}

func Position(x, y float32) Option { return func(d Node) { d.box().SetPosition(Vec2{x, y}) } }
func Margin(i Insets) Option       { return func(d Node) { d.box().SetMargin(i) } }
func Padding(i Insets) Option      { return func(d Node) { d.box().SetPadding(i) } }
func Grow(g float32) Option        { return func(d Node) { d.box().SetGrow(g) } }
func Focusable() Option            { return func(d Node) { d.box().SetFocusable(true) } }
func Hidden() Option               { return func(d Node) { d.box().SetVisible(false) } }
func Clip() Option                 { return func(d Node) { d.box().SetClip(true) } }
func AnimateLayout() Option        { return func(d Node) { d.box().SetAnimateLayout(true) } }
func Background(c gfx.Color) Option {
	return func(d Node) { d.box().SetBackground(c) }
}
func Border(c gfx.Color) Option { return func(d Node) { d.box().SetBorder(c) } }
func Radius(r int16) Option     { return func(d Node) { d.box().SetRadius(r) } }

// Gap sets the spacing between children of a Flex or Grid.
func Gap(g float32) Option {
	return func(d Node) {
		s, ok := d.(interface{ SetGap(float32) })
		if !ok {
			panic("ui.Gap: option needs a *ui.Flex or *ui.Grid")
		}
		s.SetGap(g)
	}
}

func Main(m MainAlign) Option {
	return func(d Node) { flexOf(d, "ui.Main").SetMainAlign(m) }
}

func Cross(c CrossAlign) Option {
	return func(d Node) { flexOf(d, "ui.Cross").SetCrossAlign(c) }
}

func Columns(cols ...Track) Option {
	return func(d Node) {
		g, ok := d.(*Grid)
		if !ok {
			panic("ui.Columns: option needs a *ui.Grid")
		}
		g.SetColumns(cols...)
	}
}

func flexOf(d Node, name string) *Flex {
	f, ok := d.(*Flex)
	if !ok {
		panic(name + ": option needs a *ui.Flex")
	}
	return f
}

// Stack is a plain container whose children are placed by anchor and origin.
type Stack struct{ Box }

func NewStack(opts ...Option) *Stack {
	s := &Stack{}
	Apply(s, opts)
	return s
}
