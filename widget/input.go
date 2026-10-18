package widget

import (
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Input is a one-line text field. Activating it opens the app's
// TextInputMethod, or the built-in Keyboard when none is set. Characters
// from a hardware keyboard are typed straight in.
type Input struct {
	ui.Box
	label       *Label
	text        string
	placeholder string
	onChange    func(string)
}

func NewInput(placeholder string, opts ...ui.Option) *Input {
	in := &Input{placeholder: placeholder, label: NewLabel("")}
	in.label.SetWidth(ui.Rel(1))
	ui.Attach(in, in.label)
	t := ui.CurrentTheme
	in.SetFocusable(true)
	in.SetBorder(t.Foreground)
	in.SetRadius(t.Radius)
	in.SetPadding(ui.XY(4, 1))
	in.refresh()
	ui.Apply(in, opts)
	return in
}

func (in *Input) Text() string { return in.text }

func (in *Input) SetText(s string) {
	if s == in.text {
		return
	}
	in.text = s
	in.refresh()
	if in.onChange != nil {
		in.onChange(s)
	}
}

func (in *Input) SetOnChange(fn func(string)) { in.onChange = fn }

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

func (in *Input) HandleKey(e *ui.KeyEvent) bool {
	switch e.Key {
	case ui.KeyOK:
		if e.Phase == ui.Release && !e.Canceled {
			m := ui.Current().TextInput
			if m == nil {
				m = Keyboard
			}
			m.Open(in)
		}
		return true
	case ui.KeyChar:
		in.SetText(in.text + string(e.Rune))
		return true
	case ui.KeyBackspace:
		if e.Phase != ui.Release && in.text != "" {
			in.SetText(dropLast(in.text))
		}
		return true
	}
	return false
}

// dropLast removes the final UTF-8 character.
func dropLast(s string) string {
	i := len(s) - 1
	for i > 0 && s[i]&0xC0 == 0x80 {
		i--
	}
	return s[:i]
}
