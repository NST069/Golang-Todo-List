package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TodoList struct {
	List []TodoListItem
}

type TodoListItem struct {
	ts      time.Time
	done    bool
	title   string
	comment string
}

func GetIcon(done bool) fyne.Resource {
	if done {
		return theme.CheckButtonCheckedIcon()
	} else {
		return theme.CheckButtonIcon()
	}
}

func main() {
	a := app.NewWithID("todoList")
	w := a.NewWindow("Todo List")

	var data []TodoListItem
	//data = append(data, TodoListItem{time.Now(), false, "1", "000"})

	setDone := widget.NewLabel("")

	label := widget.NewLabel("Select An Item From The List")
	hbox := container.NewHBox(setDone, label)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewCheck("", func(value bool) {}), widget.NewButtonWithIcon("", theme.CancelIcon(), nil), widget.NewLabel("dummy"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Check).Checked = data[id].done
			item.(*fyne.Container).Objects[1].(*widget.Check).OnChanged = func(value bool) { data[id].done = value }
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(data[id].title)
			item.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() { data = append(data[:id], data[id+1:]...) }
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		str := data[id].ts.Format(time.RFC822) + "\t" + data[id].title + "\n" + data[id].comment
		label.SetText(str)
		setDone.SetText(strconv.FormatBool(data[id].done))
		//icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		setDone.SetText("")
		//icon.SetResource(nil)
	}

	addNew := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		title := widget.NewEntry()
		comment := widget.NewMultiLineEntry()
		items := []*widget.FormItem{
			widget.NewFormItem("Title", title),
			widget.NewFormItem("Comment", comment),
		}
		dialog.ShowForm("Add new Item", "Add", "Cancel", items, func(b bool) {
			if !b {
				return
			}
			data = append(data, TodoListItem{time.Now(), false, title.Text, comment.Text})
		}, w)
	})

	left := container.NewBorder(addNew, nil, nil, nil, list)

	split := container.NewHSplit(left, container.NewCenter(hbox))

	split.Offset = 0.4
	w.SetContent(split)
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()

}
