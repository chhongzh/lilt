package main

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

func init() {
	register(demo{
		name:  "animation-presets",
		build: animationPresets,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.1)
			s.Shot("under way")
			s.Wait(1.1)
			s.Shot("settled right")
			s.Tap(ui.KeyOK, 1.3)
			s.Shot("settled left")
		},
	})
	register(demo{
		name:  "animation-offset",
		build: animationOffset,
		script: func(s *reel.Script) {
			s.Wait(0.5)
			s.Tap(ui.KeyOK, 0.07)
			s.Shot("kicked")
			s.Wait(0.8)
			s.Tap(ui.KeyDown, 0.4)
			s.Tap(ui.KeyOK, 0.1)
			s.Shot("bumped")
			s.Wait(0.9)
		},
	})
}

// --8<-- [start:animation-lane]
// lane draws a dot that springs between the ends of a line.
type lane struct {
	ui.Box
	pos anim.Value // 0 at the left end, 1 at the right
}

func newLane(s *anim.Spring) *lane {
	l := &lane{}
	l.pos.SetSpring(s)
	l.pos.SetPrecision(0.002) // the default 0.05 suits pixels, not 0..1
	l.SetGrow(1)
	return l
}

func (l *lane) Leaf() {}

func (l *lane) Measure(max ui.Vec2) ui.Vec2 { return ui.Vec2{X: max.X, Y: 7} }

func (l *lane) Draw(c gfx.Canvas) {
	r := l.Rect().Pixels()
	col := ui.CurrentTheme.Foreground
	left, w, y := int16(6), r.W-16, r.H/2
	c.VLine(left, y-2, 5, col)
	c.VLine(left+w, y-2, 5, col)
	c.FillCircle(left+ui.Round(l.pos.Get()*float32(w)), y, 2, col)
}

// --8<-- [end:animation-lane]

func animationPresets() ui.Node {
	// --8<-- [start:animation-presets]
	t := ui.CurrentTheme
	rows := ui.NewColumn(ui.FillWidth(), ui.Padding(ui.XY(2, 0)))
	presets := []struct {
		name   string
		spring *anim.Spring
	}{
		{"Default", anim.Default},
		{"Smooth", anim.Smooth},
		{"Snappy", anim.Snappy},
		{"Bouncy", anim.Bouncy},
		{"Interactive", anim.Interactive},
	}
	var lanes []*lane
	for _, p := range presets {
		l := newLane(p.spring)
		lanes = append(lanes, l)
		name := widget.NewLabel(p.name, widget.Font(t.SmallFont), ui.Width(ui.Px(64)))
		ui.Add(rows, ui.NewRow(ui.Cross(ui.CrossCenter), ui.Children(name, l)))
	}
	right := false
	swap := widget.NewButton("Swap", widget.Font(t.SmallFont), ui.Padding(ui.XY(4, 0)),
		ui.Align(ui.BottomRight), widget.OnClick(func() {
			ui.Do(func() {
				right = !right
				for _, l := range lanes {
					if right {
						l.pos.Set(1) // the same target, five springs
					} else {
						l.pos.Set(0)
					}
				}
			})
		}))
	// --8<-- [end:animation-presets]
	return ui.NewStack(ui.Fill(), ui.Children(rows, swap))
}

func animationOffset() ui.Node {
	// --8<-- [start:animation-offset]
	pin := widget.NewLabel("Enter PIN", ui.Align(ui.Center))
	kick := widget.NewListItem("Kick sideways", widget.OnClick(func() {
		ui.Do(func() { pin.OffsetX().Kick(-500) }) // px per second
	}))
	bump := widget.NewListItem("Bump down", widget.OnClick(func() {
		ui.Do(func() { pin.OffsetY().Bump(8) }) // about 8px away
	}))
	// --8<-- [end:animation-offset]
	list := widget.NewList(ui.Grow(1))
	ui.Add(list, kick, bump)
	top := ui.NewStack(ui.FillWidth(), ui.Height(ui.Px(24)), ui.Children(pin))
	return ui.NewColumn(ui.Fill(), ui.Children(top, list))
}
