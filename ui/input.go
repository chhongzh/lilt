package ui

// Key is a logical key. Mapping pins or scan codes to keys is up to the board code.
type Key uint8

const (
	KeyNone Key = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyOK
	KeyBack
	KeyMenu
	KeyBackspace
	// KeyChar carries a typed character in KeyEvent.Rune.
	KeyChar
	keyCount
)

func (k Key) Arrow() bool { return k >= KeyUp && k <= KeyRight }

type Phase uint8

const (
	Press Phase = iota
	Repeat
	Release
)

type KeyEvent struct {
	Key   Key
	Phase Phase
	Rune  rune
	// Canceled marks a Release that must not trigger an action, sent when a
	// touch turns into a drag.
	Canceled bool
	// Repeated marks a Release after Repeats, whose action already happened.
	Repeated bool
}

// Fires reports whether e should perform its key's action: a Release that
// was neither canceled nor preceded by Repeats, or each Repeat while held.
// Press itself only gives feedback.
func (e *KeyEvent) Fires() bool {
	return e.Phase == Repeat || e.Phase == Release && !e.Canceled && !e.Repeated
}

// KeyHandler receives keys bubbling up from the focused node. Returning true
// on Press makes the node the owner of that key's Repeat and Release.
type KeyHandler interface{ HandleKey(e *KeyEvent) bool }

// KeyPreviewer sees keys on the way down from the root, before the focused
// node, so a container can intercept them.
type KeyPreviewer interface{ PreviewKey(e *KeyEvent) bool }

type PointerPhase uint8

const (
	PointerDown PointerPhase = iota
	PointerMove
	PointerUp
)

type PointerEvent struct {
	X, Y  int16
	Phase PointerPhase
}

// PointerHandler is offered pointer events from the node under the pointer
// up to the root. Returning true on a move captures the pointer, which is
// how a scroll view takes over a drag.
type PointerHandler interface{ HandlePointer(e *PointerEvent) bool }

type event struct {
	key     KeyEvent
	ptr     PointerEvent
	pointer bool
}

// dragSlop is how far a touch may wander before it counts as a drag.
const dragSlop = 4

// inputManager queues events from any goroutine and dispatches them on the
// UI goroutine at the start of each frame.
type inputManager struct {
	queue chan event
	owner [keyCount]Node
	// nav marks keys the framework took on Press: arrows for focus
	// movement and Back for the screen stack.
	nav      [keyCount]bool
	repeated [keyCount]bool
	path     []Node

	ptrDown    bool
	ptrStart   PointerEvent
	ptrHit     Node
	ptrPressed Node
	ptrCapture PointerHandler
}

func (m *inputManager) init() {
	m.queue = make(chan event, 32)
	m.path = make([]Node, 0, 32)
}

// reset forgets held keys, the pointer and anything still queued.
func (m *inputManager) reset() {
	for {
		select {
		case <-m.queue:
			continue
		default:
		}
		break
	}
	m.owner, m.nav, m.repeated = [keyCount]Node{}, [keyCount]bool{}, [keyCount]bool{}
	m.ptrDown, m.ptrHit, m.ptrPressed, m.ptrCapture = false, nil, nil, nil
}

// Emit queues a key event. Safe to call from any goroutine.
func (m *inputManager) Emit(k Key, p Phase) { m.push(event{key: KeyEvent{Key: k, Phase: p}}) }

// EmitRune queues a typed character.
func (m *inputManager) EmitRune(r rune) {
	m.push(event{key: KeyEvent{Key: KeyChar, Phase: Press, Rune: r}})
}

// EmitPointer queues a touch or mouse event in screen pixels.
func (m *inputManager) EmitPointer(x, y int16, p PointerPhase) {
	m.push(event{ptr: PointerEvent{x, y, p}, pointer: true})
}

// push drops the event when the queue is full: a stalled UI must not block
// the goroutine polling the buttons.
func (m *inputManager) push(e event) {
	select {
	case m.queue <- e:
	default:
	}
}

func (m *inputManager) drain() bool {
	any := false
	for {
		select {
		case e := <-m.queue:
			any = true
			if e.pointer {
				m.pointer(e.ptr)
			} else {
				m.dispatch(e.key)
			}
		default:
			return any
		}
	}
}

func (m *inputManager) dispatch(e KeyEvent) {
	if e.Key == KeyOK {
		cursor.pressed(e.Phase)
	}
	if e.Key.Arrow() {
		cursor.hasPointer = false
	}
	switch e.Phase {
	case Press:
		m.repeated[e.Key] = false
	case Repeat:
		m.repeated[e.Key] = true
	case Release:
		e.Repeated = m.repeated[e.Key]
	}
	if e.Phase != Press {
		if o := m.owner[e.Key]; o != nil {
			if e.Phase == Release {
				m.owner[e.Key] = nil
			}
			deliver(o, &e)
			return
		}
		if m.nav[e.Key] {
			if e.Phase == Release {
				m.nav[e.Key] = false
			}
			m.unclaimed(&e)
			return
		}
		if e.Phase == Release {
			return
		}
	}
	m.path = m.path[:0]
	for n := focus.current; n != nil; n = n.box().parent {
		m.path = append(m.path, n)
	}
	for i := len(m.path) - 1; i >= 0; i-- {
		if p, ok := m.path[i].(KeyPreviewer); ok && p.PreviewKey(&e) {
			m.claim(e, m.path[i])
			return
		}
	}
	for _, n := range m.path {
		if h, ok := n.(KeyHandler); ok && h.HandleKey(&e) {
			m.claim(e, n)
			return
		}
	}
	if e.Phase == Press && (e.Key.Arrow() || e.Key == KeyBack) {
		m.nav[e.Key] = true
		m.unclaimed(&e)
	}
}

// unclaimed runs the framework's own behaviour for keys no node took.
func (m *inputManager) unclaimed(e *KeyEvent) {
	if e.Key == KeyBack {
		layers.back(e)
	} else {
		navigate(e)
	}
}

func (m *inputManager) claim(e KeyEvent, n Node) {
	if e.Phase == Press {
		m.owner[e.Key] = n
	}
}

func deliver(n Node, e *KeyEvent) {
	if p, ok := n.(KeyPreviewer); ok && p.PreviewKey(e) {
		return
	}
	if h, ok := n.(KeyHandler); ok {
		h.HandleKey(e)
	}
}

// forget drops every reference into a subtree that is leaving the tree.
func (m *inputManager) forget(d Node) {
	for k, o := range m.owner {
		if o != nil && IsAncestor(d, o) {
			m.owner[k] = nil
		}
	}
	if m.ptrHit != nil && IsAncestor(d, m.ptrHit) {
		m.ptrHit, m.ptrPressed, m.ptrCapture = nil, nil, nil
	}
}

func (m *inputManager) pointer(e PointerEvent) {
	switch e.Phase {
	case PointerDown:
		m.ptrDown, m.ptrStart, m.ptrCapture = true, e, nil
		m.ptrHit = hitTest(e.X, e.Y)
		cursor.pointerAt(e.X, e.Y)
		m.offer(&e)
		if t := focusableAncestor(m.ptrHit); t != nil {
			focus.Set(t)
			m.ptrPressed = t
			m.dispatch(KeyEvent{Key: KeyOK, Phase: Press})
		}
	case PointerMove:
		cursor.pointerAt(e.X, e.Y)
		if m.ptrCapture != nil {
			m.ptrCapture.HandlePointer(&e)
			return
		}
		if !m.ptrDown {
			if t := focusableAncestor(hitTest(e.X, e.Y)); t != nil {
				focus.Set(t)
			}
			return
		}
		if absf(float32(e.X-m.ptrStart.X))+absf(float32(e.Y-m.ptrStart.Y)) < dragSlop {
			return
		}
		if h := m.offer(&e); h != nil {
			m.ptrCapture = h
			if m.ptrPressed != nil {
				m.ptrPressed = nil
				m.dispatch(KeyEvent{Key: KeyOK, Phase: Release, Canceled: true})
			}
		}
	case PointerUp:
		if m.ptrCapture != nil {
			m.ptrCapture.HandlePointer(&e)
		} else {
			m.offer(&e)
		}
		if m.ptrPressed != nil {
			m.dispatch(KeyEvent{Key: KeyOK, Phase: Release})
		}
		m.ptrDown, m.ptrHit, m.ptrPressed, m.ptrCapture = false, nil, nil, nil
	}
}

func (m *inputManager) offer(e *PointerEvent) PointerHandler {
	for n := m.ptrHit; n != nil; n = n.box().parent {
		if h, ok := n.(PointerHandler); ok && h.HandlePointer(e) {
			return h
		}
	}
	return nil
}

func focusableAncestor(d Node) Node {
	for n := d; n != nil; n = n.box().parent {
		if n.box().Focusable() {
			return n
		}
	}
	return nil
}
