package main

import (
	"time"

	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "popups-drawer",
		build: popupsDrawer,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("drawer up")
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.8)
			s.Shot("closed, toast")
			s.Wait(1.5)
		},
	})
	register(demo{
		name:  "popups-side",
		build: popupsSide,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("drawer in")
			s.Tap(ui.KeyDown, 0.4)
			s.Emit(ui.KeyBack, ui.Press)
			s.Wait(0.3)
			s.Shot("back held: peek")
			s.Emit(ui.KeyBack, ui.Release)
			s.Wait(0.7)
			s.Shot("closed")
			s.Wait(1)
		},
	})
	register(demo{
		name:  "popups-dialog",
		build: popupsDialog,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.7)
			s.Shot("dialog")
			s.Tap(ui.KeyRight, 0.4)
			s.Tap(ui.KeyOK, 0.4)
			s.Shot("deleting")
			s.Wait(0.8)
			s.Shot("closed, toast")
			s.Wait(1.5)
		},
	})
	register(demo{
		name:  "popups-toast",
		build: popupsToast,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("toast")
			s.Tap(ui.KeyDown, 0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("replaced")
			s.Wait(2.2)
			s.Shot("gone")
		},
	})
}

func popupsDrawer() ui.Node {
	// --8<-- [start:popups-drawer]
	share := widget.NewListItem("Share score", widget.OnClick(func() {
		ui.Do(func() {
			d := widget.NewDrawer(widget.Bottom, "Send to")
			for _, name := range []string{"Phone", "Laptop"} {
				ui.Add(d, widget.NewListItem(name, widget.OnClick(func() {
					ui.Do(func() {
						d.Close()
						widget.Toast("Sent to " + name)
					})
				})))
			}
			d.Open()
		})
	}))
	// --8<-- [end:popups-drawer]
	l := widget.NewList()
	ui.Add(l, share, widget.NewListItem("Replay"))
	return page("Drawer", l)
}

func popupsSide() ui.Node {
	// --8<-- [start:popups-side]
	menu := widget.NewDrawer(widget.Left, "Menu")
	ui.Add(menu,
		widget.NewListItem("Home"),
		widget.NewListItem("Maps"),
		widget.NewListItem("Settings"),
	)
	open := widget.NewListItem("Open menu", widget.OnClick(func() {
		ui.Do(menu.Open)
	}))
	// --8<-- [end:popups-side]
	l := widget.NewList()
	ui.Add(l, open, widget.NewListItem("Now playing"))
	return page("Drawer", l)
}

func popupsDialog() ui.Node {
	// --8<-- [start:popups-dialog]
	del := widget.NewListItem("Delete map", widget.OnClick(func() {
		ui.Do(func() {
			widget.NewDialog("Delete map?", "Blue Zenith goes for good, scores and all",
				widget.NewButton("Keep"),
				widget.NewButton("Delete", widget.OnClick(func() {
					time.Sleep(700 * time.Millisecond) // removing the files
					ui.Do(func() { widget.Toast("Deleted") })
				})),
			).Open()
		})
	}))
	// --8<-- [end:popups-dialog]
	l := widget.NewList()
	ui.Add(l, del, widget.NewListItem("Export"))
	return page("Dialog", l)
}

func popupsToast() ui.Node {
	// --8<-- [start:popups-toast]
	save := widget.NewListItem("Save to slot 2", widget.OnClick(func() {
		ui.Do(func() { widget.Toast("Saved to slot 2") })
	}))
	load := widget.NewListItem("Load slot 2", widget.OnClick(func() {
		ui.Do(func() { widget.Toast("Loaded slot 2") })
	}))
	// --8<-- [end:popups-toast]
	l := widget.NewList()
	ui.Add(l, save, load)
	return page("Toast", l)
}
