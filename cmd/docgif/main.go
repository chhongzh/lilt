// Command docgif records the animated demos of the documentation site. Each
// demo builds a small UI, plays a key script against it on a silent screen
// and becomes docs/assets/demos/<name>.gif.
//
// Demos live in the demo_*.go files, one file per documentation page, and
// register themselves in init. The code between the snippet markers
//
//	// --8<-- [start:name]
//	// --8<-- [end:name]
//
// is included verbatim in the Markdown pages, so the examples on the site
// are the code that produced the GIFs next to them.
//
//	go run ./cmd/docgif                   record every demo
//	go run ./cmd/docgif button-*          record the demos matching a pattern
//	go run ./cmd/docgif -list             print the demo names
//	go run ./cmd/docgif -shots /tmp/shots also write a PNG of each demo's stills
package main

import (
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/driver/silent"
	"github.com/chhongzh/lilt/internal/reel"
	"github.com/chhongzh/lilt/ui"
)

// demo is one recording.
type demo struct {
	// name is the GIF's file name without extension, lower-case words
	// joined by hyphens and starting with the page's name.
	name string
	// w and h are the screen size; zero means the 128x64 of the SH1106.
	w, h int
	// build returns the root node. It runs once, right before recording.
	build func() ui.Node
	// script plays the keys. Leave time at the end for every handler to
	// return: one still sleeping when the recording stops would act on
	// the next demo.
	script func(s *reel.Script)
}

var demos []demo

func register(d demo) {
	if d.name == "" || d.build == nil || d.script == nil {
		panic("docgif: demo needs a name, build and script")
	}
	for _, o := range demos {
		if o.name == d.name {
			panic("docgif: demo " + d.name + " registered twice")
		}
	}
	demos = append(demos, d)
}

func main() {
	out := flag.String("o", "docs/assets/demos", "output directory for the GIFs")
	scale := flag.Int("scale", 3, "size of one screen pixel in the GIF")
	shots := flag.String("shots", "", "if set, write <name>.png of each demo's stills here for review")
	list := flag.Bool("list", false, "print the demo names and exit")
	flag.Parse()

	slices.SortFunc(demos, func(a, b demo) int { return strings.Compare(a.name, b.name) })
	if *list {
		for _, d := range demos {
			fmt.Println(d.name)
		}
		return
	}
	picked := pick(flag.Args())
	if len(picked) == 0 {
		fail(fmt.Errorf("no demo matches %q; see -list", flag.Args()))
	}
	for _, dir := range []string{*out, *shots} {
		if dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				fail(err)
			}
		}
	}

	start := time.Now()
	for _, d := range picked {
		res, err := record(d, *scale)
		if err != nil {
			fail(fmt.Errorf("%s: %w", d.name, err))
		}
		write(filepath.Join(*out, d.name+".gif"), res.GIF.Encode)
		if *shots != "" && len(res.Shots) > 0 {
			sheet := reel.ContactSheet(res.Shots, 4)
			write(filepath.Join(*shots, d.name+".png"), func(w io.Writer) error { return png.Encode(w, sheet) })
		}
		fmt.Printf("%-28s %4d frames  %s\n", d.name, res.GIF.Frames(), strings.Join(res.Names, ", "))
	}
	fmt.Printf("%d demos in %s\n", len(picked), time.Since(start).Round(time.Second))
}

// pick returns the demos whose names match any of the patterns, or all of
// them when there are none.
func pick(patterns []string) []demo {
	if len(patterns) == 0 {
		return demos
	}
	var out []demo
	for _, d := range demos {
		for _, p := range patterns {
			if ok, err := path.Match(p, d.name); err != nil {
				fail(err)
			} else if ok {
				out = append(out, d)
				break
			}
		}
	}
	return out
}

// record runs one demo. The theme, the spring presets, the cursor's
// included, and the text input method are global, so a demo may change
// them; they are put back afterwards.
func record(d demo, scale int) (reel.Result, error) {
	theme, input := *ui.CurrentTheme, ui.TextInput
	springs := []*anim.Spring{
		anim.Default, anim.Smooth, anim.Snappy, anim.Bouncy, anim.Interactive,
		ui.CursorLead, ui.CursorTrail, ui.CursorRelease, ui.CursorLeap, ui.CursorLean,
	}
	saved := make([]anim.Spring, len(springs))
	for i, s := range springs {
		saved[i] = *s
	}
	defer func() {
		*ui.CurrentTheme, ui.TextInput = theme, input
		for i, s := range springs {
			*s = saved[i]
		}
	}()

	w, h := d.w, d.h
	if w == 0 || h == 0 {
		w, h = 128, 64
	}
	s := &reel.Script{}
	d.script(s)
	return reel.Record(silent.New(w, h), d.build(), s, scale)
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
	fmt.Fprintln(os.Stderr, "docgif:", err)
	os.Exit(1)
}
