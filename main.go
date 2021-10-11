package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TodoListItem struct {
	Token   []byte
	Ts      time.Time
	Done    bool
	Title   string
	Comment string
}

func WriteJSON(filename string, items []TodoListItem) {
	b, err := json.Marshal(items)
	if err != nil {
		fmt.Println("Error Marshaling to json: ", err)
	} else {
		err2 := os.WriteFile(filename, b, 0666)
		if err2 != nil {
			fmt.Println("Error writing to file: ", err)
		}
	}
}

func ReadJSON(filename string, v interface{}) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading from file: ", err)
	} else {
		err2 := json.Unmarshal(b, &v)
		if err2 != nil {
			fmt.Println("Err2: ", err2)
			return err2
		}
	}
	return nil
}

func main() {
	a := app.NewWithID("todoList")
	w := a.NewWindow("Todo List")

	var data []TodoListItem

	path, err := os.Getwd()
	filename := ""
	if err == nil {
		filename = path + "\\file.json"
	}
	file := fyne.NewMenu("File")
	importJSON := fyne.NewMenuItem("Import From JSON", func() {
		_ = ReadJSON(filename, &data)
	})
	exportJSON := fyne.NewMenuItem(("Export as JSON"), func() {
		WriteJSON(filename, data)
	})
	file.Items = append(file.Items, importJSON, exportJSON)
	menu := fyne.NewMainMenu(
		file,
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
			item.(*fyne.Container).Objects[1].(*widget.Check).Checked = data[id].Done
			item.(*fyne.Container).Objects[1].(*widget.Check).OnChanged = func(value bool) { data[id].Done = value }
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(data[id].Title)
			item.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() { data = append(data[:id], data[id+1:]...) }
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		card.SetTitle(data[id].Title)
		card.SetSubTitle(data[id].Ts.Format(time.RFC822))
		content.ParseMarkdown(data[id].Comment)
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
			ts := time.Now()
			hash := sha256.Sum256([]byte(ts.String()))
			data = append(data, TodoListItem{hash[:], ts, false, title.Text, comment.Text})
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
