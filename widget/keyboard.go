package widget

import (
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Keyboard is the TextInputMethod Input falls back to.
var Keyboard = &SoftKeyboard{}

type keyKind uint8

const (
	keyChar keyKind = iota
	keySpace
	keyBackspace
	keyShift
	keyLayer
	keyDone
)

// SoftKeyboard is an on-screen keyboard shown as a sheet. Arrow keys move
// between keys with the focus finder, so rows of different lengths line up
// naturally.
type SoftKeyboard struct {
	target  ui.TextTarget
	buf     []rune
	shift   bool
	symbols bool
	preview *Label
	letters *ui.Flex
	digits  *ui.Flex
}

var (
	letterRows = [...]string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
	symbolRows = [...]string{"1234567890", "-/:;()$&@", "=.,?!'\"#"}
)

func (k *SoftKeyboard) Open(t ui.TextTarget) {
	th := ui.CurrentTheme
	k.target, k.shift, k.symbols = t, false, false
	k.buf = []rune(t.Text())
	k.preview = NewLabel("", Font(th.SmallFont),
		ui.Padding(ui.XY(2, 1)), ui.Border(th.Foreground), ui.Radius(2))
	k.letters = k.layer(letterRows, k.key(keyShift, "^"), k.key(keyLayer, "123"))
	k.digits = k.layer(symbolRows, nil, k.key(keyLayer, "abc"))
	k.digits.SetVisible(false)
	sheet := &kbSheet{k: k}
	sheet.SetAxis(ui.Vertical)
	sheet.SetGap(1)
	ui.Add(sheet, k.preview, k.letters, k.digits)
	ui.Open(sheet, ui.DrawerBottom)
	k.refresh()
}

// layer builds three rows of characters and the space row. lead, if set,
// starts the third row; the third row always ends with backspace.
func (k *SoftKeyboard) layer(rows [3]string, lead, toggle *key) *ui.Flex {
	col := ui.NewColumn(ui.Gap(1))
	for i, chars := range rows {
		row := ui.NewRow(ui.Gap(1))
		if i == 2 && lead != nil {
			ui.Add(row, lead)
		}
		for _, r := range chars {
			ui.Add(row, k.char(r))
		}
		if i == 2 {
			ui.Add(row, k.key(keyBackspace, "<"))
		}
		ui.Add(col, row)
	}
	space := k.key(keySpace, "space")
	space.SetGrow(4)
	toggle.SetGrow(1.5)
	done := k.key(keyDone, "OK")
	done.SetGrow(1.5)
	ui.Add(col, ui.NewRow(ui.Gap(1), ui.Children(toggle, space, done)))
	return col
}

func (k *SoftKeyboard) char(r rune) *key {
	up := r
	if r >= 'a' && r <= 'z' {
		up = r - 'a' + 'A'
	}
	c := k.key(keyChar, string(r))
	c.r, c.upper = r, string(up)
	return c
}

func (k *SoftKeyboard) key(kind keyKind, text string) *key {
	c := &key{k: k, kind: kind, text: text, upper: text}
	c.SetFocusable(true)
	c.SetGrow(1)
	return c
}

func (k *SoftKeyboard) press(c *key) {
	switch c.kind {
	case keyChar:
		r := c.r
		if k.shift && r >= 'a' && r <= 'z' {
			r = r - 'a' + 'A'
			k.shift = false
		}
		k.buf = append(k.buf, r)
	case keySpace:
		k.buf = append(k.buf, ' ')
	case keyBackspace:
		if len(k.buf) > 0 {
			k.buf = k.buf[:len(k.buf)-1]
		}
	case keyShift:
		k.shift = !k.shift
	case keyLayer:
		k.symbols = !k.symbols
		k.letters.SetVisible(!k.symbols)
		k.digits.SetVisible(k.symbols)
	case keyDone:
		k.target.SetText(string(k.buf))
		ui.Pop()
		return
	}
	k.refresh()
}

func (k *SoftKeyboard) refresh() {
	s := string(k.buf) + "_"
	k.preview.SetText(s)
	// Keep the caret in view once the text outgrows the box.
	fits := float32(gfx.TextWidth(k.preview.Font(), s)) <= k.preview.Rect().W-4
	if fits {
		k.preview.SetAlign(0)
	} else {
		k.preview.SetAlign(1)
	}
	ui.Invalidate()
}

// kbSheet takes characters and backspace from a hardware keyboard too.
type kbSheet struct {
	ui.Flex
	k *SoftKeyboard
}

func (s *kbSheet) HandleKey(e *ui.KeyEvent) bool {
	switch {
	case e.Key == ui.KeyChar:
		s.k.buf = append(s.k.buf, e.Rune)
	case e.Key == ui.KeyBackspace && e.Phase != ui.Release:
		if len(s.k.buf) > 0 {
			s.k.buf = s.k.buf[:len(s.k.buf)-1]
		}
	default:
		return false
	}
	s.k.refresh()
	return true
}

type key struct {
	ui.Box
	k           *SoftKeyboard
	kind        keyKind
	r           rune
	text, upper string
}

func (c *key) Leaf() {}

func (c *key) Measure(ui.Vec2) ui.Vec2 {
	f := ui.CurrentTheme.SmallFont
	return ui.Vec2{X: float32(gfx.TextWidth(f, c.text)) + 2, Y: float32(gfx.LineHeight(f)) + 1}
}

func (c *key) Draw(cv gfx.Canvas) {
	f := ui.CurrentTheme.SmallFont
	s := c.text
	if c.k.shift {
		s = c.upper
	}
	r := c.Rect().Pixels()
	x := (r.W - gfx.TextWidth(f, s)) / 2
	y := (r.H-gfx.LineHeight(f))/2 + gfx.Ascent(f)
	col := ui.CurrentTheme.Foreground
	cv.Text(f, x, y, s, col)
	if c.kind == keyShift && c.k.shift {
		cv.HLine(x, r.H-1, gfx.TextWidth(f, s), col)
	}
}

func (c *key) HandleKey(e *ui.KeyEvent) bool {
	if e.Key != ui.KeyOK {
		return false
	}
	if e.Phase == ui.Release && !e.Canceled {
		c.k.press(c)
	}
	return true
}
