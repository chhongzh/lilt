package main

import (
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "layers-push",
		build: layersPush,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.15)
			s.Shot("sliding in")
			s.Wait(0.5)
			s.Shot("display page")
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("brightness page")
			s.Tap(ui.KeyBack, 0.7)
			s.Tap(ui.KeyBack, 0.15)
			s.Shot("sliding out")
			s.Wait(0.6)
		},
	})
	register(demo{
		name:  "layers-popups",
		build: layersPopups,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("drawer open")
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("closed, toast")
			s.Wait(1)
		},
	})
	register(demo{
		name:  "layers-back",
		build: layersBack,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.7)
			s.Emit(ui.KeyBack, ui.Press)
			s.Wait(0.4)
			s.Shot("back held: peek")
			s.Emit(ui.KeyBack, ui.Release)
			s.Wait(0.7)
			s.Emit(ui.KeyBack, ui.Press)
			s.Wait(0.4)
			s.Shot("root page: half peek")
			s.Emit(ui.KeyBack, ui.Release)
			s.Wait(0.12)
			s.Shot("root page: nudge")
			s.Wait(0.6)
		},
	})
}

func layersPush() ui.Node {
	// --8<-- [start:layers-push]
	brightness := func() ui.Node {
		l := widget.NewList()
		g := widget.NewRadioGroup(1)
		ui.Add(l, widget.NewRadio(g, "Low"), widget.NewRadio(g, "Mid"), widget.NewRadio(g, "High"))
		return l
	}
	display := func() ui.Node {
		l := widget.NewList()
		ui.Add(l,
			widget.NewListItem("Brightness", widget.OnClick(func() {
				ui.Do(func() { ui.Push(brightness()) })
			})),
			widget.NewCheckBox("Dim at night", true),
		)
		return l
	}
	home := widget.NewList()
	ui.Add(home,
		widget.NewListItem("Display", widget.OnClick(func() {
			ui.Do(func() { ui.Push(display()) })
		})),
		widget.NewListItem("Sound"),
	)
	// --8<-- [end:layers-push]
	return page("Settings", home)
}

func layersPopups() ui.Node {
	// --8<-- [start:layers-popups]
	actions := func() ui.Node {
		menu := ui.NewColumn()
		for _, s := range []string{"Rename", "Share", "Delete"} {
			ui.Add(menu, widget.NewListItem(s, widget.OnClick(func() {
				ui.Do(func() {
					ui.Close(menu)
					widget.Toast(s + " picked")
				})
			})))
		}
		return menu
	}
	open := widget.NewListItem("Actions", widget.OnClick(func() {
		ui.Do(func() { ui.Open(actions(), ui.DrawerRight) }) // a new menu each time
	}))
	// --8<-- [end:layers-popups]
	list := widget.NewList()
	ui.Add(list, open, widget.NewListItem("Details"))
	return page("Map", list)
}

func layersBack() ui.Node {
	list := widget.NewList()
	ui.Add(list,
		widget.NewListItem("Recent maps", widget.OnClick(func() {
			ui.Do(func() {
				detail := widget.NewList()
				ui.Add(detail, widget.NewListItem("Blue Zenith"), widget.NewListItem("Freedom Dive"))
				ui.Push(page("Recent", detail))
			})
		})),
		widget.NewListItem("Favourites"),
	)
	return page("Library", list)
}
