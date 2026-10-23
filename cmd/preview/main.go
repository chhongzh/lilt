// Command preview runs the gallery on the host against a scripted key
// timeline and writes an animated GIF of the session plus a PNG contact
// sheet of marked moments. The recording itself is package reel.
package main

import (
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/chhongzh/lilt/driver/silent"
	"github.com/chhongzh/lilt/examples/gallery"
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
)

const width, height = 128, 64

func timeline() *reel.Script {
	s := &reel.Script{}
	s.Wait(0.6)
	s.Shot("home")

	s.Tap(ui.KeyOK, 0.6)
	s.Shot("buttons")
	s.Tap(ui.KeyOK, 0.3)
	s.Shot("tapped")
	s.Tap(ui.KeyDown, 0.4)
	s.Emit(ui.KeyOK, ui.Press)
	s.Wait(0.4)
	s.Shot("holding")
	s.Wait(0.5)
	s.Shot("hold fired")
	s.Emit(ui.KeyOK, ui.Release)
	s.Wait(1.6)
	s.Tap(ui.KeyDown, 0.4)
	s.Tap(ui.KeyOK, 0.3)
	s.Shot("loading")
	s.Wait(0.9)
	s.Shot("saved")
	s.Wait(1.5)
	s.Tap(ui.KeyDown, 0.4)
	s.Tap(ui.KeyOK, 1.2)
	s.Shot("failed")
	s.Wait(1.5)
	s.Emit(ui.KeyBack, ui.Press)
	s.Wait(0.8)
	s.Shot("back held")
	s.Emit(ui.KeyBack, ui.Release)
	s.Wait(0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyOK, 0.4)
	s.Shot("checked")
	s.Taps(ui.KeyDown, 2, 0.3)
	s.Tap(ui.KeyOK, 0.5)
	s.Shot("syncing")
	s.Wait(0.9)
	s.Shot("sync failed")
	s.Taps(ui.KeyDown, 2, 0.3)
	s.Tap(ui.KeyOK, 0.5)
	s.Shot("radio")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("sheet")
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.7)
	s.Shot("selected")
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Taps(ui.KeyDown, 6, 0.25)
	s.Wait(0.3)
	s.Shot("long sheet")
	s.Tap(ui.KeyBack, 0.6)
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("drawer bottom")
	s.Tap(ui.KeyBack, 0.5)
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("drawer top")
	s.Tap(ui.KeyBack, 0.5)
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyDown, 0.3)
	s.Shot("drawer left")
	s.Emit(ui.KeyBack, ui.Press)
	s.Wait(0.4)
	s.Shot("drawer peek")
	s.Emit(ui.KeyBack, ui.Release)
	s.Wait(0.5)
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyOK, 0.4)
	s.Shot("menu toast")
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("dialog")
	s.Tap(ui.KeyRight, 0.3)
	s.Tap(ui.KeyOK, 0.4)
	s.Shot("deleting")
	s.Wait(0.8)
	s.Shot("deleted")
	s.Wait(1.5)
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Taps(ui.KeyRight, 2, 0.3)
	s.Tap(ui.KeyOK, 0.5)
	s.Shot("choice")
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.2)
	s.Taps(ui.KeyRight, 2, 0.25)
	s.Tap(ui.KeyOK, 0.5)
	s.Shot("toggles")
	s.Taps(ui.KeyRight, 4, 0.2)
	s.Wait(0.3)
	s.Shot("chips scrolled")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyRight, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("flex center")
	s.Taps(ui.KeyRight, 2, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("flex between")
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("flex grow")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyOK, 0.2)
	s.Tap(ui.KeyRight, 0.2)
	s.Tap(ui.KeyOK, 0.2)
	s.Tap(ui.KeyDown, 0.2)
	s.Tap(ui.KeyOK, 0.5)
	s.Shot("grid")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("shelves")
	s.Taps(ui.KeyRight, 3, 0.25)
	s.Wait(0.3)
	s.Shot("shelf scrolled")
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.5)
	s.Shot("playing")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("nested")
	s.Taps(ui.KeyDown, 7, 0.2)
	s.Wait(0.3)
	s.Shot("inbox scrolled")
	s.Tap(ui.KeyDown, 0.4)
	s.Taps(ui.KeyRight, 4, 0.2)
	s.Wait(0.3)
	s.Shot("shelf")
	s.Tap(ui.KeyDown, 0.3)
	s.Taps(ui.KeyRight, 3, 0.2)
	s.Taps(ui.KeyDown, 3, 0.2)
	s.Tap(ui.KeyOK, 0.4)
	s.Shot("board")
	s.Taps(ui.KeyDown, 3, 0.25)
	s.Wait(0.3)
	s.Shot("past board")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("keyboard")
	s.Tap(ui.KeyOK, 0.25)
	s.Tap(ui.KeyDown, 0.25)
	s.Tap(ui.KeyOK, 0.25)
	s.Taps(ui.KeyRight, 2, 0.25)
	s.Tap(ui.KeyOK, 0.4)
	s.Shot("typing")
	s.Taps(ui.KeyDown, 2, 0.3)
	s.Taps(ui.KeyRight, 2, 0.3)
	s.Tap(ui.KeyOK, 0.7)
	s.Shot("typed")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("todo")
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyOK, 0.25)
	s.Tap(ui.KeyDown, 0.25)
	s.Tap(ui.KeyOK, 0.25)
	s.Taps(ui.KeyRight, 2, 0.25)
	s.Tap(ui.KeyOK, 0.25)
	s.Taps(ui.KeyDown, 2, 0.3)
	s.Taps(ui.KeyRight, 2, 0.3)
	s.Tap(ui.KeyOK, 0.7)
	s.Shot("todo added")
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 0.5)
	s.Shot("todo checked")
	s.Tap(ui.KeyRight, 0.3)
	s.Tap(ui.KeyOK, 0.6)
	s.Shot("todo sheet")
	s.Tap(ui.KeyDown, 0.3)
	s.Emit(ui.KeyOK, ui.Press)
	s.Wait(1)
	s.Emit(ui.KeyOK, ui.Release)
	s.Wait(0.8)
	s.Shot("todo deleted")
	s.Tap(ui.KeyBack, 0.6)

	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 2.5)
	s.Shot("marquee")
	s.Taps(ui.KeyDown, 12, 0.15)
	s.Wait(0.4)
	s.Shot("scrolled")
	s.Taps(ui.KeyDown, 25, 0.12)
	s.Wait(0.5)
	s.Emit(ui.KeyDown, ui.Press)
	s.Wait(0.12)
	s.Shot("edge")
	s.Emit(ui.KeyDown, ui.Release)
	s.Wait(1)
	s.Tap(ui.KeyBack, 0.8)
	return s
}

// highlights is a half-minute reel of the moments that show the feel best,
// for the README.
func highlights() *reel.Script {
	s := &reel.Script{}
	s.Wait(0.6)
	s.Taps(ui.KeyDown, 5, 0.14)
	s.Wait(0.2)
	s.Taps(ui.KeyUp, 5, 0.14)
	s.Wait(0.3)

	s.Tap(ui.KeyOK, 0.5)
	s.Tap(ui.KeyOK, 0.4)
	s.Tap(ui.KeyDown, 0.3)
	s.Emit(ui.KeyOK, ui.Press)
	s.Wait(1)
	s.Emit(ui.KeyOK, ui.Release)
	s.Wait(0.4)
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 1.6)
	s.Tap(ui.KeyDown, 0.3)
	s.Tap(ui.KeyOK, 1.6)
	s.Tap(ui.KeyBack, 0.5)

	s.Taps(ui.KeyDown, 2, 0.2)
	s.Tap(ui.KeyOK, 0.5)
	s.Tap(ui.KeyDown, 0.25)
	s.Tap(ui.KeyOK, 0.5)
	s.Taps(ui.KeyDown, 5, 0.16)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyBack, 0.5)

	s.Tap(ui.KeyDown, 0.2)
	s.Tap(ui.KeyOK, 0.5)
	s.Taps(ui.KeyDown, 2, 0.2)
	s.Tap(ui.KeyOK, 0.7)
	s.Tap(ui.KeyBack, 0.5)
	s.Taps(ui.KeyDown, 2, 0.2)
	s.Tap(ui.KeyOK, 0.5)
	s.Tap(ui.KeyRight, 0.3)
	s.Tap(ui.KeyOK, 2)
	s.Tap(ui.KeyBack, 0.5)

	s.Tap(ui.KeyDown, 0.2)
	s.Tap(ui.KeyOK, 0.5)
	s.Taps(ui.KeyRight, 2, 0.25)
	s.Tap(ui.KeyOK, 0.4)
	s.Tap(ui.KeyDown, 0.25)
	s.Tap(ui.KeyOK, 0.25)
	s.Taps(ui.KeyRight, 5, 0.2)
	s.Tap(ui.KeyOK, 0.5)
	s.Tap(ui.KeyBack, 0.5)

	s.Taps(ui.KeyDown, 3, 0.2)
	s.Tap(ui.KeyOK, 0.5)
	s.Taps(ui.KeyRight, 4, 0.25)
	s.Tap(ui.KeyDown, 0.3)
	s.Taps(ui.KeyRight, 2, 0.25)
	s.Tap(ui.KeyOK, 0.6)
	s.Tap(ui.KeyBack, 0.5)

	s.Taps(ui.KeyDown, 4, 0.2)
	s.Tap(ui.KeyOK, 2)
	// Holding Down runs to the bottom on key repeat and pushes past it.
	s.Emit(ui.KeyDown, ui.Press)
	for range 34 {
		s.Wait(0.07)
		s.Emit(ui.KeyDown, ui.Repeat)
	}
	s.Wait(0.5)
	s.Emit(ui.KeyDown, ui.Release)
	s.Wait(0.7)
	s.Tap(ui.KeyBack, 0.8)
	return s
}

func main() {
	out := flag.String("o", "preview", "output directory")
	scale := flag.Int("scale", 4, "GIF pixel scale")
	highlight := flag.Bool("highlight", false, "record the short highlight reel as highlight.gif")
	flag.Parse()
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fail(err)
	}

	sc, name := timeline(), "session.gif"
	if *highlight {
		sc, name = highlights(), "highlight.gif"
	}
	res, err := reel.Record(silent.New(width, height), gallery.Home(), sc, *scale)
	if err != nil {
		fail(err)
	}

	write(filepath.Join(*out, name), res.GIF.Encode)
	if len(res.Shots) > 0 {
		sheet := reel.ContactSheet(res.Shots, 4)
		write(filepath.Join(*out, "shots.png"), func(w io.Writer) error { return png.Encode(w, sheet) })
	}
	for i, n := range res.Names {
		fmt.Printf("%2d %s\n", i+1, n)
	}
	fmt.Printf("%d gif frames, %d shots\n", res.GIF.Frames(), len(res.Shots))
}

func write(path string, enc func(io.Writer) error) {
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
