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
		name:  "select-pick",
		build: selectPick,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("drawer open")
			s.Tap(ui.KeyDown, 0.3)
			s.Tap(ui.KeyDown, 0.5)
			s.Shot("on Ludicrous")
			s.Tap(ui.KeyOK, 0.8)
			s.Shot("picked")
			s.Wait(1.2)
		},
	})
	register(demo{
		name:  "select-fail",
		build: selectFail,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.7)
			s.Tap(ui.KeyDown, 0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("loading")
			s.Wait(0.9)
			s.Shot("old option back")
			s.Wait(1)
		},
	})
}

func selectPick() ui.Node {
	// --8<-- [start:select-pick]
	speed := widget.NewSelect("Speed", []string{"Slow", "Normal", "Fast", "Ludicrous"}, 1,
		widget.OnChange(func(from, to int) {
			ui.Do(func() { widget.Toast("Speed saved") })
		}))
	list := widget.NewList()
	ui.Add(list, speed)
	// --8<-- [end:select-pick]
	return page("Select", list)
}

func selectFail() ui.Node {
	// --8<-- [start:select-fail]
	region := widget.NewSelect("Region", []string{"Europe", "Asia", "America"}, 0,
		widget.OnChange(func(from, to int) error {
			time.Sleep(time.Second)
			return errors.New("locked by admin") // Europe comes back
		}))
	// --8<-- [end:select-fail]
	list := widget.NewList()
	ui.Add(list, region)
	return page("Select", list)
}
