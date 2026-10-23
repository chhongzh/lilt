package main

import (
	"strconv"

	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "scroll-edges",
		build: scrollEdges,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Taps(ui.KeyDown, 9, 0.18)
			s.Wait(0.3)
			s.Shot("last row")
			// What a KeyPoller sends for a held key: Press, then Repeats.
			s.Emit(ui.KeyDown, ui.Press)
			s.Wait(0.2)
			for range 6 {
				s.Emit(ui.KeyDown, ui.Repeat)
				s.Wait(0.07)
			}
			s.Wait(0.2)
			s.Shot("pushed past the end")
			s.Emit(ui.KeyDown, ui.Release)
			s.Wait(0.6)
			s.Shot("sprung back")
		},
	})
	register(demo{
		name:  "scroll-text",
		build: scrollText,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyDown, 0.5)
			s.Shot("one step")
			s.Taps(ui.KeyDown, 3, 0.4)
			s.Shot("reached Install")
			s.Wait(0.4)
		},
	})
	register(demo{
		name:  "scroll-shelf",
		build: scrollShelf,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Taps(ui.KeyRight, 3, 0.25)
			s.Wait(0.1)
			s.Shot("scrolled")
			s.Taps(ui.KeyRight, 2, 0.4)
			s.Shot("last card")
			s.Taps(ui.KeyLeft, 5, 0.12)
			s.Wait(0.5)
			s.Shot("back to start")
		},
	})
	register(demo{
		name:  "scroll-stagger",
		build: scrollStagger,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Taps(ui.KeyDown, 4, 0.12)
			s.Wait(0.05)
			s.Shot("wave")
			s.Wait(0.6)
			s.Taps(ui.KeyUp, 4, 0.12)
			s.Wait(0.05)
			s.Shot("wave back")
			s.Wait(0.6)
		},
	})
	register(demo{
		name:  "scroll-nested",
		build: scrollNested,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Taps(ui.KeyDown, 4, 0.3)
			s.Shot("inbox at its end")
			s.Tap(ui.KeyDown, 0.5)
			s.Shot("page took over")
			s.Taps(ui.KeyDown, 2, 0.4)
			s.Shot("below the inbox")
		},
	})
}

func scrollEdges() ui.Node {
	// --8<-- [start:scroll-edges]
	tracks := ui.NewColumn()
	for i := 1; i <= 10; i++ {
		ui.Add(tracks, widget.NewListItem("Track "+strconv.Itoa(i)))
	}
	// The right padding keeps the thumb clear of the cursor.
	s := widget.NewScroll(tracks, widget.ScrollVertical,
		ui.FillWidth(), ui.Padding(ui.LTRB(0, 0, 3, 0)))
	// --8<-- [end:scroll-edges]
	return page("Scroll", s)
}

func scrollText() ui.Node {
	// --8<-- [start:scroll-text]
	small := widget.Font(ui.CurrentTheme.SmallFont)
	notes := ui.NewColumn()
	ui.Add(notes, widget.NewButton("Later", small))
	for _, line := range []string{
		"Firmware 2.1", "Boots 2x faster", "New sleep timer",
		"Quiet fan at night", "Fewer Wi-Fi drops", "Battery % fixed",
		"Bigger clock", "Thanks, testers!",
	} {
		ui.Add(notes, widget.NewLabel(line, small, ui.Padding(ui.CurrentTheme.RowPadding)))
	}
	ui.Add(notes, widget.NewButton("Install", small))
	s := widget.NewScroll(notes, widget.ScrollVertical,
		ui.FillWidth(), ui.Padding(ui.LTRB(2, 1, 3, 1)))
	// --8<-- [end:scroll-text]
	return page("Update", s)
}

func scrollShelf() ui.Node {
	small := widget.Font(ui.CurrentTheme.SmallFont)
	// --8<-- [start:scroll-shelf]
	shelf := ui.NewRow(ui.Gap(3), ui.Padding(ui.XY(2, 1)))
	for _, m := range []string{"Blue Zenith", "Freedom Dive", "Teo", "Glory Days", "Airman"} {
		ui.Add(shelf, widget.NewButton(m, small, ui.Padding(ui.XY(6, 4))))
	}
	// The bottom padding keeps the thumb clear of the cursor.
	s := widget.NewScroll(shelf, widget.ScrollHorizontal,
		ui.FillWidth(), ui.Padding(ui.LTRB(0, 0, 0, 3)))
	// --8<-- [end:scroll-shelf]
	l := widget.NewList()
	ui.Add(l, mutedCaption("Recent"), s)
	return page("Scroll", l)
}

func scrollStagger() ui.Node {
	// --8<-- [start:scroll-stagger]
	rows := ui.NewColumn()
	for i := 1; i <= 12; i++ {
		ui.Add(rows, widget.NewListItem("Slot "+strconv.Itoa(i)))
	}
	s := widget.NewScroll(rows, widget.ScrollVertical,
		ui.FillWidth(), ui.Padding(ui.LTRB(0, 0, 3, 0)))
	s.SetStagger(0.6) // the default is 0.2; 0 moves the rows as one piece
	// --8<-- [end:scroll-stagger]
	return page("Stagger 0.6", s)
}

func scrollNested() ui.Node {
	t := ui.CurrentTheme
	small := widget.Font(t.SmallFont)
	// --8<-- [start:scroll-nested]
	// A bordered inbox, short enough to scroll on its own inside the page.
	inset := t.PanelPadding
	inset.Right += 3 // room for the thumb
	inbox := widget.NewList(ui.Height(ui.Px(34)), ui.Border(t.Foreground),
		ui.Radius(t.Radius), ui.Margin(ui.XY(2, 0)), ui.Padding(inset))
	for i := 1; i <= 5; i++ {
		ui.Add(inbox, widget.NewListItem("Message "+strconv.Itoa(i), small))
	}
	outer := widget.NewList()
	ui.Add(outer, inbox,
		widget.NewListItem("Mark all read"),
		widget.NewListItem("Settings"),
	)
	// --8<-- [end:scroll-nested]
	return page("Nested", outer)
}
