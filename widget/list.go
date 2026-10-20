package widget

import (
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

// ListItem is a focusable row with a title and an optional trailing widget.
type ListItem struct {
	ui.Flex
	title    *Label
	trailing ui.Node
	onClick  func()
	// activate is the embedding widget's own behaviour, run before onClick.
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
	li.SetPadding(ui.XY(4, 1))
	li.SetGap(4)
	li.SetCrossAlign(ui.CrossCenter)
}

func (li *ListItem) Title() *Label { return li.title }

func (li *ListItem) SetFont(f *gfx.Font) { li.title.SetFont(f) }

func (li *ListItem) SetOnClick(fn func()) { li.onClick = fn }

func (li *ListItem) SetTrailing(d ui.Node) {
	if li.trailing != nil {
		ui.Remove(li.trailing)
	}
	li.trailing = d
	if d != nil {
		ui.Attach(li, d)
	}
}

// HandleKey activates on release, like iOS: pressing only squeezes the
// cursor, and a press that turns into a drag does nothing.
func (li *ListItem) HandleKey(e *ui.KeyEvent) bool {
	if e.Key != ui.KeyOK {
		return false
	}
	if e.Phase == ui.Release && !e.Canceled {
		if li.activate != nil {
			li.activate()
		}
		if li.onClick != nil {
			li.onClick()
		}
	}
	return true
}
