package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "button-click",
		build: buttonClick,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Emit(ui.KeyOK, ui.Press)
			s.Wait(0.3)
			s.Shot("pressed")
			s.Emit(ui.KeyOK, ui.Release)
			s.Wait(0.5)
			s.Shot("tapped")
			s.Tap(ui.KeyOK, 0.6)
		},
	})
	register(demo{
		name:  "button-states",
		build: buttonStates,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.4)
			s.Shot("loading")
			s.Wait(0.8)
			s.Shot("saved")
			s.Wait(1)
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 1.3)
			s.Shot("failed")
			s.Wait(1.5)
		},
	})
	register(demo{
		name:  "button-hold",
		build: buttonHold,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Hold(ui.KeyOK, 0.4)
			s.Shot("let go early")
			s.Wait(0.8)
			s.Emit(ui.KeyOK, ui.Press)
			s.Wait(0.5)
			s.Shot("holding")
			s.Wait(0.6)
			s.Emit(ui.KeyOK, ui.Release)
			s.Wait(0.3)
			s.Shot("deleted")
			s.Wait(1)
		},
	})
}

func buttonClick() ui.Node {
	// --8<-- [start:button-click]
	taps := 0
	var tap *widget.Button
	tap = widget.NewButton("Tap me", widget.OnClick(func() {
		taps++
		ui.Do(func() { tap.Label().SetText("Tapped " + strconv.Itoa(taps)) })
	}))
	// --8<-- [end:button-click]
	return page("Button", padded(tap))
}

func buttonStates() ui.Node {
	// --8<-- [start:button-states]
	save := widget.NewButton("Save", widget.OnClick(func() error {
		time.Sleep(time.Second) // longer than 150ms: dots, then a check
		return nil
	}))
	upload := widget.NewButton("Upload", widget.OnClick(func() error {
		time.Sleep(time.Second)
		return errors.New("offline") // a shake and a cross
	}))
	// --8<-- [end:button-states]
	return page("Button", padded(save, upload))
}

func buttonHold() ui.Node {
	// --8<-- [start:button-hold]
	del := widget.NewButton("Hold to delete",
		widget.HoldToConfirm(800*time.Millisecond),
		widget.OnClick(func() {
			ui.Do(func() { widget.Toast("Deleted") })
		}),
	)
	// --8<-- [end:button-hold]
	return page("Button", padded(del))
}
