package gfx

import "tinygo.org/x/tinyfont"

// Font is a tinyfont bitmap font. Glyphs must be sorted by rune, which is
// how tinyfont's generator emits them.
type Font = tinyfont.Font

// Ascent is the distance from the top of the tallest glyph to the baseline.
func Ascent(f *Font) int16 { return -int16(f.BBox[3]) }

// LineHeight is the height a single line of text occupies.
func LineHeight(f *Font) int16 { return max(int16(f.YAdvance), int16(f.BBox[1])) }

// TextWidth is the advance width of s.
func TextWidth(f *Font, s string) int16 {
	var w int16
	for _, r := range s {
		w += int16(glyph(f, r).XAdvance)
	}
	return w
}

// glyph searches Glyphs directly: going through tinyfont's GetGlyph would box
// every glyph in an interface.
func glyph(f *Font, r rune) *tinyfont.Glyph {
	lo, hi := 0, len(f.Glyphs)
	for lo < hi {
		m := int(uint(lo+hi) >> 1)
		if f.Glyphs[m].Rune < r {
			lo = m + 1
		} else {
			hi = m
		}
	}
	if lo < len(f.Glyphs) && f.Glyphs[lo].Rune == r {
		return &f.Glyphs[lo]
	}
	return &f.EmptyGlyph
}

// Text draws s with its baseline at y and returns the x after the last glyph.
func (c Canvas) Text(f *Font, x, y int16, s string, col Color) int16 {
	for _, r := range s {
		g := glyph(f, r)
		gx := x + int16(g.XOffset)
		gy := y + int16(g.YOffset)
		if c.Visible(Rect{gx, gy, int16(g.Width), int16(g.Height)}) {
			c.glyph(g, gx, gy, col)
		}
		x += int16(g.XAdvance)
	}
	return x
}

// glyph blits a bitmap packed as one continuous MSB-first bit stream, rows
// not padded to bytes.
func (c Canvas) glyph(g *tinyfont.Glyph, x, y int16, col Color) {
	bit := 0
	w, h := int(g.Width), int(g.Height)
	for j := 0; j < h; j++ {
		for i := 0; i < w; i++ {
			if g.Bitmaps[bit>>3]&(0x80>>(bit&7)) != 0 {
				c.Pixel(x+int16(i), y+int16(j), col)
			}
			bit++
		}
	}
}
