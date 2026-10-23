package ui

// KeyPoller turns raw button levels into debounced Press, Repeat and Release
// events with an accelerating repeat. Call Poll for each button from the
// loop that reads the pins; polling avoids interrupt bounce entirely.
type KeyPoller struct {
	// Debounce is how many identical consecutive samples confirm a change.
	Debounce       uint8
	RepeatDelay    float32
	RepeatInterval float32
	MinInterval    float32
	// Accel multiplies the interval after every repeat.
	Accel float32

	keys [keyCount]pollState
}

type pollState struct {
	raw, down      bool
	stable         uint8
	held, next, iv float32
}

func NewKeyPoller() *KeyPoller {
	return &KeyPoller{Debounce: 2, RepeatDelay: 0.2, RepeatInterval: 0.07, MinInterval: 0.03, Accel: 0.9}
}

// Poll feeds one sample for k; dt is the time since the previous sample.
func (p *KeyPoller) Poll(k Key, down bool, dt float32) {
	s := &p.keys[k]
	if down != s.raw {
		s.raw, s.stable = down, 0
	} else if s.stable < p.Debounce {
		s.stable++
	}
	if s.stable >= p.Debounce && s.raw != s.down {
		s.down = s.raw
		if s.down {
			s.held, s.next, s.iv = 0, p.RepeatDelay, p.RepeatInterval
			input.Emit(k, Press)
		} else {
			input.Emit(k, Release)
		}
		return
	}
	if !s.down {
		return
	}
	s.held += dt
	if s.held >= s.next {
		input.Emit(k, Repeat)
		s.next += s.iv
		s.iv = max(s.iv*p.Accel, p.MinInterval)
	}
}
