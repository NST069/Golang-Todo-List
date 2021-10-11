package main

import (
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

	menu := fyne.NewMainMenu(
		fyne.NewMenu("File"),
	)
	w.SetMainMenu(menu)

	content := widget.NewRichTextFromMarkdown(``)
	content.Scroll = container.ScrollBoth
	card := widget.NewCard("Title", "Subtitle", content)

	tbIsEditing := false
	tbText := ""
	entry := widget.NewMultiLineEntry()
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			tbIsEditing = true
			tbText = content.String()
			entry.Text = tbText
			card.SetContent(entry)
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			if tbIsEditing {
				tbIsEditing = false
				content.ParseMarkdown(entry.Text)
				card.SetContent(content)
				tbText = ""
			}
		}),
		widget.NewToolbarAction(theme.CancelIcon(), func() {
			if tbIsEditing {
				tbIsEditing = false
				card.SetContent(content)
				tbText = ""
			}
		}),
	)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewCheck("", func(value bool) {}), widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), nil), widget.NewLabel("dummy"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Check).Checked = data[id].done
			item.(*fyne.Container).Objects[1].(*widget.Check).OnChanged = func(value bool) { data[id].done = value }
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(data[id].title)
			item.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() { data = append(data[:id], data[id+1:]...) }
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		card.SetTitle(data[id].title)
		card.SetSubTitle(data[id].ts.Format(time.RFC822))
		content.ParseMarkdown(data[id].comment)
	}
	list.OnUnselected = func(id widget.ListItemID) {
		card.SetTitle("Select An Item From The List")
	}

	addNew := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		title := widget.NewEntry()
		comment := widget.NewMultiLineEntry()
		items := []*widget.FormItem{
			widget.NewFormItem("Title", title),
			widget.NewFormItem("Comment", comment),
		}
		d := dialog.NewForm("Add new Item", "Add", "Cancel", items, func(b bool) {
			if !b {
				return
			}
			data = append(data, TodoListItem{time.Now(), false, title.Text, comment.Text})
		}, w)
		d.Resize(fyne.NewSize(400, 200))
		d.Show()
	})

	left := container.NewBorder(addNew, nil, nil, nil, list)
	right := container.NewBorder(toolbar, nil, nil, nil, card)
	split := container.NewHSplit(left, right)

	split.Offset = 0.4
	w.SetContent(split)
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()

}
