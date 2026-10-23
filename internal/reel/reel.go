// Package reel replays a scripted key timeline against a UI on a silent
// screen and records it as an animated GIF plus stills of marked moments.
// cmd/preview records the gallery with it and cmd/docgif the documentation
// demos.
//
// Frames are stepped by hand through internal/core instead of lilt.Start,
// but paced to the wall clock, so event callbacks, which run on their own
// goroutines and sleep in real time, finish in step with the frames.
package reel

import (
	"image"
	"time"

	"github.com/chhongzh/lilt/driver/silent"
	"github.com/chhongzh/lilt/internal/core"
	"github.com/chhongzh/lilt/ui"
)

const (
	fps = 100
	// GIF frames are sampled at half the frame rate, 2/100 s apart: browsers
	// slow down any shorter GIF delay.
	gifEvery = 2
	// tail keeps recording after the last step so motion can settle.
	tail = 0.5
)

type step struct {
	at    float32
	key   ui.Key
	phase ui.Phase
	shot  string
}

// Script is a timeline of key events and named stills, built by calls in
// the order they happen. Times are in seconds.
type Script struct {
	t     float32
	steps []step
}

// Wait moves the timeline on by d.
func (s *Script) Wait(d float32) { s.t += d }

// Emit sends one key event at the current time.
func (s *Script) Emit(k ui.Key, p ui.Phase) {
	s.steps = append(s.steps, step{at: s.t, key: k, phase: p})
}

// Shot marks a still named name. It shows the frame before any key emitted
// at the same moment takes effect.
func (s *Script) Shot(name string) { s.steps = append(s.steps, step{at: s.t, shot: name}) }

// Tap presses and releases k, then waits gap.
func (s *Script) Tap(k ui.Key, gap float32) {
	s.Emit(k, ui.Press)
	s.Wait(0.08)
	s.Emit(k, ui.Release)
	s.Wait(gap)
}

// Taps taps k n times, gap apart.
func (s *Script) Taps(k ui.Key, n int, gap float32) {
	for range n {
		s.Tap(k, gap)
	}
}

// Hold presses k, keeps it down for d and releases it, without repeats.
func (s *Script) Hold(k ui.Key, d float32) {
	s.Emit(k, ui.Press)
	s.Wait(d)
	s.Emit(k, ui.Release)
}

// Duration is the time of the last step.
func (s *Script) Duration() float32 { return s.t }

// Result is what Record captured.
type Result struct {
	GIF *silent.Recorder
	// Shots are the stills, at the scale passed to Record, named by Names.
	Shots []*image.Paletted
	Names []string
}

// Record shows root on scr, plays s against it and records every other
// frame into a GIF with the given pixel scale. It resets the UI afterwards,
// so several recordings can run in one process, one at a time.
func Record(scr *silent.Screen, root ui.Node, s *Script, scale int) (Result, error) {
	core.Boot(scr, root)
	defer core.Reset()
	res := Result{GIF: &silent.Recorder{Scale: scale}}
	dt := float32(1) / fps
	start := time.Now()
	next := 0
	for frame := 0; ; frame++ {
		t := float32(frame) * dt
		if next == len(s.steps) && t > s.t+tail {
			break
		}
		for next < len(s.steps) && s.steps[next].at <= t {
			st := s.steps[next]
			if st.shot != "" {
				res.Shots = append(res.Shots, scr.Image(scale))
				res.Names = append(res.Names, st.shot)
			} else {
				ui.Emit(st.key, st.phase)
			}
			next++
		}
		if _, err := core.Frame(dt); err != nil {
			return res, err
		}
		if frame%gifEvery == 0 {
			res.GIF.Capture(scr, 100*gifEvery/fps)
		}
		if ahead := time.Duration(t*float32(time.Second)) - time.Since(start); ahead > 0 {
			time.Sleep(ahead)
		}
	}
	return res, nil
}

// ContactSheet lays shots out cols to a row, left to right, on a grey
// background.
func ContactSheet(shots []*image.Paletted, cols int) *image.RGBA {
	const gap = 6
	if len(shots) == 0 {
		return image.NewRGBA(image.Rect(0, 0, 1, 1))
	}
	cols = min(cols, len(shots))
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
