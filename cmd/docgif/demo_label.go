package main

import (
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "label-styles",
		build: labelStyles,
		script: func(s *reel.Script) {
			s.Wait(1)
			s.Shot("styles")
			s.Wait(1)
		},
	})
	register(demo{
		name:  "label-marquee",
		build: labelMarquee,
		script: func(s *reel.Script) {
			s.Wait(2.6)
			s.Shot("focused row scrolls")
			s.Tap(ui.KeyDown, 0.4)
			s.Shot("first row back at start")
			s.Wait(2.2)
			s.Shot("second row scrolls")
			s.Wait(0.5)
		},
	})
}

func labelStyles() ui.Node {
	// --8<-- [start:label-styles]
	t := ui.CurrentTheme
	small := widget.Font(t.SmallFont)
	title := widget.NewLabel("Blue Zenith", widget.TextAlign(0.5), ui.FillWidth())
	artist := widget.NewLabel("xi", small, widget.TextColor(t.Muted),
		widget.TextAlign(0.5), ui.FillWidth())
	// Too long and right-aligned: it rests showing its end.
	file := widget.NewLabel("sd:/songs/blue-zenith/audio.mp3", small,
		widget.TextAlign(1), ui.FillWidth(), ui.Padding(t.RowPadding))
	play := widget.NewListItem("Play", small,
		widget.Trailing(widget.NewLabel("6:12", small, widget.TextColor(t.Muted))))
	// --8<-- [end:label-styles]
	body := ui.NewColumn(ui.Children(title, artist, file, play))
	return page("Label", body)
}

func labelMarquee() ui.Node {
	// --8<-- [start:label-marquee]
	t := ui.CurrentTheme
	news := widget.NewLabel("Firmware 2.4 is ready: faster boot, two fixes",
		widget.Font(t.SmallFont), widget.MarqueeMode(widget.MarqueeAlways),
		ui.Padding(t.RowPadding))
	list := widget.NewList()
	ui.Add(list,
		news,
		widget.NewListItem("Harumachi Clover (Swing)"),
		widget.NewListItem("Kani Do-Raku (Extended)"),
	)
	// --8<-- [end:label-marquee]
	return page("Label", list)
}
