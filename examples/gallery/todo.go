package gallery

import (
	"strconv"
	"time"

	"github.com/chhongzh/lilt/ui"
	"github.com/chhongzh/lilt/widget"
)

// todos is a to-do list made only of stock widgets. The field on top adds a
// task, OK on a task checks it off, and its Edit button opens a sheet to
// rename it or, held, to delete it.
func todos() ui.Node {
	t := ui.CurrentTheme
	tl := &todoList{
		list: padded(),
		summary: widget.NewLabel("", widget.Font(t.SmallFont), widget.TextColor(t.Muted),
			widget.TextAlign(0.5)),
	}
	add := widget.NewInput("New task")
	add.SetOnChange(func(s string) {
		if s == "" {
			return
		}
		tl.add(s)
		add.SetText("")
	})
	ui.Add(tl.list, add, tl.summary)
	for _, s := range []string{"Clear Blue Zenith", "Charge the battery", "Rewire the SET key"} {
		tl.add(s)
	}
	return Page("Todo", tl.list)
}

// todoList keeps the tasks between the field and the summary line, which
// counts what is left.
type todoList struct {
	list    *widget.List
	summary *widget.Label
	n, done int
}

func (tl *todoList) add(text string) {
	t := ui.CurrentTheme
	task := widget.NewCheckBox(text, false, ui.Grow(1))
	edit := widget.NewButton("Edit", widget.Font(t.SmallFont))
	row := ui.NewRow(ui.Gap(2), ui.Cross(ui.CrossCenter), ui.AnimateLayout(),
		ui.Children(task, edit))
	task.SetOnChange(func(done bool) {
		col := t.Foreground
		if done {
			col = t.Muted
			tl.done++
		} else {
			tl.done--
		}
		task.Title().SetColor(col)
		tl.refresh()
	})
	edit.SetOnClick(func() { tl.edit(row, task) })
	// Before the summary, which stays last.
	ui.Insert(tl.list, len(tl.list.Column().Children())-1, row)
	tl.n++
	tl.refresh()
}

func (tl *todoList) edit(row ui.Node, task *widget.CheckBox) {
	t := ui.CurrentTheme
	header := widget.NewLabel(task.Title().Text(), widget.Font(t.SmallFont),
		ui.Padding(ui.XY(2, 0)))
	rename := widget.NewButton("Rename", widget.OnClick(func() {
		ui.Pop()
		textInput().Open(taskTitle{task})
	}))
	del := widget.NewButton("Hold to delete", widget.HoldToConfirm(800*time.Millisecond),
		widget.OnClick(func() {
			// Popping first hands focus back to this row's Edit button, so
			// removing the row then moves it to a neighbour.
			ui.Pop()
			if task.Checked() {
				tl.done--
			}
			tl.n--
			ui.Remove(row)
			tl.refresh()
		}))
	ui.Push(ui.NewColumn(ui.Gap(2), ui.Padding(ui.LTRB(2, 2, 2, 2)),
		ui.Children(header, rename, del)), ui.Sheet)
}

func (tl *todoList) refresh() {
	switch {
	case tl.n == 0:
		tl.summary.SetText("Nothing to do")
	case tl.done == tl.n:
		tl.summary.SetText("All done")
	default:
		tl.summary.SetText(strconv.Itoa(tl.n-tl.done) + " of " + strconv.Itoa(tl.n) + " left")
	}
}

// taskTitle lets the keyboard edit a task's title; clearing it keeps the old one.
type taskTitle struct{ task *widget.CheckBox }

func (tt taskTitle) Text() string { return tt.task.Title().Text() }

func (tt taskTitle) SetText(s string) {
	if s != "" {
		tt.task.Title().SetText(s)
	}
}

// textInput is the editor Input would open: the app's, or the soft keyboard.
func textInput() ui.TextInputMethod {
	if m := ui.Current().TextInput; m != nil {
		return m
	}
	return widget.Keyboard
}
