// Package gallery is a demo of every standard widget, shared by the
// firmware example and the host preview.
package gallery

import (
	"errors"
	"strconv"
	"time"

	"osu-mc/ui"
	"osu-mc/widget"
)

// Build pushes the gallery's home screen onto the current app.
func Build() {
	home := widget.NewList()
	for _, d := range demos {
		build := d.build
		ui.Add(home, widget.NewListItem(d.title, widget.OnClick(func() {
			ui.Push(build(), ui.SlideLeft)
		})))
	}
	ui.Push(page("Gallery", home), ui.SlideLeft)
}

var demos = []struct {
	title string
	build func() ui.Drawable
}{
	{"Buttons", buttons},
	{"Toggles", toggles},
	{"Select", selects},
	{"Input", inputs},
	{"Long text", longText},
}

// page stacks a small muted header above body, which takes the rest.
func page(title string, body ui.Drawable) ui.Drawable {
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

func buttons() ui.Drawable {
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
	return page("Buttons", l)
}

func toggles() ui.Drawable {
	l := widget.NewList()
	ui.Add(l,
		widget.NewCheckBox("Wi-Fi", false),
		widget.NewCheckBox("Bluetooth", true),
	)
	g := widget.NewRadioGroup(0, nil)
	for _, s := range []string{"Dark", "Light", "Auto"} {
		ui.Add(l, widget.NewRadio(g, s))
	}
	return page("Toggles", l)
}

func selects() ui.Drawable {
	l := widget.NewList()
	ui.Add(l,
		widget.NewSelect("Speed", []string{"Slow", "Normal", "Fast", "Ludicrous"}, 1),
		widget.NewSelect("Language", []string{
			"English", "Deutsch", "Espanol", "Francais", "Italiano",
			"Nederlands", "Polski", "Portugues", "Svenska",
		}, 0),
	)
	return page("Select", l)
}

func inputs() ui.Drawable {
	l := padded()
	ui.Add(l, widget.NewInput("Name"), widget.NewInput("Password"))
	return page("Input", l)
}

func longText() ui.Drawable {
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
	return page("Long text", l)
}
