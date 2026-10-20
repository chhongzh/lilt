package widget

import "github.com/chhongzh/lilt/ui"

// Dialog is a box in the middle of the screen with a title, a one-line
// message and a row of buttons sharing its width. A button's OnClick
// handler runs as usual, loading dots and all; once it returns without an
// error the dialog closes. Back closes it as well.
type Dialog struct {
	ui.Flex
	title   *Label
	message *Label
	actions *ui.Flex
}

func NewDialog(title, message string, buttons ...*Button) *Dialog {
	t := ui.CurrentTheme
	d := &Dialog{
		title:   NewLabel(title, TextAlign(0.5)),
		message: NewLabel(message, Font(t.SmallFont), TextAlign(0.5), MarqueeMode(MarqueeAlways)),
		actions: ui.NewRow(ui.Gap(t.Gap)),
	}
	d.SetAxis(ui.Vertical)
	d.SetGap(t.Gap)
	d.SetPadding(ui.All(3))
	d.SetWidth(ui.Rel(1))
	for _, b := range buttons {
		b.SetGrow(1)
		b.after = func(err error) {
			if err == nil {
				d.Close()
			}
		}
		ui.Attach(d.actions, b)
	}
	ui.Attach(d, d.title)
	if message != "" {
		ui.Attach(d, d.message)
	}
	ui.Attach(d, d.actions)
	return d
}

func (d *Dialog) Title() *Label   { return d.title }
func (d *Dialog) Message() *Label { return d.message }

// Open shows the dialog, its first button focused. Call it on the UI
// goroutine.
func (d *Dialog) Open() { ui.Open(d, ui.Dialog) }

func (d *Dialog) Close() { ui.Close(d) }

// toastHold is how long a toast stays up.
const toastHold = 2

// toast is the label Toast shows.
type toast struct {
	Label
	timer float32
}

var lastToast *toast

// Toast shows text near the bottom of the screen for two seconds. It takes
// neither focus nor keys, and replaces a toast still showing. Call it on the
// UI goroutine.
func Toast(text string) {
	if lastToast != nil {
		ui.Close(lastToast)
	}
	t := &toast{}
	t.text = text
	t.SetFont(ui.CurrentTheme.SmallFont)
	t.SetMarquee(MarqueeAlways)
	t.SetPadding(ui.CurrentTheme.RowPadding)
	lastToast = t
	ui.Open(t, ui.Toast)
}

func (t *toast) Update(dt float32) {
	t.Label.Update(dt)
	t.timer += dt
	if t.timer >= toastHold && t.timer-dt < toastHold {
		ui.Close(t)
		if lastToast == t {
			lastToast = nil
		}
	}
}
