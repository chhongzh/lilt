package ui

import (
	"sync"

	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/internal/core"
)

// TextTarget is a field the system keyboard edits.
type TextTarget interface {
	Text() string
	SetText(s string)
}

// TextInputMethod opens an editor for t. The widget package provides an
// on-screen keyboard; boards with a real keyboard can supply their own.
type TextInputMethod interface{ Open(t TextTarget) }

// TextInput is the editor Input opens; nil means the widget package's
// soft keyboard.
var TextInput TextInputMethod

// The UI is a singleton: one screen, one tree, one of each manager.
var (
	screen     gfx.Screen
	scrW, scrH int16
	root       = &Stack{}

	layers layerStack
	focus  focusManager
	input  inputManager
	cursor Cursor

	// frameMu is held for the whole of each frame; Do takes it between frames.
	frameMu sync.Mutex
)

func init() {
	input.init()
	cursor.init()
	core.Boot = boot
	core.Frame = frame
	core.Reset = reset
}

func boot(s gfx.Screen, first any) {
	screen = s
	scrW, scrH = s.Size()
	root.width, root.height = Px(float32(scrW)), Px(float32(scrH))
	Relayout()
	if n, ok := first.(Node); ok && n != nil {
		Push(n)
	}
}

func reset() {
	frameMu.Lock()
	defer frameMu.Unlock()
	for len(layers.list) > 0 {
		h := layers.list[len(layers.list)-1]
		layers.list = layers.list[:len(layers.list)-1]
		Remove(h)
	}
	focus.Set(nil)
	input.reset()
	cursor.has = false
}

// Do runs fn between frames, so code on other goroutines can change the
// tree. Event callbacks run off the UI goroutine and must wrap their UI
// changes in Do. fn must not call Do.
func Do(fn func()) {
	frameMu.Lock()
	defer frameMu.Unlock()
	fn()
}

// ScreenSize is the size of the bound screen in pixels.
func ScreenSize() (w, h int16) { return scrW, scrH }

func frame(dt float32) (bool, error) {
	frameMu.Lock()
	defer frameMu.Unlock()
	if input.drain() {
		Invalidate()
	}
	update(root, dt)
	layers.update()
	if anim.Step(dt) {
		Invalidate()
	}
	if needsLayout {
		needsLayout = false
		Arrange(root, Rect{0, 0, float32(scrW), float32(scrH)})
		layers.focusFirst()
	}
	cursor.update()
	if !needsDraw {
		return false, nil
	}
	needsDraw = false
	c := gfx.NewCanvas(screen)
	c.FillRect(gfx.Rect{W: scrW, H: scrH}, CurrentTheme.Background)
	paint(root, c)
	cursor.draw(c)
	return true, screen.Flush()
}

func update(d Node, dt float32) {
	b := d.box()
	if !b.visible() {
		return
	}
	if u, ok := d.(Updater); ok {
		u.Update(dt)
	}
	for _, c := range b.children {
		update(c, dt)
	}
}

func paint(d Node, c gfx.Canvas) {
	b := d.box()
	if !b.visible() {
		return
	}
	p := pixelRect(b)
	c = c.Translate(p.X, p.Y)
	local := gfx.Rect{W: p.W, H: p.H}
	if b.flags&flagClip != 0 {
		if !c.Visible(local) {
			return
		}
		c = c.Clip(local)
	}
	if b.flags&flagBackground != 0 {
		c.FillRoundRect(local, b.radius, b.bg)
	}
	if b.flags&flagBorder != 0 {
		c.StrokeRoundRect(local, b.radius, b.border)
	}
	if dr, ok := d.(Painter); ok {
		dr.Draw(c)
	}
	for _, ch := range b.children {
		paint(ch, c)
	}
	if o, ok := d.(OverPainter); ok {
		o.DrawOver(c)
	}
}

// hitTest finds the deepest visible node under a screen point in the top screen.
func hitTest(x, y int16) Node {
	s := layers.scope()
	if s == nil {
		return nil
	}
	ox, oy := pixelOrigin(s)
	p := pixelRect(s.box())
	return hit(s, x, y, ox-p.X, oy-p.Y)
}

func hit(d Node, x, y, ox, oy int16) Node {
	b := d.box()
	if !b.visible() {
		return nil
	}
	r := pixelRect(b).Offset(ox, oy)
	in := r.Contains(x, y)
	if b.flags&flagClip != 0 && !in {
		return nil
	}
	for i := len(b.children) - 1; i >= 0; i-- {
		if h := hit(b.children[i], x, y, r.X, r.Y); h != nil {
			return h
		}
	}
	if in {
		return d
	}
	return nil
}

// SetFocus moves focus to d.
func SetFocus(d Node) { focus.Set(d) }

// Focused returns the focused node, or nil.
func Focused() Node { return focus.current }

// Emit queues a key event. Safe to call from any goroutine.
func Emit(k Key, p Phase) { input.Emit(k, p) }

// EmitRune queues a typed character. Safe to call from any goroutine.
func EmitRune(r rune) { input.EmitRune(r) }

// EmitPointer queues a touch or mouse event in screen pixels. Safe to call
// from any goroutine.
func EmitPointer(x, y int16, p PointerPhase) { input.EmitPointer(x, y, p) }
