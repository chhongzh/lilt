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
		name:  "threading-do",
		build: threadingDo,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("checking")
			s.Wait(0.9)
			s.Shot("updated")
			s.Wait(0.4)
		},
	})
	register(demo{
		name:  "threading-rollback",
		build: threadingRollback,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.4)
			s.Shot("waiting")
			s.Wait(0.75)
			s.Shot("rolled back")
			s.Wait(0.6)
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.1)
			s.Shot("quick handler")
			s.Wait(0.8)
		},
	})
}

func threadingDo() ui.Node {
	// --8<-- [start:threading-do]
	status := widget.NewLabel("Version 1.2", widget.Font(ui.CurrentTheme.SmallFont),
		ui.Padding(ui.CurrentTheme.RowPadding))
	check := widget.NewListItem("Check updates", widget.OnClick(func() {
		time.Sleep(time.Second) // the network; no frame waits for it
		ui.Do(func() { status.SetText("Version 1.3 ready") })
	}))
	// --8<-- [end:threading-do]
	list := widget.NewList()
	ui.Add(list, check, status)
	return page("Firmware", list)
}

func threadingRollback() ui.Node {
	// --8<-- [start:threading-rollback]
	cloud := widget.NewCheckBox("Cloud sync", false,
		widget.OnChange(func(from, to bool) error {
			time.Sleep(800 * time.Millisecond)
			return errors.New("offline") // the box flips back
		}))
	vibrate := widget.NewCheckBox("Vibrate", true,
		widget.OnChange(func(from, to bool) {
			time.Sleep(50 * time.Millisecond) // under 150ms: no dots
		}))
	// --8<-- [end:threading-rollback]
	list := widget.NewList()
	ui.Add(list, cloud, vibrate)
	return page("Settings", list)
}
