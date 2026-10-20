package widget

import (
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Edge is the side of the screen a Drawer slides in from.
type Edge uint8

const (
	Bottom Edge = iota
	Top
	Left
	Right
)

// Drawer is a column that slides in from one edge of the screen, with an
// optional small title. Top and bottom drawers span the width and are as
// tall as their content; side drawers span the height and are as wide as
// it. Back closes a drawer.
type Drawer struct {
	ui.Flex
	edge  Edge
	title *Label
}

func NewDrawer(edge Edge, title string, opts ...ui.Option) *Drawer {
	d := &Drawer{edge: edge}
	d.SetAxis(ui.Vertical)
	d.SetPadding(ui.LTRB(0, 2, 0, 1))
	if title != "" {
		t := ui.CurrentTheme
		d.title = NewLabel(title, Font(t.SmallFont), TextColor(t.Muted), ui.Padding(t.RowPadding))
		ui.Attach(d, d.title)
	}
	ui.Apply(d, opts)
	return d
}

func (d *Drawer) Edge() Edge { return d.edge }

// Title is the label above the content, or nil.
func (d *Drawer) Title() *Label { return d.title }

// Open slides the drawer in. Call it on the UI goroutine.
func (d *Drawer) Open() { ui.Open(d, ui.DrawerBottom+ui.Popup(d.edge)) }

// Close slides the drawer out again.
func (d *Drawer) Close() { ui.Close(d) }

// TitleHeight is the height the title takes, 0 without one.
func (d *Drawer) TitleHeight() float32 {
	if d.title == nil {
		return 0
	}
	return float32(gfx.LineHeight(d.title.Font())) + d.title.Padding().Vertical()*ui.Unit()
}
