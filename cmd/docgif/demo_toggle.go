package main

import (
	"errors"
	"time"

	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "toggle-checkbox",
		build: toggleCheckbox,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("wi-fi on")
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.4)
			s.Shot("flipped, loading")
			s.Wait(0.9)
			s.Shot("flipped back")
			s.Wait(1)
		},
	})
	register(demo{
		name:  "toggle-radio",
		build: toggleRadio,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("light")
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("auto")
			s.Wait(0.4)
		},
	})
}

func toggleCheckbox() ui.Node {
	// --8<-- [start:toggle-checkbox]
	wifi := widget.NewCheckBox("Wi-Fi", false)
	sync := widget.NewCheckBox("Cloud sync", false,
		widget.OnChange(func(from, to bool) error {
			time.Sleep(time.Second)
			return errors.New("offline") // the box flips back
		}))
	list := widget.NewList()
	ui.Add(list, wifi, sync)
	// --8<-- [end:toggle-checkbox]
	return page("CheckBox", list)
}

func toggleRadio() ui.Node {
	// --8<-- [start:toggle-radio]
	t := ui.CurrentTheme
	names := []string{"Dark", "Light", "Auto"}
	status := widget.NewLabel("Using Dark", widget.Font(t.SmallFont),
		widget.TextColor(t.Muted), ui.Padding(t.RowPadding))

	theme := widget.NewRadioGroup(0)
	theme.OnChange(func(from, to int) {
		ui.Do(func() { status.SetText("Using " + names[to]) })
	})
	list := widget.NewList()
	ui.Add(list, status)
	for _, name := range names {
		ui.Add(list, widget.NewRadio(theme, name))
	}
	// --8<-- [end:toggle-radio]
	return page("Radio", list)
}
