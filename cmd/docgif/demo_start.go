package main

import (
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "start-hello",
		build: hello,
		script: func(s *reel.Script) {
			s.Wait(0.6)
			s.Taps(ui.KeyDown, 2, 0.4)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("checked")
			s.Taps(ui.KeyUp, 2, 0.4)
			s.Tap(ui.KeyOK, 0.8)
			s.Shot("pushed")
			s.Tap(ui.KeyBack, 0.8)
		},
	})
}

// --8<-- [start:start-hello]
func hello() ui.Node {
	list := widget.NewList()
	ui.Add(list,
		widget.NewListItem("Open a page", widget.OnClick(func() {
			// Callbacks run on their own goroutine: change the tree in ui.Do.
			ui.Do(func() { ui.Push(widget.NewLabel("Hello, lilt", widget.TextAlign(0.5))) })
		})),
		widget.NewCheckBox("Wi-Fi", true),
		widget.NewCheckBox("Bluetooth", false),
	)
	return list
}

// --8<-- [end:start-hello]
