package main

import (
	"strconv"

	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "tree-ref",
		build: treeRef,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.4)
			s.Tap(ui.KeyOK, 0.5)
			s.Shot("louder")
			s.Tap(ui.KeyDown, 0.4)
			s.Taps(ui.KeyOK, 3, 0.4)
			s.Wait(0.2)
			s.Shot("quieter")
		},
	})
	register(demo{
		name:  "tree-add-remove",
		build: treeAddRemove,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Tap(ui.KeyOK, 0.5)
			s.Tap(ui.KeyOK, 0.6)
			s.Shot("added")
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyDown, 0.5)
			s.Shot("picked")
			s.Tap(ui.KeyOK, 0.8)
			s.Shot("removed")
		},
	})
}

func treeRef() ui.Node {
	// --8<-- [start:tree-ref]
	level := 5
	var volume *widget.Label // filled in by ui.Ref while the tree is built
	screen := ui.NewColumn(ui.Fill(), ui.Children(
		ui.NewRow(ui.Padding(ui.CurrentTheme.RowPadding), ui.Children(
			widget.NewLabel("Speaker", ui.Grow(1)),
			widget.NewLabel("Vol 5", ui.Ref(&volume)),
		)),
		widget.NewList(ui.Grow(1), ui.Children(
			widget.NewListItem("Louder", widget.OnClick(func() {
				ui.Do(func() {
					level++
					volume.SetText("Vol " + strconv.Itoa(level))
				})
			})),
			widget.NewListItem("Quieter", widget.OnClick(func() {
				ui.Do(func() {
					level--
					volume.SetText("Vol " + strconv.Itoa(level))
				})
			})),
		)),
	))
	// --8<-- [end:tree-ref]
	return screen
}

func treeAddRemove() ui.Node {
	// --8<-- [start:tree-add-remove]
	list := widget.NewList()
	added := 0
	ui.Add(list, widget.NewListItem("+ New alarm", widget.OnClick(func() {
		ui.Do(func() {
			added++
			var row *widget.ListItem
			row = widget.NewListItem("Alarm "+strconv.Itoa(added), widget.OnClick(func() {
				ui.Do(func() { ui.Remove(row) }) // a row may remove itself
			}))
			ui.Insert(list, 1, row) // right below "+ New alarm"
		})
	})))
	// --8<-- [end:tree-add-remove]
	return page("Alarms", list)
}
