package main

import (
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "field-keyboard",
		build: fieldKeyboard,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("keyboard open")
			s.Taps(ui.KeyRight, 4, 0.12)
			s.Tap(ui.KeyOK, 0.2)
			s.Taps(ui.KeyLeft, 2, 0.12)
			s.Tap(ui.KeyOK, 0.2)
			s.Tap(ui.KeyDown, 0.12)
			s.Taps(ui.KeyLeft, 2, 0.12)
			s.Tap(ui.KeyOK, 0.3)
			s.Shot("typed tea")
			s.Taps(ui.KeyDown, 2, 0.12)
			s.Taps(ui.KeyRight, 2, 0.2)
			s.Shot("on OK")
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("done")
		},
	})
	register(demo{
		name:  "field-method",
		build: fieldMethod,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("presets open")
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("Hallway picked")
		},
	})
}

func fieldKeyboard() ui.Node {
	// --8<-- [start:field-keyboard]
	results := widget.NewLabel("", widget.Font(ui.CurrentTheme.SmallFont),
		ui.Padding(ui.CurrentTheme.RowPadding))
	search := widget.NewInput("Search", widget.OnChange(func(from, to string) {
		ui.Do(func() { results.SetText("Results for " + to) })
	}))
	// --8<-- [end:field-keyboard]
	return page("Input", padded(search, results))
}

// --8<-- [start:field-method-type]
// presets is a TextInputMethod for boards with too few keys for a
// keyboard: it offers a handful of ready-made names in a drawer.
type presets []string

func (p presets) Open(t ui.TextTarget) {
	d := widget.NewDrawer(widget.Bottom, "Pick a name")
	for _, name := range p {
		ui.Add(d, widget.NewListItem(name, widget.OnClick(func() {
			ui.Do(func() {
				d.Close()
				t.SetText(name)
			})
		})))
	}
	d.Open()
}

// --8<-- [end:field-method-type]

func fieldMethod() ui.Node {
	// --8<-- [start:field-method]
	ui.TextInput = presets{"Kitchen", "Hallway", "Garage"}
	room := widget.NewInput("Room")
	// --8<-- [end:field-method]
	return page("Input", padded(room))
}
