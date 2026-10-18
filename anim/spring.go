// Package anim drives every moving value in the framework with one shared
// spring model, so motion feels consistent across widgets.
package anim

// Spring describes motion the way SwiftUI does: Response is the period of
// the undamped oscillation in seconds, Damping the fraction of critical
// damping (1 settles without overshoot, lower values bounce).
type Spring struct {
	Response float32
	Damping  float32
	// Ease drops the spring for exponential easing: full speed from the
	// first frame, no momentum and no overshoot. The value closes 1-1/e of
	// the gap every Response/2π seconds and Damping is ignored.
	Ease bool
}

const twoPi = 6.2831855

func (s *Spring) coefficients() (k, c float32) {
	w := s.rate()
	return w * w, 2 * s.Damping * w
}

// rate is the angular frequency, or the easing rate, in 1/s.
func (s *Spring) rate() float32 { return twoPi / s.Response }

// Presets are pointers so retuning one at runtime changes every value that
// uses it, which is how the global physics config is meant to be adjusted.
// Each overshoots a little: the bounce is the point.
var (
	// Default moves list items and layout, with a light overshoot.
	Default = &Spring{Response: 0.29, Damping: 0.75}
	// Smooth carries whole screens in and out.
	Smooth = &Spring{Response: 0.39, Damping: 0.85}
	// Snappy is for thumbs and progress that should track closely.
	Snappy = &Spring{Response: 0.31, Damping: 0.8}
	// Bouncy pops marks in so they land with a snap.
	Bouncy = &Spring{Response: 0.25, Damping: 0.45}
	// Interactive follows a finger or a key going down: quick, no bounce.
	Interactive = &Spring{Response: 0.14, Damping: 1}
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
