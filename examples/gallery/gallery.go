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

// Build pushes the gallery's home screen onto the current app. Extra demos,
// such as board-specific pages, are listed before the standard ones.
func Build(extra ...Demo) {
	home := widget.NewList()
	for _, d := range append(extra, demos...) {
		build := d.Build
		ui.Add(home, widget.NewListItem(d.Title, widget.OnClick(func() {
			ui.Push(build(), ui.SlideLeft)
		})))
	}
	ui.Push(Page("Gallery", home), ui.SlideLeft)
}

var demos = []Demo{
	{"Buttons", buttons},
	{"Toggles", toggles},
	{"Select", selects},
	{"Chips", chips},
	{"Flex", flex},
	{"Grid", grid},
	{"Scroll", scrolls},
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
	l.Column().SetGap(2)
	return l
}

func buttons() ui.Node {
	l := padded()
	taps := 0
	var tap *widget.Button
	tap = widget.NewButton("Tap me", widget.OnClick(func() {
		taps++
		tap.Label().SetText("Tapped " + strconv.Itoa(taps))
	}))
	ui.Add(l,
		tap,
		widget.NewButton("Hold to delete", widget.HoldToConfirm(800*time.Millisecond),
			widget.OnClick(func() {})),
		widget.NewButton("Save", widget.OnAction(func() error {
			time.Sleep(time.Second)
			return nil
		})),
		widget.NewButton("Upload", widget.OnAction(func() error {
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
	)
	g := widget.NewRadioGroup(0, nil)
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
		small, widget.OnChange(func(i int) { row.SetMainAlign(aligns[i]) }))
	grow := widget.NewToggles([]string{"Grow 2"}, small)
	grow.SetOnChange(func(int) {
		var g float32
		if grow.On(0) {
			g = 1
		}
		boxes[1].SetGrow(g)
	})

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
			widget.OnClick(func() { press(k) })))
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
				widget.OnClick(func() { playing.SetText("Playing " + m) })))
		}
		// The bottom padding keeps the thumb clear of the cursor.
		ui.Add(l, caption(sh.name), widget.NewScroll(row, widget.ScrollHorizontal,
			ui.FillWidth(), ui.Padding(ui.LTRB(0, 0, 0, 3))))
	}
	return Page("Scroll", ui.NewColumn(ui.Children(playing, l)))
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
			ui.Padding(ui.XY(4, 1))),
	)
	for i := 1; i <= 30; i++ {
		ui.Add(l, widget.NewListItem("Row "+strconv.Itoa(i)))
	}
	return Page("Long text", l)
}
