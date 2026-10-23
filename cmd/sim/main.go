// Command sim runs the lilt gallery in a window on the desktop, or in a
// browser when built for wasm. See package desktop for the keys.
package main

import (
	"fmt"
	"os"

	"github.com/chhongzh/lilt/driver/desktop"
	"github.com/chhongzh/lilt/examples/gallery"
)

func main() {
	if err := desktop.Run(desktop.Config{Title: "lilt"}, gallery.Home()); err != nil {
		fmt.Fprintln(os.Stderr, "sim:", err)
		os.Exit(1)
	}
}
