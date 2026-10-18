package ui

import (
	"time"

	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
)

// TextTarget is a field the system keyboard edits.
type TextTarget interface {
	Text() string
	SetText(s string)
}

// TextInputMethod opens an editor for t. The widget package provides an
// on-screen keyboard; boards with a real keyboard can supply their own.
type TextInputMethod interface{ Open(t TextTarget) }

// App owns the screen, the node tree and the frame loop. There is one per
// program; NewApp makes it the current one.
type App struct {
	Screens   ScreenStack
	Input     InputManager
	Focus     FocusManager
	Cursor    Cursor
	TextInput TextInputMethod
	// FPS caps the frame rate of Run.
	FPS int

	screen gfx.Screen
	root   *Stack
	posts  chan func()
	w, h   int16
}

var app *App

// NewApp creates the app drawing to s and makes it current.
func NewApp(s gfx.Screen) *App {
	w, h := s.Size()
	a := &App{screen: s, w: w, h: h, FPS: 60, posts: make(chan func(), 32)}
	app = a
	a.Input.init()
	a.Cursor.init()
	a.root = &Stack{}
	a.root.width, a.root.height = Px(float32(w)), Px(float32(h))
	Relayout()
	return a
}

// Current returns the app created last.
func Current() *App { return app }

func (a *App) Size() (w, h int16) { return a.w, a.h }

// Frame advances the UI by dt seconds and redraws if anything changed. It
// reports whether the screen was flushed.
func (a *App) Frame(dt float32) (bool, error) {
	for done := false; !done; {
		select {
		case fn := <-a.posts:
			fn()
		default:
			done = true
		}
	}
	if a.Input.drain() {
		Invalidate()
	}
	update(a.root, dt)
	a.Screens.update()
	if anim.Step(dt) {
		Invalidate()
	}
	if needsLayout {
		needsLayout = false
		Arrange(a.root, Rect{0, 0, float32(a.w), float32(a.h)})
		a.Screens.focusFirst()
	}
	a.Cursor.update()
	if !needsDraw {
		return false, nil
	}
	needsDraw = false
	c := gfx.NewCanvas(a.screen)
	c.FillRect(gfx.Rect{W: a.w, H: a.h}, CurrentTheme.Background)
	paint(a.root, c)
	a.Cursor.draw(c)
	return true, a.screen.Flush()
}

// Run calls Frame forever at up to FPS frames per second.
func (a *App) Run() error {
	last := time.Now()
	for {
		period := time.Second / time.Duration(a.FPS)
		now := time.Now()
		dt := min(float32(now.Sub(last).Seconds()), 0.1)
		last = now
		if _, err := a.Frame(dt); err != nil {
			return err
		}
		if spent := time.Since(now); spent < period {
			time.Sleep(period - spent)
		} else {
			// Still yield so polling goroutines get to run.
			time.Sleep(0)
		}
	}
}

func update(d Drawable, dt float32) {
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

func paint(d Drawable, c gfx.Canvas) {
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
	if dr, ok := d.(Drawer); ok {
		dr.Draw(c)
	}
	for _, ch := range b.children {
		paint(ch, c)
	}
	if o, ok := d.(OverDrawer); ok {
		o.DrawOver(c)
	}
}

// hitTest finds the deepest visible node under a screen point in the top screen.
func (a *App) hitTest(x, y int16) Drawable {
	s := a.Screens.scope()
	if s == nil {
		return nil
	}
	ox, oy := pixelOrigin(s)
	p := pixelRect(s.box())
	return hit(s, x, y, ox-p.X, oy-p.Y)
}

func hit(d Drawable, x, y, ox, oy int16) Drawable {
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

// Post runs fn on the UI goroutine at the start of the next frame. Use it to
// hand results from other goroutines back to widgets.
func Post(fn func()) { app.posts <- fn }

// Push shows d as a new screen.
func Push(d Drawable, t Transition) { app.Screens.Push(d, t) }

// Pop dismisses the top screen.
func Pop() bool { return app.Screens.Pop() }

// SetFocus moves focus to d.
func SetFocus(d Drawable) { app.Focus.Set(d) }

// Focused returns the focused node, or nil.
func Focused() Drawable { return app.Focus.current }
