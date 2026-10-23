// Package desktop runs a lilt UI in a window, standing in for the handheld:
// a one-bit screen scaled up, and the keys on the keyboard.
//
//	W A S D, arrow keys   up, left, down, right
//	Space, Enter          OK (the joystick press)
//	J                     SET, which is Back
//	K                     RST, which is Menu
//
// It is built on Ebitengine, which needs no cgo and also builds for
// GOOS=js GOARCH=wasm, so the same program runs in a browser.
package desktop

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/chhongzh/lilt"
	"github.com/chhongzh/lilt/driver/silent"
	"github.com/chhongzh/lilt/ui"
)

// Config sizes the window. Zero fields take the defaults.
type Config struct {
	Title string // "lilt"
	// Width and Height are the screen in pixels, 128x64 like the SH1106.
	Width, Height int
	// Scale is the initial size of a screen pixel in the window.
	Scale int // 6
}

var binds = [...]struct {
	key   ui.Key
	codes []ebiten.Key
}{
	{ui.KeyUp, []ebiten.Key{ebiten.KeyW, ebiten.KeyArrowUp}},
	{ui.KeyDown, []ebiten.Key{ebiten.KeyS, ebiten.KeyArrowDown}},
	{ui.KeyLeft, []ebiten.Key{ebiten.KeyA, ebiten.KeyArrowLeft}},
	{ui.KeyRight, []ebiten.Key{ebiten.KeyD, ebiten.KeyArrowRight}},
	{ui.KeyOK, []ebiten.Key{ebiten.KeySpace, ebiten.KeyEnter}},
	{ui.KeyBack, []ebiten.Key{ebiten.KeyJ}},
	{ui.KeyMenu, []ebiten.Key{ebiten.KeyK}},
}

// Run shows root in a window and blocks until the window closes. Call it from
// main: windowing must own the main goroutine.
func Run(cfg Config, root ui.Node) error {
	if cfg.Title == "" {
		cfg.Title = "lilt"
	}
	cfg.Width = orElse(cfg.Width, 128)
	cfg.Height = orElse(cfg.Height, 64)
	cfg.Scale = orElse(cfg.Scale, 6)

	p := newPanel(cfg.Width, cfg.Height)
	g := &game{p: p, poller: ui.NewKeyPoller(), last: time.Now()}
	// Ebitengine samples the keyboard once per update, which is already clean.
	g.poller.Debounce = 0

	lilt.Start(p, root)
	ebiten.SetWindowTitle(cfg.Title)
	ebiten.SetWindowSize(cfg.Width*cfg.Scale, cfg.Height*cfg.Scale)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(g)
	if stopErr := lilt.Stop(); err == nil {
		err = stopErr
	}
	return err
}

func orElse(v, def int) int {
	if v > 0 {
		return v
	}
	return def
}

// panel is the screen lilt draws on. Drawing happens on the UI goroutine and
// the window reads on its own, so Flush publishes a finished frame into front
// under a lock instead of letting the window see a half-drawn buffer.
type panel struct {
	*silent.Screen
	w, h int

	mu    sync.Mutex
	front []byte // RGBA
}

func newPanel(w, h int) *panel {
	p := &panel{Screen: silent.New(w, h), w: w, h: h, front: make([]byte, w*h*4)}
	p.publish()
	return p
}

func (p *panel) Flush() error {
	p.mu.Lock()
	p.publish()
	p.mu.Unlock()
	return nil
}

func (p *panel) publish() {
	off, on := silent.Palette[0], silent.Palette[1]
	for i, lit := range p.Pixels() {
		c := off
		if lit {
			c = on
		}
		r, g, b, a := c.RGBA()
		p.front[4*i], p.front[4*i+1], p.front[4*i+2], p.front[4*i+3] = byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8)
	}
}

type game struct {
	p      *panel
	poller *ui.KeyPoller
	last   time.Time
	img    *ebiten.Image
}

func (g *game) Update() error {
	now := time.Now()
	dt := min(float32(now.Sub(g.last).Seconds()), 0.1)
	g.last = now
	for _, b := range binds {
		down := false
		for _, c := range b.codes {
			down = down || ebiten.IsKeyPressed(c)
		}
		g.poller.Poll(b.key, down, dt)
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.img == nil {
		g.img = ebiten.NewImage(g.p.w, g.p.h)
	}
	g.p.mu.Lock()
	g.img.WritePixels(g.p.front)
	g.p.mu.Unlock()

	screen.Fill(silent.Palette[0])
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	// Whole pixels only, so the grid stays crisp; centred in the window.
	scale := max(min(sw/g.p.w, sh/g.p.h), 1)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(float64(scale), float64(scale))
	op.GeoM.Translate(float64((sw-g.p.w*scale)/2), float64((sh-g.p.h*scale)/2))
	screen.DrawImage(g.img, op)
}

// Layout asks for one buffer pixel per device pixel, so scaling the screen up
// by whole numbers stays sharp on high-DPI displays too.
func (g *game) Layout(w, h int) (int, int) {
	f := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(w) * f), int(float64(h) * f)
}
