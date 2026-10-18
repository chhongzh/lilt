package anim

// substep keeps semi-implicit Euler stable for stiff springs at low frame rates.
const substep = 1.0 / 240

// Value is a float that springs toward its target. The zero value is ready to
// use, rests at 0 and follows Default. Values must not be copied once animating
// because the scheduler holds a pointer to them.
type Value struct {
	current   float32
	target    float32
	velocity  float32
	precision float32
	spring    *Spring
	parent    *Value
	active    bool
}

// NewValue returns a value resting at v.
func NewValue(v float32, s *Spring) *Value {
	return &Value{current: v, target: v, spring: s}
}

// Get returns the displayed value, including the parent's.
func (v *Value) Get() float32 {
	if v.parent != nil {
		return v.current + v.parent.Get()
	}
	return v.current
}

// Current returns the value without the parent's contribution.
func (v *Value) Current() float32 { return v.current }

// Target returns where the value is heading, without the parent's contribution.
func (v *Value) Target() float32 { return v.target }

// Velocity is in units per second.
func (v *Value) Velocity() float32 { return v.velocity }

// Settled reports whether the value is at rest.
func (v *Value) Settled() bool { return !v.active }

// Set animates toward x.
func (v *Value) Set(x float32) {
	if x == v.target {
		return
	}
	v.target = x
	v.wake()
}

// Jump moves to x immediately and stops any motion.
func (v *Value) Jump(x float32) {
	v.current, v.target, v.velocity = x, x, 0
}

// Shift moves the current value by d while keeping the target, so the value
// glides back. Layout uses it to animate position changes.
func (v *Value) Shift(d float32) {
	if d == 0 {
		return
	}
	v.current += d
	v.wake()
}

// Kick adds velocity, for feedback such as a nudge at a list edge.
func (v *Value) Kick(velocity float32) {
	v.velocity += velocity
	v.wake()
}

// Bump kicks the value hard enough that, critically damped and at rest, it
// swings about px away before coming back. Bouncier springs swing somewhat
// further. Eased values ignore it.
func (v *Value) Bump(px float32) {
	v.Kick(px * v.curve().rate() * e)
}

// SetSpring changes the physics; nil means Default.
func (v *Value) SetSpring(s *Spring) { v.spring = s }

// SetParent makes Get add p's value, so a child follows its parent's motion
// on top of its own.
func (v *Value) SetParent(p *Value) { v.parent = p }

// SetPrecision sets the distance under which the value counts as settled.
// The default of 0.05 suits pixel values; use smaller ones for 0..1 progress.
func (v *Value) SetPrecision(p float32) { v.precision = p }

func (v *Value) wake() {
	if !v.active {
		v.active = true
		active = append(active, v)
	}
}

func (v *Value) curve() *Spring {
	if v.spring == nil {
		return Default
	}
	return v.spring
}

func (v *Value) step(dt float32) {
	s := v.curve()
	x, vel := v.current-v.target, v.velocity
	if s.Ease {
		x, vel = x*expNeg(s.rate()*dt), 0
	} else {
		k, c := s.coefficients()
		n := int(dt/substep) + 1
		h := dt / float32(n)
		for i := 0; i < n; i++ {
			vel += (-k*x - c*vel) * h
			x += vel * h
		}
	}
	p := v.precision
	if p == 0 {
		p = 0.05
	}
	if abs(x) < p && abs(vel) < p*10 {
		v.current, v.velocity, v.active = v.target, 0, false
		return
	}
	v.current, v.velocity = v.target+x, vel
}

const e = 2.7182817

// expNeg approximates e^(-x) for x >= 0. math.Exp is float64, which the
// ESP32-S3 FPU does in software, so the argument is halved until small, a
// (2,2) Padé approximant applied, and the result squared back up.
func expNeg(x float32) float32 {
	if x <= 0 {
		return 1
	}
	if x > 16 {
		return 0
	}
	n := 0
	for x > 0.5 {
		x *= 0.5
		n++
	}
	x2 := x * x / 12
	h := x / 2
	r := (1 - h + x2) / (1 + h + x2)
	for ; n > 0; n-- {
		r *= r
	}
	return r
}

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
