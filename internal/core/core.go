// Package core holds the hooks package ui installs for driving the frame
// loop. Package lilt runs the loop with them and cmd/preview steps frames
// by hand; neither belongs in the public API.
package core

import "github.com/chhongzh/lilt/gfx"

var (
	// Boot binds the screen and pushes root, if any, as the first screen.
	Boot func(s gfx.Screen, root any)
	// Frame advances the UI by dt seconds and reports whether it flushed.
	Frame func(dt float32) (bool, error)
	// Reset drops every screen, the focus and any held keys.
	Reset func()
)
