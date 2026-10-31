package gallery

import (
	"github.com/chhongzh/lilt/ui"

	"tinygo.org/x/tinyfont/freesans"
)

// FitFonts swaps the theme's fonts for FreeSans sizes close to the defaults
// times ui.ScaleFor(w, h). Fonts draw at their own pixel size, so a scaled
// screen needs larger ones; below a scale of 1.5 the defaults stay.
func FitFonts(w, h int16) {
	t := ui.CurrentTheme
	switch s := ui.ScaleFor(w, h); {
	case s >= 3.25:
		t.Font, t.SmallFont = &freesans.Regular24pt7b, &freesans.Regular18pt7b
	case s >= 2.25:
		t.Font, t.SmallFont = &freesans.Regular18pt7b, &freesans.Regular12pt7b
	case s >= 1.5:
		t.Font, t.SmallFont = &freesans.Regular12pt7b, &freesans.Regular9pt7b
	}
}
