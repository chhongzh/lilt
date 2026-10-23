package widget

import (
	"time"

	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Capabilities are interfaces: an option applies to any node with the
// method it needs, whatever the node embeds. Widgets get the methods by
// embedding Label, ListItem, Control and the like.
type (
	// Texter is a widget with a line of text, such as Label.
	Texter interface {
		Text() string
		SetText(s string)
	}
	// FontSetter takes a font for its text.
	FontSetter interface{ SetFont(f *gfx.Font) }
	// Colorer takes a text colour.
	Colorer interface{ SetColor(c gfx.Color) }
	// Aligner places its text horizontally, 0 left to 1 right.
	Aligner interface{ SetAlign(x float32) }
	// Trailer shows a node at the end of its row.
	Trailer interface{ SetTrailing(d ui.Node) }
)

// Font sets the font of a Label, or of the label inside a widget.
func Font(f *gfx.Font) ui.Option {
	return func(d ui.Node) { must[FontSetter](d, "Font").SetFont(f) }
}

func TextColor(c gfx.Color) ui.Option {
	return func(d ui.Node) { must[Colorer](d, "TextColor").SetColor(c) }
}

// TextAlign places text horizontally, 0 left to 1 right.
func TextAlign(x float32) ui.Option {
	return func(d ui.Node) { must[Aligner](d, "TextAlign").SetAlign(x) }
}

func MarqueeMode(m Marquee) ui.Option {
	return func(d ui.Node) { must[*Label](d, "MarqueeMode").SetMarquee(m) }
}

// OnClick sets the handler run when the widget is activated. It runs on its
// own goroutine; see Control.
func OnClick[F ClickFunc](fn F) ui.Option {
	return func(d ui.Node) { must[controller](d, "OnClick").control().OnClick(fn) }
}

// OnChange sets the handler run with the old and new value after the user
// changes it. It runs on its own goroutine; see Control.
func OnChange[F ChangeFunc](fn F) ui.Option {
	return func(d ui.Node) { must[controller](d, "OnChange").control().OnChange(fn) }
}

// HoldToConfirm makes a Button fire only after OK is held for dur.
func HoldToConfirm(dur time.Duration) ui.Option {
	return func(d ui.Node) { must[*Button](d, "HoldToConfirm").SetHold(dur) }
}

// Trailing puts d at the end of a ListItem row.
func Trailing(t ui.Node) ui.Option {
	return func(d ui.Node) { must[Trailer](d, "Trailing").SetTrailing(t) }
}

// controller is any widget embedding Control.
type controller interface{ control() *Control }

func must[T any](d ui.Node, opt string) T {
	t, ok := d.(T)
	if !ok {
		panic("widget: option " + opt + " does not apply to this node")
	}
	return t
}
