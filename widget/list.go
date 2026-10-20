package widget

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// List is a vertical scrolling column. Add appends rows to the column.
type List struct {
	Scroll
	column *ui.Flex
}

func NewList(opts ...ui.Option) *List {
	l := &List{column: ui.NewColumn()}
	l.init(l.column, ScrollVertical)
	l.SetWidth(ui.Rel(1))
	l.SetHeight(ui.Rel(1))
	// Room for the scroll thumb so it never sits under the cursor.
	l.SetPadding(ui.LTRB(0, 0, 3, 0))
	ui.Apply(l, opts)
	return l
}

// Column is the flex holding the rows, for adjusting gap or alignment.
func (l *List) Column() *ui.Flex { return l.column }

// ListItem is a focusable row with a title and an optional trailing
// widget. Releasing OK over it runs its OnClick handler; while a slow one
// runs, loading dots take the trailing widget's place. Rows with a value,
// such as CheckBox and Select, build on it.
type ListItem struct {
	ui.Flex
	Control
	title    *Label
	trailing ui.Node
	waiting  bool
	// activate replaces the click handler for rows with their own
	// behaviour; it runs on the UI goroutine.
	activate func()
}

func NewListItem(title string, opts ...ui.Option) *ListItem {
	li := &ListItem{}
	li.init(title)
	ui.Apply(li, opts)
	return li
}

func (li *ListItem) init(title string) {
	li.title = NewLabel(title)
	li.title.SetGrow(1)
	ui.Attach(li, li.title)
	li.SetFocusable(true)
	li.SetClip(true)
	li.SetPadding(ui.CurrentTheme.RowPadding)
	li.SetGap(4)
	li.SetCrossAlign(ui.CrossCenter)
	li.OffsetX().SetSpring(anim.Bouncy)
}

func (li *ListItem) Title() *Label { return li.title }

func (li *ListItem) SetFont(f *gfx.Font) { li.title.SetFont(f) }

func (li *ListItem) SetTrailing(d ui.Node) {
	if li.trailing != nil {
		ui.Remove(li.trailing)
	}
	li.trailing = d
	if d != nil {
		ui.BoxOf(d).OffsetY().SetSpring(anim.Snappy)
		ui.Attach(li, d)
	}
}

// HandleKey activates on release, like iOS: pressing only squeezes the
// cursor, and a press that turns into a drag does nothing.
func (li *ListItem) HandleKey(e *ui.KeyEvent) bool {
	if e.Key != ui.KeyOK {
		return false
	}
	if e.Phase == ui.Release && !e.Canceled && !li.Busy() {
		if li.activate != nil {
			li.activate()
		} else {
			li.run(clickCall(li.click), func(err error) {
				if err != nil {
					shake(li)
				}
			})
		}
	}
	return true
}

// Update slides the trailing widget up and out while loading and back in
// from below once the handler returns.
func (li *ListItem) Update(dt float32) {
	li.tick(dt)
	if li.Loading() == li.waiting {
		return
	}
	li.waiting = !li.waiting
	if li.trailing == nil {
		return
	}
	o := ui.BoxOf(li.trailing).OffsetY()
	h := li.Rect().H
	if li.waiting {
		o.Set(-h)
	} else {
		o.Jump(h)
		o.Set(0)
	}
}

func (li *ListItem) DrawOver(c gfx.Canvas) {
	if !li.waiting {
		return
	}
	r := li.Rect()
	p := li.Padding()
	x := ui.Round(r.W - p.Right*ui.Unit() - 9*ui.Unit())
	drawWaiting(c, x, ui.Round(r.H/2), li.elapsed)
}
