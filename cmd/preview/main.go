// Command preview runs the gallery on the host against a scripted key
// timeline and writes an animated GIF of the session plus a PNG contact
// sheet of marked moments. It runs in real time so background actions,
// which sleep on the wall clock, finish in step with the frames.
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/chhongzh/lilt/examples/gallery"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

const (
	width, height = 128, 64
	fps           = 120
	// GIF frames are sampled at half the frame rate to keep the file small.
	gifEvery = 2
)

var palette = color.Palette{
	color.RGBA{8, 10, 14, 255},
	color.RGBA{200, 230, 255, 255},
}

// screen is a one-bit frame buffer, like the OLED it stands in for.
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

func (s *screen) image(scale int) *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, width*scale, height*scale), palette)
	for y := 0; y < height*scale; y++ {
		for x := 0; x < width*scale; x++ {
			if s.px[(y/scale)*width+x/scale] {
				img.Pix[y*img.Stride+x] = 1
			}
		}
	}
	return img
}

type step struct {
	at    float32
	key   ui.Key
	phase ui.Phase
	shot  string
}

type script struct {
	t     float32
	steps []step
}

func (s *script) wait(d float32) { s.t += d }

func (s *script) emit(k ui.Key, p ui.Phase) {
	s.steps = append(s.steps, step{at: s.t, key: k, phase: p})
}

func (s *script) shot(name string) { s.steps = append(s.steps, step{at: s.t, shot: name}) }

func (s *script) tap(k ui.Key, gap float32) {
	s.emit(k, ui.Press)
	s.wait(0.08)
	s.emit(k, ui.Release)
	s.wait(gap)
}

func (s *script) taps(k ui.Key, n int, gap float32) {
	for range n {
		s.tap(k, gap)
	}
}

func timeline() *script {
	s := &script{}
	s.wait(0.6)
	s.shot("home")

	s.tap(ui.KeyOK, 0.6)
	s.shot("buttons")
	s.tap(ui.KeyOK, 0.3)
	s.shot("tapped")
	s.tap(ui.KeyDown, 0.4)
	s.emit(ui.KeyOK, ui.Press)
	s.wait(0.4)
	s.shot("holding")
	s.wait(0.5)
	s.shot("hold fired")
	s.emit(ui.KeyOK, ui.Release)
	s.wait(1.6)
	s.tap(ui.KeyDown, 0.4)
	s.tap(ui.KeyOK, 0.3)
	s.shot("loading")
	s.wait(0.9)
	s.shot("saved")
	s.wait(1.5)
	s.tap(ui.KeyDown, 0.4)
	s.tap(ui.KeyOK, 1.2)
	s.shot("failed")
	s.wait(1.5)
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyOK, 0.4)
	s.shot("checked")
	s.taps(ui.KeyDown, 3, 0.3)
	s.tap(ui.KeyOK, 0.5)
	s.shot("radio")
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyOK, 0.6)
	s.shot("sheet")
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.7)
	s.shot("selected")
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.taps(ui.KeyDown, 6, 0.25)
	s.wait(0.3)
	s.shot("long sheet")
	s.tap(ui.KeyBack, 0.6)
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.taps(ui.KeyRight, 2, 0.3)
	s.tap(ui.KeyOK, 0.5)
	s.shot("choice")
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.2)
	s.taps(ui.KeyRight, 2, 0.25)
	s.tap(ui.KeyOK, 0.5)
	s.shot("toggles")
	s.taps(ui.KeyRight, 4, 0.2)
	s.wait(0.3)
	s.shot("chips scrolled")
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyRight, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.shot("flex center")
	s.taps(ui.KeyRight, 2, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.shot("flex between")
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.shot("flex grow")
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyOK, 0.2)
	s.tap(ui.KeyRight, 0.2)
	s.tap(ui.KeyOK, 0.2)
	s.tap(ui.KeyDown, 0.2)
	s.tap(ui.KeyOK, 0.5)
	s.shot("grid")
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.shot("shelves")
	s.taps(ui.KeyRight, 3, 0.25)
	s.wait(0.3)
	s.shot("shelf scrolled")
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.5)
	s.shot("playing")
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyOK, 0.6)
	s.shot("keyboard")
	s.tap(ui.KeyOK, 0.25)
	s.tap(ui.KeyDown, 0.25)
	s.tap(ui.KeyOK, 0.25)
	s.taps(ui.KeyRight, 2, 0.25)
	s.tap(ui.KeyOK, 0.4)
	s.shot("typing")
	s.taps(ui.KeyDown, 2, 0.3)
	s.taps(ui.KeyRight, 2, 0.3)
	s.tap(ui.KeyOK, 0.7)
	s.shot("typed")
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.shot("todo")
	s.tap(ui.KeyOK, 0.6)
	s.tap(ui.KeyOK, 0.25)
	s.tap(ui.KeyDown, 0.25)
	s.tap(ui.KeyOK, 0.25)
	s.taps(ui.KeyRight, 2, 0.25)
	s.tap(ui.KeyOK, 0.25)
	s.taps(ui.KeyDown, 2, 0.3)
	s.taps(ui.KeyRight, 2, 0.3)
	s.tap(ui.KeyOK, 0.7)
	s.shot("todo added")
	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 0.5)
	s.shot("todo checked")
	s.tap(ui.KeyRight, 0.3)
	s.tap(ui.KeyOK, 0.6)
	s.shot("todo sheet")
	s.tap(ui.KeyDown, 0.3)
	s.emit(ui.KeyOK, ui.Press)
	s.wait(1)
	s.emit(ui.KeyOK, ui.Release)
	s.wait(0.8)
	s.shot("todo deleted")
	s.tap(ui.KeyBack, 0.6)

	s.tap(ui.KeyDown, 0.3)
	s.tap(ui.KeyOK, 2.5)
	s.shot("marquee")
	s.taps(ui.KeyDown, 12, 0.15)
	s.wait(0.4)
	s.shot("scrolled")
	s.taps(ui.KeyDown, 25, 0.12)
	s.wait(0.5)
	s.emit(ui.KeyDown, ui.Press)
	s.wait(0.12)
	s.shot("edge")
	s.emit(ui.KeyDown, ui.Release)
	s.wait(1)
	s.tap(ui.KeyBack, 0.8)
	return s
}

func main() {
	out := flag.String("o", "preview", "output directory")
	scale := flag.Int("scale", 4, "GIF pixel scale")
	flag.Parse()
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fail(err)
	}

	scr := &screen{}
	app := ui.NewApp(scr)
	gallery.Build()
	sc := timeline()

	var anim gif.GIF
	var shots []*image.Paletted
	var names []string
	var last []bool
	dt := float32(1) / fps
	start := time.Now()
	next := 0
	for frame := 0; ; frame++ {
		t := float32(frame) * dt
		if next == len(sc.steps) && t > sc.t+0.5 {
			break
		}
		for next < len(sc.steps) && sc.steps[next].at <= t {
			st := sc.steps[next]
			if st.shot != "" {
				// The buffer still shows the last frame, before any key
				// emitted at the same moment takes effect.
				shots = append(shots, scr.image(3))
				names = append(names, st.shot)
			} else {
				app.Input.Emit(st.key, st.phase)
			}
			next++
		}
		if _, err := app.Frame(dt); err != nil {
			fail(err)
		}
		if frame%gifEvery == 0 {
			delay := 100 * gifEvery / fps
			if last != nil && equal(last, scr.px[:]) {
				anim.Delay[len(anim.Delay)-1] += delay
			} else {
				anim.Image = append(anim.Image, scr.image(*scale))
				anim.Delay = append(anim.Delay, delay)
				last = append(last[:0], scr.px[:]...)
			}
		}
		// Pace to the wall clock; see the package comment.
		if ahead := time.Duration(t*float32(time.Second)) - time.Since(start); ahead > 0 {
			time.Sleep(ahead)
		}
	}

	write(filepath.Join(*out, "session.gif"), func(f *os.File) error { return gif.EncodeAll(f, &anim) })
	sheet := contactSheet(shots)
	write(filepath.Join(*out, "shots.png"), func(f *os.File) error { return png.Encode(f, sheet) })
	for i, n := range names {
		fmt.Printf("%2d %s\n", i+1, n)
	}
	fmt.Printf("%d gif frames, %d shots\n", len(anim.Image), len(shots))
}

func equal(a, b []bool) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// contactSheet lays shots out four to a row, left to right.
func contactSheet(shots []*image.Paletted) *image.RGBA {
	const cols, gap = 4, 6
	if len(shots) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	w, h := shots[0].Bounds().Dx(), shots[0].Bounds().Dy()
	rows := (len(shots) + cols - 1) / cols
	sheet := image.NewRGBA(image.Rect(0, 0, cols*(w+gap)+gap, rows*(h+gap)+gap))
	for i := range sheet.Pix {
		sheet.Pix[i] = 60
	}
	for i, s := range shots {
		x, y := gap+(i%cols)*(w+gap), gap+(i/cols)*(h+gap)
		for py := 0; py < h; py++ {
			for px := 0; px < w; px++ {
				sheet.Set(x+px, y+py, s.At(px, py))
			}
		}
	}
	return sheet
}

func write(path string, enc func(*os.File) error) {
	f, err := os.Create(path)
	if err != nil {
		fail(err)
	}
	if err := enc(f); err != nil {
		fail(err)
	}
	if err := f.Close(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "preview:", err)
	os.Exit(1)
}
