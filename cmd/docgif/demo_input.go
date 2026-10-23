package main

import (
	"strconv"

	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "input-press",
		build: inputPress,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Emit(ui.KeyDown, ui.Press)
			s.Wait(0.35)
			s.Shot("down held: lean")
			s.Emit(ui.KeyDown, ui.Release)
			s.Wait(0.5)
			s.Shot("released: moved")
			s.Emit(ui.KeyOK, ui.Press)
			s.Wait(0.3)
			s.Shot("ok held: squeeze")
			s.Emit(ui.KeyOK, ui.Release)
			s.Wait(0.5)
			s.Shot("released: toggled")
			s.Wait(0.3)
		},
	})
	register(demo{
		name:  "input-repeat",
		build: inputRepeat,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			// What a KeyPoller sends for a held key: Press, a pause, then
			// Repeats coming faster, and a Release that does not move again.
			s.Emit(ui.KeyDown, ui.Press)
			s.Wait(0.2)
			iv := float32(0.07)
			for range 7 {
				s.Emit(ui.KeyDown, ui.Repeat)
				s.Wait(iv)
				iv = max(iv*0.9, 0.03)
			}
			s.Wait(0.3)
			s.Shot("held")
			s.Emit(ui.KeyDown, ui.Release)
			s.Wait(0.6)
			s.Shot("released")
		},
	})
	register(demo{
		name:  "input-focus",
		build: inputFocus,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyRight, 0.06)
			s.Shot("stretching right")
			s.Wait(0.4)
			s.Tap(ui.KeyRight, 0.45)
			s.Tap(ui.KeyDown, 0.06)
			s.Shot("stretching down")
			s.Wait(0.4)
			s.Tap(ui.KeyLeft, 0.45)
			s.Tap(ui.KeyLeft, 0.45)
			s.Tap(ui.KeyUp, 0.5)
			s.Shot("back home")
		},
	})
}

func inputPress() ui.Node {
	list := widget.NewList()
	ui.Add(list,
		widget.NewCheckBox("Wi-Fi", true),
		widget.NewCheckBox("Bluetooth", false),
		widget.NewCheckBox("Airplane mode", false),
	)
	return page("Settings", list)
}

func inputRepeat() ui.Node {
	list := widget.NewList()
	for i := 1; i <= 20; i++ {
		ui.Add(list, widget.NewListItem("Track "+strconv.Itoa(i)))
	}
	return page("Playlist", list)
}

func inputFocus() ui.Node {
	// --8<-- [start:input-focus]
	small := widget.Font(ui.CurrentTheme.SmallFont)
	apps := ui.NewGrid(ui.Columns(ui.Fr(1), ui.Fr(1), ui.Fr(1)), ui.Gap(2), ui.Padding(ui.XY(2, 1)))
	for _, name := range []string{"Clock", "Maps", "Notes", "Music", "Timer", "Radio"} {
		ui.Add(apps, widget.NewButton(name, small, ui.Padding(ui.XY(2, 3))))
	}
	// --8<-- [end:input-focus]
	return page("Apps", apps)
}
