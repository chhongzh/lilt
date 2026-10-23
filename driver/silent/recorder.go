package silent

import (
	"image"
	"image/gif"
	"io"
)

// Recorder collects frames of a Screen into an animated GIF. Frames are drawn
// over the previous one, so each stores only the pixels that changed, and a
// frame identical to the last one just extends its delay.
type Recorder struct {
	// Scale is the size of one screen pixel in the GIF; 0 means 1.
	Scale int

	anim gif.GIF
	last []bool
}

// Frames is how many GIF frames the recording holds so far.
func (r *Recorder) Frames() int { return len(r.anim.Image) }

// Capture adds the current state of s, shown for delay hundredths of a second.
func (r *Recorder) Capture(s *Screen, delay int) {
	scale := max(r.Scale, 1)
	var box image.Rectangle
	if r.last == nil {
		box = image.Rect(0, 0, s.w, s.h)
	} else {
		box = s.Diff(r.last)
	}
	if box.Empty() {
		r.anim.Delay[len(r.anim.Delay)-1] += delay
		return
	}
	r.anim.Image = append(r.anim.Image, s.Region(box, scale))
	r.anim.Delay = append(r.anim.Delay, delay)
	r.anim.Disposal = append(r.anim.Disposal, gif.DisposalNone)
	r.last = append(r.last[:0], s.px...)
}

// Encode writes the recording to w.
func (r *Recorder) Encode(w io.Writer) error { return gif.EncodeAll(w, &r.anim) }
