package widget

import (
	"time"

	"osu-mc/gfx"
	"osu-mc/ui"
)

// Font sets the font of a Label, or of the label inside a widget.
func Font(f *gfx.Font) ui.Option {
	return func(d ui.Drawable) { must[interface{ SetFont(*gfx.Font) }](d, "Font").SetFont(f) }
}

func TextColor(c gfx.Color) ui.Option {
	return func(d ui.Drawable) { must[interface{ SetColor(gfx.Color) }](d, "TextColor").SetColor(c) }
}

// TextAlign places text horizontally, 0 left to 1 right.
func TextAlign(x float32) ui.Option {
	return func(d ui.Drawable) { must[interface{ SetAlign(float32) }](d, "TextAlign").SetAlign(x) }
}

func MarqueeMode(m Marquee) ui.Option {
	return func(d ui.Drawable) { must[interface{ SetMarquee(Marquee) }](d, "MarqueeMode").SetMarquee(m) }
}

// OnClick runs fn when the widget is activated.
func OnClick(fn func()) ui.Option {
	return func(d ui.Drawable) { must[interface{ SetOnClick(func()) }](d, "OnClick").SetOnClick(fn) }
}

// OnAction runs fn on its own goroutine when a Button is activated and shows
// loading, then success or failure depending on the returned error.
func OnAction(fn func() error) ui.Option {
	return func(d ui.Drawable) { must[*Button](d, "OnAction").SetOnAction(fn) }
}

// HoldToConfirm makes a Button fire only after OK is held for dur.
func HoldToConfirm(dur time.Duration) ui.Option {
	return func(d ui.Drawable) { must[*Button](d, "HoldToConfirm").SetHold(dur) }
}

// Trailing puts d at the end of a ListItem row.
func Trailing(t ui.Drawable) ui.Option {
	return func(d ui.Drawable) {
		must[interface{ SetTrailing(ui.Drawable) }](d, "Trailing").SetTrailing(t)
	}
}

func must[T any](d ui.Drawable, opt string) T {
	t, ok := d.(T)
	if !ok {
		panic("widget: option " + opt + " does not apply to this node")
	}
	return t
}

// OnChange registers a value callback: bool for CheckBox, int for Select,
// string for Input.
func OnChange[T any](fn func(T)) ui.Option {
	return func(d ui.Drawable) { must[interface{ SetOnChange(func(T)) }](d, "OnChange").SetOnChange(fn) }
}
