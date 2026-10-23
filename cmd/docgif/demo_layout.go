package main

import (
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "layout-anchor",
		build: layoutAnchor,
		script: func(s *reel.Script) {
			s.Wait(0.6)
			s.Tap(ui.KeyOK, 0.7)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("top-right")
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("bottom-right")
			s.Tap(ui.KeyOK, 0.7)
			s.Tap(ui.KeyOK, 0.7)
			s.Tap(ui.KeyOK, 0.8)
		},
	})
	register(demo{
		name:  "layout-flex",
		build: layoutFlex,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			for _, shot := range []string{"center", "end", "between"} {
				s.Tap(ui.KeyRight, 0.3)
				s.Tap(ui.KeyOK, 0.7)
				s.Shot(shot)
			}
			s.Taps(ui.KeyLeft, 3, 0.25)
			s.Tap(ui.KeyOK, 0.8)
		},
	})
	register(demo{
		name:  "layout-grow",
		build: layoutGrow,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyRight, 0.3)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("one-grows")
			s.Tap(ui.KeyRight, 0.3)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("two-share")
			s.Tap(ui.KeyOK, 0.4)
			s.Tap(ui.KeyLeft, 0.3)
			s.Tap(ui.KeyOK, 0.8)
		},
	})
	register(demo{
		name:  "layout-grid",
		build: layoutGrid,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyRight, 0.5)
			s.Shot("right")
			s.Taps(ui.KeyDown, 2, 0.5)
			s.Shot("down")
			s.Tap(ui.KeyLeft, 0.5)
			s.Taps(ui.KeyUp, 2, 0.6)
		},
	})
	register(demo{
		name:  "layout-visible",
		build: layoutVisible,
		script: func(s *reel.Script) {
			s.Wait(0.6)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("hidden")
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("shown")
		},
	})
}

func layoutAnchor() ui.Node {
	// --8<-- [start:layout-anchor]
	t := ui.CurrentTheme
	spots := []ui.Vec2{ui.TopLeft, ui.TopCenter, ui.TopRight, ui.BottomRight, ui.BottomCenter, ui.BottomLeft}
	names := []string{"TopLeft", "TopCenter", "TopRight", "BottomRight", "BottomCenter", "BottomLeft"}
	tag := widget.NewLabel(names[0], widget.Font(t.SmallFont),
		ui.Border(t.Foreground), ui.Radius(t.Radius), ui.Padding(ui.XY(3, 0)),
		ui.AnimateLayout()) // glide to the new spot instead of jumping
	at := 0
	move := widget.NewButton("Move", ui.Align(ui.Center), widget.OnClick(func() {
		ui.Do(func() {
			at = (at + 1) % len(spots)
			tag.SetAnchor(spots[at]) // the point of the parent...
			tag.SetOrigin(spots[at]) // ...and the point of the tag placed on it
			tag.SetText(names[at])
		})
	}))
	stage := ui.NewStack(ui.Fill(), ui.Padding(ui.All(2)), ui.Children(move, tag))
	// --8<-- [end:layout-anchor]
	return stage
}

func layoutFlex() ui.Node {
	// --8<-- [start:layout-flex]
	t := ui.CurrentTheme
	small := widget.Font(t.SmallFont)
	row := ui.NewRow(ui.Gap(2), ui.Padding(ui.XY(4, 2)))
	for _, s := range []string{"1", "2", "3"} {
		ui.Add(row, widget.NewLabel(s, small, ui.Border(t.Foreground), ui.Radius(t.Radius),
			ui.Padding(ui.XY(5, 1)), ui.AnimateLayout()))
	}
	aligns := []ui.MainAlign{ui.MainStart, ui.MainCenter, ui.MainEnd, ui.MainSpaceBetween}
	align := widget.NewChoice([]string{"Start", "Center", "End", "Between"}, 0, small,
		widget.OnChange(func(from, to int) {
			ui.Do(func() { row.SetMainAlign(aligns[to]) })
		}))
	// --8<-- [end:layout-flex]
	return page("Flex", ui.NewColumn(ui.Children(row, align)))
}

func layoutGrow() ui.Node {
	// --8<-- [start:layout-grow]
	t := ui.CurrentTheme
	small := widget.Font(t.SmallFont)
	row := ui.NewRow(ui.Gap(2), ui.Padding(ui.XY(4, 2)))
	boxes := make([]*widget.Label, 3)
	for i, s := range []string{"A", "B", "C"} {
		boxes[i] = widget.NewLabel(s, small, widget.TextAlign(0.5), ui.Border(t.Foreground),
			ui.Radius(t.Radius), ui.Padding(ui.XY(5, 1)), ui.AnimateLayout())
		ui.Add(row, boxes[i])
	}
	grow := widget.NewToggles([]string{"Grow A", "Grow B", "Grow C"}, small,
		widget.OnChange(func(from, to uint32) {
			ui.Do(func() {
				for i, b := range boxes {
					b.SetGrow(float32(to >> i & 1)) // bit i is chip i
				}
			})
		}))
	// --8<-- [end:layout-grow]
	return page("Grow", ui.NewColumn(ui.Children(row, grow)))
}

func layoutGrid() ui.Node {
	// --8<-- [start:layout-grid]
	small := widget.Font(ui.CurrentTheme.SmallFont)
	grid := ui.NewGrid(ui.Columns(ui.Fixed(34), ui.Fr(1), ui.Fr(1)), ui.Gap(2),
		ui.Padding(ui.XY(2, 1)))
	for _, room := range []string{"Hall", "Desk", "Porch"} {
		ui.Add(grid,
			widget.NewLabel(room, small),
			widget.NewButton("On", small),
			widget.NewButton("Off", small),
		)
	}
	// --8<-- [end:layout-grid]
	return page("Lights", grid)
}

func layoutVisible() ui.Node {
	// --8<-- [start:layout-visible]
	small := widget.Font(ui.CurrentTheme.SmallFont)
	network := widget.NewListItem("Network",
		widget.Trailing(widget.NewLabel("Home", small)))
	wifi := widget.NewCheckBox("Wi-Fi", true, widget.OnChange(func(from, to bool) {
		ui.Do(func() { network.SetVisible(to) }) // hidden rows take no space
	}))
	bluetooth := widget.NewCheckBox("Bluetooth", false, ui.AnimateLayout())
	list := widget.NewList(ui.Children(wifi, network, bluetooth))
	// --8<-- [end:layout-visible]
	return page("Wireless", list)
}
