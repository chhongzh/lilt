package widget

import (
	"osu-mc/gfx"
	"osu-mc/ui"
)

// sheetMax is the share of the screen height a Select sheet may take.
const sheetMax = 0.8

// Select is a list row showing the chosen option; activating it opens a
// sheet of radios.
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
	s.value = NewLabel(s.optionText(), Font(ui.CurrentTheme.SmallFont), MarqueeMode(MarqueeOff))
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
	t := ui.CurrentTheme
	_, sh := ui.Current().Size()
	header := NewLabel(s.title.Text(), Font(t.SmallFont), ui.Padding(ui.XY(4, 1)))
	list := NewList()
	group := NewRadioGroup(s.selected, func(i int) {
		s.SetSelected(i)
		if s.onChange != nil {
			s.onChange(i)
		}
		ui.Pop()
	})
	for _, o := range s.options {
		ui.Add(list, NewRadio(group, o))
	}
	row := float32(gfx.LineHeight(t.Font)) + 2
	head := float32(gfx.LineHeight(t.SmallFont)) + 2
	h := min(row*float32(len(s.options)), float32(sh)*sheetMax-head)
	list.SetHeight(ui.Px(h))
	sheet := ui.NewColumn(ui.Padding(ui.LTRB(0, 2, 0, 1)), ui.Children(header, list))
	ui.Push(sheet, ui.Sheet)
	if s.selected >= 0 && s.selected < len(group.Radios()) {
		ui.SetFocus(group.Radios()[s.selected])
	}
}
