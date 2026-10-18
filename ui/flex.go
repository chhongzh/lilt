package ui

type Axis uint8

const (
	Horizontal Axis = iota
	Vertical
)

// MainAlign distributes free space along a Flex's axis when no child grows.
type MainAlign uint8

const (
	MainStart MainAlign = iota
	MainCenter
	MainEnd
	MainSpaceBetween
)

// CrossAlign places children across a Flex's axis.
type CrossAlign uint8

const (
	// CrossStretch fills the cross axis for children whose size there is Auto.
	CrossStretch CrossAlign = iota
	CrossStart
	CrossCenter
	CrossEnd
)

// Flex lines children up along one axis.
type Flex struct {
	Box
	axis  Axis
	gap   float32
	main  MainAlign
	cross CrossAlign
}

func NewRow(opts ...Option) *Flex    { return newFlex(Horizontal, opts) }
func NewColumn(opts ...Option) *Flex { return newFlex(Vertical, opts) }

func newFlex(a Axis, opts []Option) *Flex {
	f := &Flex{axis: a}
	Apply(f, opts)
	return f
}

func (f *Flex) Axis() Axis { return f.axis }

// SetAxis lets widgets that embed a Flex by value become columns.
func (f *Flex) SetAxis(a Axis) {
	if f.axis != a {
		f.axis = a
		Relayout()
	}
}

func (f *Flex) Gap() float32 { return f.gap }

func (f *Flex) SetGap(g float32) {
	if f.gap != g {
		f.gap = g
		Relayout()
	}
}

func (f *Flex) SetMainAlign(m MainAlign) {
	if f.main != m {
		f.main = m
		Relayout()
	}
}

func (f *Flex) SetCrossAlign(c CrossAlign) {
	if f.cross != c {
		f.cross = c
		Relayout()
	}
}

// split and join convert between x/y and main/cross for the flex axis.
func (f *Flex) split(v Vec2) (main, cross float32) {
	if f.axis == Horizontal {
		return v.X, v.Y
	}
	return v.Y, v.X
}

func (f *Flex) join(main, cross float32) Vec2 {
	if f.axis == Horizontal {
		return Vec2{main, cross}
	}
	return Vec2{cross, main}
}

func (f *Flex) margins(b *Box) (main, cross, lead, crossLead float32) {
	m := b.margin
	if f.axis == Horizontal {
		return m.Horizontal(), m.Vertical(), m.Left, m.Top
	}
	return m.Vertical(), m.Horizontal(), m.Top, m.Left
}

func (f *Flex) crossAuto(b *Box) bool {
	if f.axis == Horizontal {
		return b.height.mode == autoSize
	}
	return b.width.mode == autoSize
}

func (f *Flex) MeasureChildren(max Vec2) Vec2 {
	var main, cross float32
	n := 0
	for _, c := range f.children {
		cb := c.box()
		if !cb.visible() {
			continue
		}
		s := Measure(c, shrink(max, cb.margin))
		mm, cm, _, _ := f.margins(cb)
		sm, sc := f.split(s)
		main += sm + mm
		cross = max32(cross, sc+cm)
		n++
	}
	if n > 1 {
		main += f.gap * float32(n-1)
	}
	return f.join(main, cross)
}

func (f *Flex) ArrangeChildren(content Rect) {
	availMain, availCross := f.split(Vec2{content.W, content.H})
	var used, grow float32
	n := 0
	for _, c := range f.children {
		cb := c.box()
		if !cb.visible() {
			continue
		}
		cb.measured = Measure(c, shrink(Vec2{content.W, content.H}, cb.margin))
		mm, _, _, _ := f.margins(cb)
		sm, _ := f.split(cb.measured)
		used += sm + mm
		grow += cb.grow
		n++
	}
	if n > 1 {
		used += f.gap * float32(n-1)
	}
	free := availMain - used
	pos, gap := float32(0), f.gap
	if grow == 0 && free > 0 {
		switch f.main {
		case MainCenter:
			pos = free / 2
		case MainEnd:
			pos = free
		case MainSpaceBetween:
			if n > 1 {
				gap += free / float32(n-1)
			}
		}
	}
	// Growing children also absorb overflow, so a long label shrinks
	// instead of pushing its siblings out of the row.
	if grow == 0 {
		free = max32(free, 0)
	}
	for _, c := range f.children {
		cb := c.box()
		if !cb.visible() {
			continue
		}
		mm, cm, lead, crossLead := f.margins(cb)
		sm, sc := f.split(cb.measured)
		if grow > 0 {
			sm = max32(sm+free*cb.grow/grow, 0)
		}
		room := availCross - cm
		var cp float32
		switch {
		case f.cross == CrossStretch && f.crossAuto(cb):
			sc = room
		case f.cross == CrossCenter:
			cp = (room - sc) / 2
		case f.cross == CrossEnd:
			cp = room - sc
		}
		o := f.join(pos+lead, cp+crossLead)
		s := f.join(sm, sc)
		Arrange(c, Rect{content.X + o.X, content.Y + o.Y, s.X, s.Y})
		pos += sm + mm + gap
	}
}
