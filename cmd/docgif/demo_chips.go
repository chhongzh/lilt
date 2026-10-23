package main

import (
	"errors"
	"math/bits"
	"strconv"
	"time"

	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "chips-choice",
		build: chipsChoice,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Taps(ui.KeyRight, 2, 0.3)
			s.Shot("cursor on Hard")
			s.Emit(ui.KeyOK, ui.Press)
			s.Wait(0.2)
			s.Shot("pressed")
			s.Emit(ui.KeyOK, ui.Release)
			s.Wait(0.6)
			s.Shot("Hard chosen")
			s.Taps(ui.KeyRight, 2, 0.4)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("scrolled to Expert")
		},
	})
	register(demo{
		name:  "chips-toggles",
		build: chipsToggles,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Taps(ui.KeyRight, 2, 0.3)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("HD and DT on")
			s.Tap(ui.KeyRight, 0.3)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("FL on too")
			s.Tap(ui.KeyLeft, 0.3)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("DT off")
		},
	})
	register(demo{
		name:  "chips-fail",
		build: chipsFail,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyRight, 0.3)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("Eco loading")
			s.Wait(0.9)
			s.Shot("Eco on")
			s.Tap(ui.KeyRight, 0.3)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("Turbo loading")
			s.Wait(0.6)
			s.Shot("failed: back to Eco")
			s.Wait(0.8)
		},
	})
}

func chipsChoice() ui.Node {
	// --8<-- [start:chips-choice]
	levels := []string{"Easy", "Normal", "Hard", "Insane", "Expert"}
	stars := widget.NewLabel("2 stars", widget.Font(ui.CurrentTheme.SmallFont),
		ui.Padding(ui.CurrentTheme.RowPadding))
	level := widget.NewChoice(levels, 1, widget.OnChange(func(from, to int) {
		ui.Do(func() { stars.SetText(strconv.Itoa(to+1) + " stars") })
	}))
	// --8<-- [end:chips-choice]
	l := widget.NewList()
	ui.Add(l, mutedCaption("Difficulty"), level, stars)
	return page("Choice", l)
}

func chipsToggles() ui.Node {
	// --8<-- [start:chips-toggles]
	summary := widget.NewLabel("No mods", widget.Font(ui.CurrentTheme.SmallFont),
		ui.Padding(ui.CurrentTheme.RowPadding))
	mods := widget.NewToggles([]string{"HD", "HR", "DT", "FL"},
		widget.OnChange(func(from, to uint32) {
			n := bits.OnesCount32(to) // bit i is chip i
			ui.Do(func() { summary.SetText(strconv.Itoa(n) + " mods on") })
		}))
	// --8<-- [end:chips-toggles]
	l := widget.NewList()
	ui.Add(l, mutedCaption("Mods"), mods, summary)
	return page("Toggles", l)
}

func chipsFail() ui.Node {
	// --8<-- [start:chips-fail]
	fan := widget.NewChoice([]string{"Off", "Eco", "Turbo"}, 0,
		widget.OnChange(func(from, to int) error {
			time.Sleep(time.Second) // talking to the fan controller
			if to == 2 {
				return errors.New("too hot")
			}
			return nil
		}))
	// --8<-- [end:chips-fail]
	l := widget.NewList()
	ui.Add(l, mutedCaption("Fan"), fan)
	return page("Choice", l)
}
