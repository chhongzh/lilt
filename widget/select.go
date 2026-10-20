package widget

import (
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
)

// Select is a list row showing the chosen option; activating it opens a
// bottom Drawer of radios. Picking one closes the drawer.
type Select struct {
	ListItem
	value    *Label
	options  []string
	selected int
	onChange func(int)
}

func NewSelect(title string, options []string, selected int, opts ...ui.Option) *Select {
	s := &Select{options: options, selected: selected}
	s.init(title)
	s.value = NewLabel(s.optionText(), Font(ui.CurrentTheme.SmallFont))
	s.SetTrailing(s.value)
	s.activate = s.open
	ui.Apply(s, opts)
	return s
}

func (s *Select) Selected() int { return s.selected }

func (s *Select) SetSelected(i int) {
	s.selected = i
	s.value.SetText(s.optionText())
}

func (s *Select) SetOnChange(fn func(int)) { s.onChange = fn }

func (s *Select) optionText() string {
	if s.selected < 0 || s.selected >= len(s.options) {
		return ""
	}
	return s.options[s.selected]
}

func (s *Select) open() {
	d := NewDrawer(Bottom, s.title.Text())
	list := NewList()
	group := NewRadioGroup(s.selected, nil)
	for i, o := range s.options {
		ui.Add(list, NewRadio(group, o, OnClick(func() {
			d.Close()
			if i != s.selected {
				s.SetSelected(i)
				if s.onChange != nil {
					s.onChange(i)
				}
			}
		})))
	}
	// The list is sized to its rows so a short one does not fill the drawer.
	_, sh := ui.ScreenSize()
	row := float32(gfx.LineHeight(ui.CurrentTheme.Font)) + 2
	room := float32(sh)*0.8 - d.TitleHeight() - d.Padding().Vertical()
	list.SetHeight(ui.Px(min(row*float32(len(s.options)), room)))
	ui.Add(d, list)
	d.Open()
	if s.selected >= 0 && s.selected < len(group.Radios()) {
		ui.SetFocus(group.Radios()[s.selected])
	}
}
