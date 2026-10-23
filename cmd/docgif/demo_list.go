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
		name:  "list-items",
		build: listItems,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("loading")
			s.Wait(0.9)
			s.Shot("updated")
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 1.4)
			s.Shot("failed")
			s.Wait(1)
		},
	})
	register(demo{
		name:  "list-wave",
		build: listWave,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Taps(ui.KeyDown, 7, 0.12)
			s.Shot("rows trailing the focus")
			s.Wait(0.8)
			s.Shot("settled")
			s.Taps(ui.KeyUp, 7, 0.12)
			s.Shot("going back up")
			s.Wait(1)
		},
	})
}

func listItems() ui.Node {
	// --8<-- [start:list-items]
	t := ui.CurrentTheme
	small, muted := widget.Font(t.SmallFont), widget.TextColor(t.Muted)
	version := widget.NewLabel("2.3", small, muted)
	update := widget.NewListItem("Update", widget.Trailing(version),
		widget.OnClick(func() {
			time.Sleep(time.Second) // past 150ms the dots replace "2.3"
			ui.Do(func() { version.SetText("2.4") })
		}))
	backup := widget.NewListItem("Back up", widget.Trailing(widget.NewLabel("Never", small, muted)),
		widget.OnClick(func() error {
			time.Sleep(time.Second)
			return errors.New("no SD card") // the row shakes
		}))
	list := widget.NewList()
	ui.Add(list, update, backup)
	// --8<-- [end:list-items]
	return page("List", list)
}

func listWave() ui.Node {
	// --8<-- [start:list-wave]
	small := widget.Font(ui.CurrentTheme.SmallFont)
	list := widget.NewList()
	for i := 1; i <= 24; i++ {
		ui.Add(list, widget.NewListItem("Track "+strconv.Itoa(i), small))
	}
	// --8<-- [end:list-wave]
	return page("List", list)
}
