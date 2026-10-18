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

func (v *Value) step(dt float32) {
	s := v.spring
	if s == nil {
		s = Default
	}
	k, c := s.coefficients()
	n := int(dt/substep) + 1
	h := dt / float32(n)
	x, vel := v.current-v.target, v.velocity
	for i := 0; i < n; i++ {
		vel += (-k*x - c*vel) * h
		x += vel * h
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

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
