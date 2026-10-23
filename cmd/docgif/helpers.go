package main

import (
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

// page puts a small muted title above body, which takes the rest of the
// screen, like the pages of the gallery. Demos that are about something
// else wrap their content in it and keep it out of their snippets.
func page(title string, body ui.Node) ui.Node {
	t := ui.CurrentTheme
	header := widget.NewLabel(title, widget.Font(t.SmallFont), widget.TextColor(t.Muted),
		ui.Padding(ui.LTRB(4, 1, 4, 0)))
	ui.BoxOf(body).SetGrow(1)
	return ui.NewColumn(ui.Fill(), ui.Children(header, body))
}

// padded is a list whose rows keep clear of the edges, for bordered
// widgets such as buttons and fields.
func padded(rows ...ui.Node) *widget.List {
	l := widget.NewList()
	l.Column().SetPadding(ui.XY(2, 1))
	l.Column().SetGap(ui.CurrentTheme.Gap)
	ui.Add(l, rows...)
	return l
}

// mutedCaption is a small muted line naming the row below it.
func mutedCaption(text string) *widget.Label {
	t := ui.CurrentTheme
	return widget.NewLabel(text, widget.Font(t.SmallFont), widget.TextColor(t.Muted),
		ui.Padding(ui.LTRB(4, 2, 4, 0)))
}
