package widget

import (
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Select is a list row showing the chosen option; activating it opens a
// bottom Drawer of radios. Picking one closes the drawer and reports the
// change through OnChange, while the row shows its loading; a handler
// that fails puts the old option back.
type Select struct {
	ListItem
	value    *Label
	options  []string
	selected int
}

func NewSelect(title string, options []string, selected int, opts ...ui.Option) *Select {
	s := &Select{options: options, selected: selected}
	s.init(title)
	s.kind = 0
	s.value = NewLabel(s.optionText(), Font(ui.CurrentTheme.SmallFont))
	s.SetTrailing(s.value)
	s.activate = s.open
	ui.Apply(s, opts)
	return s
}

func (s *Select) Selected() int { return s.selected }

// SetSelected shows option i without calling OnChange.
func (s *Select) SetSelected(i int) {
	s.selected = i
	s.value.SetText(s.optionText())
}

func (s *Select) optionText() string {
	if s.selected < 0 || s.selected >= len(s.options) {
		return ""
	}
	return s.options[s.selected]
}

func (s *Select) open() {
	d := NewDrawer(Bottom, s.title.Text())
	list := NewList()
	group := NewRadioGroup(s.selected)
	group.picked = func(from, to int) {
		d.Close()
		s.SetSelected(to)
		s.run(changeCall(s.change, from, to), func(err error) {
			if err != nil {
				s.SetSelected(from)
				shake(s)
			}
		})
	}
	for _, o := range s.options {
		ui.Add(list, NewRadio(group, o))
	}
	// The list is sized to its rows so a short one does not fill the drawer.
	_, sh := ui.ScreenSize()
	row := float32(gfx.LineHeight(ui.CurrentTheme.Font)) + ui.CurrentTheme.RowPadding.Vertical()*ui.Unit()
	room := float32(sh)*0.8 - d.TitleHeight() - d.Padding().Vertical()*ui.Unit()
	list.SetHeight(ui.Px(min(row*float32(len(s.options)), room)))
	ui.Add(d, list)
	d.Open()
	if s.selected >= 0 && s.selected < len(group.Radios()) {
		ui.SetFocus(group.Radios()[s.selected])
	}
}
