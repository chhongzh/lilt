// Package lilt runs the UI. Start binds a screen and spins up the UI
// goroutine, which owns the tree from then on; Pause freezes it and Stop
// ends it. Everything else lives in packages ui and widget.
//
//	lilt.Start(screen, home)
//	select {} // or run board code here
package lilt

import (
	"sync"
	"time"

	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/internal/core"
	"github.com/chhongzh/lilt/ui"
)

// FPS caps the frame rate. It is read every frame.
var FPS = 60

type command uint8

const (
	pause command = iota
	resume
	stop
)

var (
	mu     sync.Mutex
	ctl    chan command
	done   chan error
	paused bool
)

// Start shows root, if not nil, as the first screen of s and starts the UI
// goroutine. Build root before calling Start; afterwards change the tree
// only from inside ui.Do. Called while paused, Start resumes and ignores
// its arguments.
func Start(s gfx.Screen, root ui.Node) {
	mu.Lock()
	defer mu.Unlock()
	if ctl != nil {
		if paused {
			paused = false
			send(resume)
		}
		return
	}
	core.Boot(s, root)
	ctl, done = make(chan command), make(chan error, 1)
	go run(ctl, done)
}

// Pause stops updating and drawing until Start is called again. Keys
// arriving meanwhile wait in the input queue, dropped once it is full.
func Pause() {
	mu.Lock()
	defer mu.Unlock()
	if ctl != nil && !paused {
		paused = true
		send(pause)
	}
}

// Stop ends the UI goroutine and drops every screen, so Start can begin
// afresh. It returns the error that stopped the loop early, if the screen
// failed to flush.
func Stop() error {
	mu.Lock()
	defer mu.Unlock()
	if ctl == nil {
		return nil
	}
	send(stop)
	err := <-done
	ctl, done, paused = nil, nil, false
	core.Reset()
	return err
}

// send delivers c, unless the loop already ended on an error, which then
// stays in done for Stop.
func send(c command) {
	select {
	case ctl <- c:
	case err := <-done:
		done <- err
	}
}

func run(ctl <-chan command, done chan<- error) {
	last := time.Now()
	for {
		select {
		case c := <-ctl:
			for c == pause {
				c = <-ctl
			}
			if c == stop {
				done <- nil
				return
			}
			last = time.Now()
		default:
		}
		period := time.Second / time.Duration(max(FPS, 1))
		now := time.Now()
		dt := min(float32(now.Sub(last).Seconds()), 0.1)
		last = now
		if _, err := core.Frame(dt); err != nil {
			done <- err
			return
		}
		if spent := time.Since(now); spent < period {
			time.Sleep(period - spent)
		} else {
			// Still yield so polling goroutines get to run.
			time.Sleep(0)
		}
	}
}
