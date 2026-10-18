// Command capture records the gallery for the promo video. It runs a
// scripted session at the video's 60 fps and writes every frame as a
// packed one-bit bitmap, plus a JSON index of named marks and key events
// so the video can cut and click in step with the footage.
//
// -stiff replays the same keys with anim.TimeScale turned up, so every
// spring lands within a frame or two: the "springs off" half of the
// comparison.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/examples/gallery"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"

	"tinygo.org/x/tinyfont/notosans"
)

const (
	width, height = 128, 64
	fps           = 60
	frameBytes    = width * height / 8
)

type screen struct{ px [width * height]bool }

func (s *screen) Size() (w, h int16) { return width, height }

func (s *screen) SetPixel(x, y int16, c gfx.Color) {
	if x >= 0 && y >= 0 && x < width && y < height {
		s.px[int(y)*width+int(x)] = gfx.IsLit(c)
	}
}

func (s *screen) FillRect(x, y, w, h int16, c gfx.Color) {
	on := gfx.IsLit(c)
	s.each(x, y, w, h, func(i int) { s.px[i] = on })
}

func (s *screen) InvertRect(x, y, w, h int16) {
	s.each(x, y, w, h, func(i int) { s.px[i] = !s.px[i] })
}

func (s *screen) each(x, y, w, h int16, fn func(i int)) {
	x0, y0 := max(x, 0), max(y, 0)
	x1, y1 := min(x+w, width), min(y+h, height)
	for j := y0; j < y1; j++ {
		for i := x0; i < x1; i++ {
			fn(int(j)*width + int(i))
		}
	}
}

func (s *screen) Flush() error { return nil }

// pack appends the frame MSB first, row by row.
func (s *screen) pack(out []byte) []byte {
	var b byte
	for i, on := range s.px {
		if on {
			b |= 0x80 >> (i & 7)
		}
		if i&7 == 7 {
			out = append(out, b)
			b = 0
		}
	}
	return out
}

type step struct {
	at    float32
	key   ui.Key
	phase ui.Phase
	mark  string
}

type script struct {
	t     float32
	steps []step
}

func (s *script) wait(d float32) { s.t += d }
func (s *script) mark(n string)  { s.steps = append(s.steps, step{at: s.t, mark: n}) }
func (s *script) emit(k ui.Key, p ui.Phase) {
	s.steps = append(s.steps, step{at: s.t, key: k, phase: p})
}

func (s *script) tap(k ui.Key, gap float32) {
	s.emit(k, ui.Press)
	s.wait(0.09)
	s.emit(k, ui.Release)
	s.wait(gap)
}

func (s *script) taps(k ui.Key, n int, gap float32) {
	for range n {
		s.tap(k, gap)
	}
}

// reel walks through the moments the video features. Each mark opens a
// segment the video cuts to; the gaps before marks let the screen settle.
func reel() *script {
	s := &script{}
	s.wait(0.5)

	s.mark("cursor")
	s.taps(ui.KeyDown, 4, 0.32)
	s.wait(0.15)
	s.taps(ui.KeyUp, 4, 0.32)
	s.wait(0.4)

	s.mark("push")
	s.tap(ui.KeyOK, 0.7)
	s.mark("press")
	s.emit(ui.KeyOK, ui.Press)
	s.wait(0.35)
	s.emit(ui.KeyOK, ui.Release)
	s.wait(0.7)
	s.emit(ui.KeyOK, ui.Press)
	s.wait(0.3)
	s.emit(ui.KeyOK, ui.Release)
	s.wait(0.6)

	s.mark("hold")
	s.tap(ui.KeyDown, 0.4)
	s.emit(ui.KeyOK, ui.Press)
	s.wait(1.05)
	s.emit(ui.KeyOK, ui.Release)
	s.wait(0.8)

	s.mark("action")
	s.tap(ui.KeyDown, 0.35)
	s.tap(ui.KeyOK, 1.7)
	s.tap(ui.KeyDown, 0.35)
	s.tap(ui.KeyOK, 1.7)

	s.mark("back")
	s.emit(ui.KeyBack, ui.Press)
	s.wait(0.45)
	s.emit(ui.KeyBack, ui.Release)
	s.wait(0.8)

	s.taps(ui.KeyDown, 2, 0.25)
	s.tap(ui.KeyOK, 0.6)
	s.mark("sheet")
	s.tap(ui.KeyOK, 0.8)
	s.taps(ui.KeyDown, 2, 0.3)
	s.tap(ui.KeyOK, 0.8)
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.25)
	s.tap(ui.KeyOK, 0.6)
	s.mark("chips")
	s.taps(ui.KeyRight, 2, 0.3)
	s.tap(ui.KeyOK, 0.4)
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.25)
	s.taps(ui.KeyRight, 2, 0.25)
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.25)
	s.tap(ui.KeyOK, 0.6)
	s.mark("flex")
	s.tap(ui.KeyRight, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.taps(ui.KeyRight, 2, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.8)
	s.tap(ui.KeyBack, 0.6)

	s.taps(ui.KeyDown, 2, 0.25)
	s.tap(ui.KeyOK, 0.6)
	s.mark("shelves")
	s.taps(ui.KeyRight, 4, 0.25)
	s.tap(ui.KeyDown, 0.3)
	s.taps(ui.KeyRight, 2, 0.25)
	s.wait(0.4)
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.25)
	s.tap(ui.KeyOK, 0.6)
	s.mark("keyboard")
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyOK, 0.25)
	s.tap(ui.KeyDown, 0.25)
	s.tap(ui.KeyOK, 0.25)
	s.taps(ui.KeyRight, 2, 0.25)
	s.tap(ui.KeyOK, 0.4)
	s.taps(ui.KeyDown, 2, 0.3)
	s.taps(ui.KeyRight, 2, 0.3)
	s.tap(ui.KeyOK, 0.8)
	s.tap(ui.KeyBack, 0.6)

	s.taps(ui.KeyDown, 2, 0.25)
	s.tap(ui.KeyOK, 0.9)
	s.mark("wave")
	s.emit(ui.KeyDown, ui.Press)
	for range 34 {
		s.wait(0.07)
		s.emit(ui.KeyDown, ui.Repeat)
	}
	s.wait(0.55)
	s.emit(ui.KeyDown, ui.Release)
	s.wait(0.9)
	s.mark("end")
	return s
}

// compare is the short run recorded twice, with and without springs.
func compare() *script {
	s := &script{}
	s.wait(0.4)
	s.mark("compare")
	s.taps(ui.KeyDown, 3, 0.3)
	s.tap(ui.KeyUp, 0.3)
	s.tap(ui.KeyOK, 0.75)
	s.taps(ui.KeyDown, 2, 0.3)
	s.tap(ui.KeyBack, 0.7)
	s.mark("end")
	return s
}

// slow moves the cursor on the home list for the 4x slow-motion close-up.
// Keys are held and spaced four times longer to match.
func slow() *script {
	s := &script{}
	s.wait(0.5)
	s.mark("slow")
	hold := func(k ui.Key, gap float32) {
		s.emit(k, ui.Press)
		s.wait(0.36)
		s.emit(k, ui.Release)
		s.wait(gap)
	}
	hold(ui.KeyDown, 1.5)
	hold(ui.KeyDown, 1.5)
	hold(ui.KeyUp, 1.6)
	s.mark("end")
	return s
}

type keyEvent struct {
	Frame int    `json:"frame"`
	Key   string `json:"key"`
	Phase string `json:"phase"`
}

type index struct {
	FPS    int            `json:"fps"`
	Width  int            `json:"width"`
	Height int            `json:"height"`
	Frames int            `json:"frames"`
	Marks  map[string]int `json:"marks"`
	Keys   []keyEvent     `json:"keys"`
}

var keyNames = map[ui.Key]string{
	ui.KeyUp: "up", ui.KeyDown: "down", ui.KeyLeft: "left", ui.KeyRight: "right",
	ui.KeyOK: "ok", ui.KeyBack: "back",
}

var phaseNames = map[ui.Phase]string{ui.Press: "press", ui.Repeat: "repeat", ui.Release: "release"}

func main() {
	out := flag.String("o", "public/footage", "output directory")
	which := flag.String("script", "reel", "reel, compare or slow")
	stiff := flag.Bool("stiff", false, "turn springs off by running animation 40x faster")
	logo := flag.Bool("logo", false, "write the lilt wordmark as drawn by the on-device font")
	flag.Parse()
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fail(err)
	}
	if *logo {
		writeLogo(*out)
		return
	}

	scripts := map[string]func() *script{"reel": reel, "compare": compare, "slow": slow}
	build, ok := scripts[*which]
	if !ok {
		fail(fmt.Errorf("unknown script %q", *which))
	}
	sc := build()
	name := *which
	if name == "slow" {
		anim.TimeScale = 0.25
	}
	if *stiff {
		name += "-stiff"
		anim.TimeScale = 40
	}

	scr := &screen{}
	app := ui.NewApp(scr)
	gallery.Build()

	idx := index{FPS: fps, Width: width, Height: height, Marks: map[string]int{}}
	data := make([]byte, 0, frameBytes*fps*60)
	dt := float32(1) / fps
	start := time.Now()
	next := 0
	for frame := 0; ; frame++ {
		t := float32(frame) * dt
		if next == len(sc.steps) {
			break
		}
		for next < len(sc.steps) && sc.steps[next].at <= t {
			st := sc.steps[next]
			if st.mark != "" {
				idx.Marks[st.mark] = frame
			} else {
				app.Input.Emit(st.key, st.phase)
				idx.Keys = append(idx.Keys, keyEvent{frame, keyNames[st.key], phaseNames[st.phase]})
			}
			next++
		}
		if _, err := app.Frame(dt); err != nil {
			fail(err)
		}
		data = scr.pack(data)
		idx.Frames++
		// Background actions sleep on the wall clock, so pace to it.
		if ahead := time.Duration(t*float32(time.Second)) - time.Since(start); ahead > 0 {
			time.Sleep(ahead)
		}
	}

	write(filepath.Join(*out, name+".bin"), data)
	js, err := json.Marshal(idx)
	if err != nil {
		fail(err)
	}
	write(filepath.Join(*out, name+".json"), js)
	fmt.Printf("%s: %d frames, marks %v\n", name, idx.Frames, idx.Marks)
}

// writeLogo draws "lilt" with the framework's own font and stores the lit
// pixels, cropped, so the video can assemble the wordmark pixel by pixel.
func writeLogo(dir string) {
	scr := &screen{}
	c := gfx.NewCanvas(scr)
	f := &notosans.Notosans12pt
	w := gfx.TextWidth(f, "lilt")
	c.Text(f, (width-w)/2, 40, "lilt", gfx.White)

	x0, y0, x1, y1 := width, height, 0, 0
	for i, on := range scr.px {
		if on {
			x, y := i%width, i/width
			x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x+1), max(y1, y+1)
		}
	}
	type logo struct {
		W, H   int
		Pixels [][2]int
	}
	l := logo{W: x1 - x0, H: y1 - y0}
	for i, on := range scr.px {
		if on {
			l.Pixels = append(l.Pixels, [2]int{i%width - x0, i/width - y0})
		}
	}
	js, err := json.Marshal(l)
	if err != nil {
		fail(err)
	}
	write(filepath.Join(dir, "logo.json"), js)
	fmt.Printf("logo: %dx%d, %d pixels\n", l.W, l.H, len(l.Pixels))
}

func write(path string, b []byte) {
	if err := os.WriteFile(path, b, 0o644); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "capture:", err)
	os.Exit(1)
}
