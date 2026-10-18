// Package anim drives every moving value in the framework with one shared
// spring model, so motion feels consistent across widgets.
package anim

// Spring describes motion the way SwiftUI does: Response is the period of
// the undamped oscillation in seconds, Damping the fraction of critical
// damping (1 settles without overshoot, lower values bounce).
type Spring struct {
	Response float32
	Damping  float32
}

const twoPi = 6.2831855

func (s *Spring) coefficients() (k, c float32) {
	w := twoPi / s.Response
	return w * w, 2 * s.Damping * w
}

// Presets are pointers so retuning one at runtime changes every value that
// uses it, which is how the global physics config is meant to be adjusted.
var (
	Default     = &Spring{Response: 0.55, Damping: 0.825}
	Smooth      = &Spring{Response: 0.5, Damping: 1}
	Snappy      = &Spring{Response: 0.5, Damping: 0.85}
	Bouncy      = &Spring{Response: 0.5, Damping: 0.7}
	Interactive = &Spring{Response: 0.15, Damping: 0.86}
)

// TimeScale slows or speeds up all animation; useful when inspecting motion.
var TimeScale float32 = 1

// RubberBandCoefficient matches UIScrollView's overscroll resistance.
var RubberBandCoefficient float32 = 0.55

// RubberBand maps a raw overscroll distance to the displayed one. The result
// approaches but never reaches dim, so pushing harder yields less movement.
func RubberBand(x, dim float32) float32 {
	if dim <= 0 {
		return 0
	}
	neg := x < 0
	if neg {
		x = -x
	}
	y := (1 - 1/(x*RubberBandCoefficient/dim+1)) * dim
	if neg {
		return -y
	}
	return y
}
