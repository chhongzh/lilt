package main

import (
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"

	"tinygo.org/x/tinyfont/freesans"
)

func init() {
	script := func(s *reel.Script) {
		s.Wait(0.5)
		s.Tap(ui.KeyDown, 0.5)
		s.Shot("button")
		s.Tap(ui.KeyDown, 0.5)
		s.Shot("row")
		s.Taps(ui.KeyUp, 2, 0.5)
	}
	register(demo{
		name:   "theme-default",
		build:  themeDefault,
		script: script,
	})
	register(demo{
		name:   "theme-large",
		w:      192,
		h:      96,
		build:  themeLarge,
		script: script,
	})
}

// --8<-- [start:theme-settings]
func settings() ui.Node {
	t := ui.CurrentTheme
	list := widget.NewList()
	list.Column().SetPadding(ui.XY(2, 1)) // units, like every inset
	list.Column().SetGap(t.Gap)
	ui.Add(list,
		widget.NewCheckBox("Wi-Fi", true),
		widget.NewButton("Pair device"),
		widget.NewListItem("About"),
	)
	return list
}

// --8<-- [end:theme-settings]

func themeDefault() ui.Node { return settings() }

func themeLarge() ui.Node {
	// --8<-- [start:theme-large]
	// Change the theme before building: widgets read it when constructed.
	ui.CurrentTheme.Font = &freesans.Regular9pt7b // 22 px lines, so 1 unit = 22/15 px
	root := settings()
	// --8<-- [end:theme-large]
	return root
}
