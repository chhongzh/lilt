package anim

// active holds only moving values. A value joins when its target changes and
// leaves once settled, so idle values cost nothing and nothing needs releasing.
var active = make([]*Value, 0, 64)

// Step advances every moving value by dt seconds and reports whether any
// value moved, which callers use to decide whether to redraw.
func Step(dt float32) bool {
	if len(active) == 0 {
		return false
	}
	dt *= TimeScale
	for i := 0; i < len(active); {
		v := active[i]
		v.step(dt)
		if v.active {
			i++
			continue
		}
		last := len(active) - 1
		active[i] = active[last]
		active[last] = nil
		active = active[:last]
	}
	return true
}

// Busy reports whether any value is still moving.
func Busy() bool { return len(active) > 0 }
