// Package gallery is a demo of every standard widget, shared by the
// firmware example and the host preview.
package gallery

import (
	"errors"
	"strconv"
	"time"

	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

// Demo is a page the home list links to.
type Demo struct {
	Title string
	Build func() ui.Node
}

// Home builds the gallery's home screen. Extra demos, such as board-specific
// pages, are listed before the standard ones.
func Home(extra ...Demo) ui.Node {
	home := widget.NewList()
	for _, d := range append(extra, demos...) {
		build := d.Build
		ui.Add(home, widget.NewListItem(d.Title, widget.OnClick(func() {
			ui.Do(func() { ui.Push(build()) })
		})))
	}
	return Page("Gallery", home)
}

var demos = []Demo{
	{"Buttons", buttons},
	{"Toggles", toggles},
	{"Select", selects},
	{"Popups", popups},
	{"Chips", chips},
	{"Flex", flex},
	{"Grid", grid},
	{"Scroll", scrolls},
	{"Nested", nested},
	{"Input", inputs},
	{"Todo", todos},
	{"Long text", longText},
}

// Page stacks a small muted header above body, which takes the rest.
func Page(title string, body ui.Node) ui.Node {
	t := ui.CurrentTheme
	header := widget.NewLabel(title, widget.Font(t.SmallFont), widget.TextColor(t.Muted),
		ui.Padding(ui.LTRB(4, 1, 4, 0)))
	ui.BoxOf(body).SetGrow(1)
	return ui.NewColumn(ui.Fill(), ui.Children(header, body))
}

// padded is a list whose rows keep clear of the edges, for bordered widgets.
func padded() *widget.List {
	l := widget.NewList()
	l.Column().SetPadding(ui.XY(2, 1))
	l.Column().SetGap(ui.CurrentTheme.Gap)
	return l
}

func buttons() ui.Node {
	l := padded()
	taps := 0
	var tap *widget.Button
	tap = widget.NewButton("Tap me", widget.OnClick(func() {
		taps++
		ui.Do(func() { tap.Label().SetText("Tapped " + strconv.Itoa(taps)) })
	}))
	ui.Add(l,
		tap,
		widget.NewButton("Hold to delete", widget.HoldToConfirm(800*time.Millisecond),
			widget.OnClick(func() {})),
		widget.NewButton("Save", widget.OnClick(func() error {
			time.Sleep(time.Second)
			return nil
		})),
		widget.NewButton("Upload", widget.OnClick(func() error {
			time.Sleep(time.Second)
			return errors.New("offline")
		})),
	)
	return Page("Buttons", l)
}

func toggles() ui.Node {
	l := widget.NewList()
	ui.Add(l,
		widget.NewCheckBox("Wi-Fi", false),
		widget.NewCheckBox("Bluetooth", true),
		// A slow handler that fails: dots, a shake, and the box flips back.
		widget.NewCheckBox("Cloud sync", false, widget.OnChange(func(from, to bool) error {
			time.Sleep(time.Second)
			return errors.New("offline")
		})),
	)
	g := widget.NewRadioGroup(0)
	for _, s := range []string{"Dark", "Light", "Auto"} {
		ui.Add(l, widget.NewRadio(g, s))
	}
	return Page("Toggles", l)
}

func selects() ui.Node {
	l := widget.NewList()
	ui.Add(l,
		widget.NewSelect("Speed", []string{"Slow", "Normal", "Fast", "Ludicrous"}, 1),
		widget.NewSelect("Language", []string{
			"English", "Deutsch", "Espanol", "Francais", "Italiano",
			"Nederlands", "Polski", "Portugues", "Svenska",
		}, 0),
	)
	return Page("Select", l)
}

// popups opens a drawer from each edge, a dialog whose action takes a
// moment, and a toast.
func popups() ui.Node {
	l := widget.NewList()
	edges := []struct {
		name string
		edge widget.Edge
	}{
		{"Drawer from bottom", widget.Bottom},
		{"Drawer from top", widget.Top},
		{"Drawer from left", widget.Left},
		{"Drawer from right", widget.Right},
	}
	for _, e := range edges {
		ui.Add(l, widget.NewListItem(e.name, widget.OnClick(func() {
			ui.Do(func() { menu(e.edge).Open() })
		})))
	}
	ui.Add(l,
		widget.NewListItem("Dialog", widget.OnClick(func() {
			ui.Do(func() {
				widget.NewDialog("Delete map?", "Blue Zenith goes for good, scores and all",
					widget.NewButton("Keep"),
					widget.NewButton("Delete", widget.OnClick(func() {
						time.Sleep(700 * time.Millisecond)
						ui.Do(func() { widget.Toast("Deleted") })
					})),
				).Open()
			})
		})),
		widget.NewListItem("Toast", widget.OnClick(func() {
			ui.Do(func() { widget.Toast("Saved to slot 2") })
		})),
	)
	return Page("Popups", l)
}

// menu is a drawer of three entries, each closing it and saying which was picked.
func menu(edge widget.Edge) *widget.Drawer {
	d := widget.NewDrawer(edge, "Menu")
	for _, s := range []string{"Profile", "Settings", "Sign out"} {
		ui.Add(d, widget.NewListItem(s, widget.OnClick(func() {
			ui.Do(func() {
				d.Close()
				widget.Toast(s)
			})
		})))
	}
	return d
}

// caption is a small muted line naming the row below it.
func caption(s string) *widget.Label {
	t := ui.CurrentTheme
	return widget.NewLabel(s, widget.Font(t.SmallFont), widget.TextColor(t.Muted),
		ui.Padding(ui.LTRB(4, 2, 4, 0)))
}

func chips() ui.Node {
	l := widget.NewList()
	ui.Add(l,
		caption("Difficulty"),
		widget.NewChoice([]string{"Easy", "Normal", "Hard", "Insane", "Expert"}, 1),
		caption("Mods"),
		widget.NewToggles([]string{"HD", "HR", "DT", "FL", "NF", "EZ"}),
		caption("Rate"),
		widget.NewChoice([]string{"0.5x", "0.75x", "1x", "1.25x", "1.5x", "2x"}, 2),
	)
	return Page("Chips", l)
}

// flex lays three boxes out in a row, with chips below switching the row's
// main alignment and making the middle box grow. The boxes glide to their
// new places.
func flex() ui.Node {
	t := ui.CurrentTheme
	small := widget.Font(t.SmallFont)
	row := ui.NewRow(ui.Gap(2), ui.Padding(ui.XY(4, 2)))
	boxes := make([]*widget.Label, 3)
	for i, s := range []string{"1", "2", "3"} {
		boxes[i] = widget.NewLabel(s, small, widget.TextAlign(0.5), ui.Border(t.Foreground),
			ui.Radius(t.Radius), ui.Padding(ui.XY(5, 1)), ui.AnimateLayout())
		ui.Add(row, boxes[i])
	}
	aligns := []ui.MainAlign{ui.MainStart, ui.MainCenter, ui.MainEnd, ui.MainSpaceBetween}
	main := widget.NewChoice([]string{"Start", "Center", "End", "Between"}, 0,
		small, widget.OnChange(func(from, to int) {
			ui.Do(func() { row.SetMainAlign(aligns[to]) })
		}))
	grow := widget.NewToggles([]string{"Grow 2"}, small, widget.OnChange(func(from, to uint32) {
		ui.Do(func() { boxes[1].SetGrow(float32(to & 1)) })
	}))

	// The row stays put above the controls so every change stays in view.
	// Growing rather than filling keeps the list's bottom on screen.
	l := widget.NewList(ui.Grow(1))
	ui.Add(l, caption("Main align"), main, caption("Grow"), grow)
	return Page("Flex", ui.NewColumn(ui.Children(row, l)))
}

// grid is a PIN pad: twelve buttons in three equal columns.
func grid() ui.Node {
	t := ui.CurrentTheme
	small := widget.Font(t.SmallFont)
	const prompt = "Enter a PIN"
	pin := ""
	display := widget.NewLabel(prompt, small, widget.TextAlign(0.5))
	press := func(k string) {
		switch {
		case k == "C":
			pin = ""
		case k == "<":
			if pin != "" {
				pin = pin[:len(pin)-1]
			}
		case len(pin) < 8:
			pin += k
		}
		if pin == "" {
			display.SetText(prompt)
		} else {
			display.SetText(pin)
		}
	}
	g := ui.NewGrid(ui.Columns(ui.Fr(1), ui.Fr(1), ui.Fr(1)), ui.Gap(2))
	for _, k := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "C", "0", "<"} {
		ui.Add(g, widget.NewButton(k, small, ui.Padding(ui.XY(2, 1)),
			widget.OnClick(func() { ui.Do(func() { press(k) }) })))
	}
	l := padded()
	ui.Add(l, display, g)
	return Page("Grid", l)
}

// scrolls shows two shelves of maps that scroll sideways, like a TV menu:
// up and down pick a shelf, left and right move along it, and the cards
// follow the cursor as a wave.
func scrolls() ui.Node {
	small := widget.Font(ui.CurrentTheme.SmallFont)
	playing := widget.NewLabel("Pick a map", small, ui.Padding(ui.XY(4, 1)))
	shelves := []struct {
		name string
		maps []string
	}{
		{"Recent", []string{"Blue Zenith", "Freedom Dive", "Teo", "Glory Days",
			"Harumachi Clover", "Senbonzakura"}},
		{"Favourites", []string{"Kani Do-Raku", "Snow Drive", "Big Black",
			"Night of Knights", "Airman"}},
	}
	l := widget.NewList(ui.Grow(1))
	for _, sh := range shelves {
		row := ui.NewRow(ui.Gap(3), ui.Padding(ui.XY(2, 1)))
		for _, m := range sh.maps {
			ui.Add(row, widget.NewButton(m, small, ui.Padding(ui.XY(6, 4)),
				widget.OnClick(func() { ui.Do(func() { playing.SetText("Playing " + m) }) })))
		}
		// The bottom padding keeps the thumb clear of the cursor.
		ui.Add(l, caption(sh.name), widget.NewScroll(row, widget.ScrollHorizontal,
			ui.FillWidth(), ui.Padding(ui.LTRB(0, 0, 0, 3))))
	}
	return Page("Scroll", ui.NewColumn(ui.Children(playing, l)))
}

// nested puts scrolls inside a scrolling page: a short framed list that
// scrolls on its own, a shelf that scrolls sideways and a grid that scrolls
// both ways. At the end of an inner scroll the keys pass to the page.
func nested() ui.Node {
	t := ui.CurrentTheme
	small := widget.Font(t.SmallFont)
	picked := widget.NewLabel("Pick a cell", small, ui.Padding(t.RowPadding))

	inbox := widget.NewList(ui.Height(ui.Px(40)), ui.Border(t.Foreground), ui.Radius(t.Radius),
		ui.Margin(ui.XY(2, 0)), ui.Padding(ui.LTRB(1, 1, 3, 1)))
	for i := 1; i <= 8; i++ {
		ui.Add(inbox, widget.NewListItem("Message "+strconv.Itoa(i), small))
	}

	shelf := ui.NewRow(ui.Gap(t.Gap), ui.Padding(ui.XY(2, 1)))
	for _, s := range []string{"HD", "HR", "DT", "FL", "NF", "EZ", "SO", "PF"} {
		ui.Add(shelf, widget.NewButton(s, small, ui.Padding(ui.XY(5, 1))))
	}

	// A board of cells too wide and too tall for the screen.
	board := ui.NewGrid(ui.Gap(t.Gap), ui.Padding(ui.All(2)),
		ui.Columns(ui.Fixed(28), ui.Fixed(28), ui.Fixed(28), ui.Fixed(28), ui.Fixed(28), ui.Fixed(28)))
	for r := 'A'; r <= 'F'; r++ {
		for c := 1; c <= 6; c++ {
			cell := string(r) + strconv.Itoa(c)
			ui.Add(board, widget.NewButton(cell, small, ui.Padding(ui.XY(1, 1)),
				widget.OnClick(func() { ui.Do(func() { picked.SetText("Picked " + cell) }) })))
		}
	}

	l := widget.NewList(ui.Grow(1))
	ui.Add(l,
		caption("Inbox"), inbox,
		caption("Shelf"), widget.NewScroll(shelf, widget.ScrollHorizontal, ui.FillWidth(),
			ui.Padding(ui.LTRB(0, 0, 0, 3))),
		caption("Board"), widget.NewScroll(board, widget.ScrollBoth, ui.FillWidth(),
			ui.Height(ui.Px(44)), ui.Padding(ui.LTRB(0, 0, 3, 3))),
		widget.NewListItem("Below the board"),
	)
	return Page("Nested", ui.NewColumn(ui.Children(picked, l)))
}

func inputs() ui.Node {
	l := padded()
	ui.Add(l, widget.NewInput("Name"), widget.NewInput("Password"))
	return Page("Input", l)
}

func longText() ui.Node {
	t := ui.CurrentTheme
	l := widget.NewList()
	ui.Add(l,
		widget.NewListItem("A title far too long to fit, scrolling while focused"),
		widget.NewLabel("This one scrolls all the time, focus or not.",
			widget.Font(t.SmallFont), widget.MarqueeMode(widget.MarqueeAlways),
			ui.Padding(t.RowPadding)),
	)
	for i := 1; i <= 30; i++ {
		ui.Add(l, widget.NewListItem("Row "+strconv.Itoa(i)))
	}
	return Page("Long text", l)
}
