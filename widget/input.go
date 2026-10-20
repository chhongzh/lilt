package widget

import (
	"github.com/chhongzh/lilt/anim"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Input is a one-line text field. Activating it opens ui.TextInput, or the
// built-in Keyboard when none is set. Characters from a hardware keyboard
// are typed straight in.
//
// Edits are reported through OnChange with the text before and after. Edits
// made while a handler runs are reported together once it returns; one that
// fails puts back the text it was given.
type Input struct {
	ui.Box
	Control
	label       *Label
	text        string
	reported    string
	placeholder string
	waiting     bool
}

func NewInput(placeholder string, opts ...ui.Option) *Input {
	in := &Input{placeholder: placeholder, label: NewLabel("")}
	in.kind = ""
	in.label.SetWidth(ui.Rel(1))
	ui.Attach(in, in.label)
	t := ui.CurrentTheme
	in.SetFocusable(true)
	in.SetBorder(t.Foreground)
	in.SetRadius(t.Radius)
	in.SetPadding(t.RowPadding)
	in.OffsetX().SetSpring(anim.Bouncy)
	in.refresh()
	ui.Apply(in, opts)
	return in
}

func (in *Input) Text() string { return in.text }

// SetText replaces the text without calling OnChange.
func (in *Input) SetText(s string) {
	in.text, in.reported = s, s
	in.refresh()
}

func (in *Input) SetFont(f *gfx.Font) { in.label.SetFont(f) }

func (in *Input) refresh() {
	t := ui.CurrentTheme
	if in.text == "" {
		in.label.SetText(in.placeholder)
		in.label.SetColor(t.Muted)
		return
	}
	in.label.SetText(in.text)
	in.label.SetColor(t.Foreground)
}

// edit is a change made by the user, which OnChange hears about.
func (in *Input) edit(s string) {
	in.text = s
	in.refresh()
	in.report()
}

func (in *Input) report() {
	if in.Busy() || in.text == in.reported {
		return
	}
	from, to := in.reported, in.text
	in.reported = to
	in.run(changeCall(in.change, from, to), func(err error) {
		if err != nil {
			in.SetText(from)
			shake(in)
			return
		}
		in.report()
	})
}

// inputTarget is what the editor gets, so its result counts as an edit.
type inputTarget struct{ in *Input }

func (t inputTarget) Text() string     { return t.in.text }
func (t inputTarget) SetText(s string) { t.in.edit(s) }

func (in *Input) HandleKey(e *ui.KeyEvent) bool {
	switch e.Key {
	case ui.KeyOK:
		if e.Phase == ui.Release && !e.Canceled {
			in.run(clickCall(in.click), func(error) { in.report() })
			m := ui.TextInput
			if m == nil {
				m = Keyboard
			}
			m.Open(inputTarget{in})
		}
		return true
	case ui.KeyChar:
		in.edit(in.text + string(e.Rune))
		return true
	case ui.KeyBackspace:
		if e.Phase != ui.Release && in.text != "" {
			in.edit(dropLast(in.text))
		}
		return true
	}
	return false
}

func (in *Input) Update(dt float32) {
	in.tick(dt)
	if in.Loading() != in.waiting {
		in.waiting = !in.waiting
		ui.Invalidate()
	}
}

func (in *Input) DrawOver(c gfx.Canvas) {
	if !in.waiting {
		return
	}
	r := in.ContentRect()
	drawWaiting(c, ui.Round(r.Right()-9*ui.Unit()), ui.Round(r.Center().Y), in.elapsed)
}

// dropLast removes the final UTF-8 character.
func dropLast(s string) string {
	i := len(s) - 1
	for i > 0 && s[i]&0xC0 == 0x80 {
		i--
	}
	return s[:i]
}
