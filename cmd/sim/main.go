// Command sim runs the lilt gallery in a window on the desktop, or in a
// browser when built for wasm. See package desktop for the keys.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chhongzh/lilt/driver/desktop"
	"github.com/chhongzh/lilt/examples/gallery"
	"github.com/chhongzh/lilt/ui"
)

func main() {
	w := flag.Int("w", 128, "screen width in physical pixels")
	h := flag.Int("h", 64, "screen height in physical pixels")
	zoom := flag.Float64("zoom", 1, "multiplies the scale the screen size gives")
	flag.Parse()

	// The fonts depend on the scale, so the zoom has to be known first.
	ui.SetZoom(float32(*zoom))
	gallery.FitFonts(int16(*w), int16(*h))
	cfg := desktop.Config{Title: "lilt", Width: *w, Height: *h}
	if err := desktop.Run(cfg, gallery.Home()); err != nil {
		fmt.Fprintln(os.Stderr, "sim:", err)
		os.Exit(1)
	}
}
