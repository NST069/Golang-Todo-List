package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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

	var list *widget.List

	var data []TodoListItem
	var luri fyne.ListableURI
	path, perr := os.Getwd()
	if perr != nil {
		log.Fatal("Perr: ", perr)
	}
	luri, lerr := storage.ListerForURI(storage.NewFileURI(path))
	if lerr != nil {
		log.Fatal("Lerr: ", lerr, "\n\t", storage.NewFileURI(path))
	}

	file := fyne.NewMenu("File")
	importJSON := fyne.NewMenuItem("Import From JSON", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			_ = ReadJSON(reader.URI().Path(), &data)
			list.Refresh()
		}, w)
		fd.SetLocation(luri)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fd.Show()
	})
	exportJSON := fyne.NewMenuItem(("Export as JSON"), func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil {
				log.Println("Cancelled")
				return
			}
			WriteJSON(writer.URI().Path(), data)
		}, w)
		fd.SetDismissText("test")
		fd.SetFileName("TodoList.json")
		fd.SetLocation(luri)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fd.Show()
	})

	file.Items = append(file.Items, importJSON, exportJSON)
	menu := fyne.NewMainMenu(
		file,
	)
	w.SetMainMenu(menu)

	var currentItemId int
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
				tbText = ""
				data[currentItemId].Comment = entry.Text
				content.ParseMarkdown(data[currentItemId].Comment)
				card.SetContent(content)
				card.Refresh()
			}
		}),
		widget.NewToolbarAction(theme.CancelIcon(), func() {
			if tbIsEditing {
				tbIsEditing = false
				tbText = ""
				card.SetContent(content)
				card.Refresh()
			}
		}),
	)

	list = widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, widget.NewCheck("", func(value bool) {}), widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), nil), widget.NewLabel("dummy"))
		},
		func(id widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[1].(*widget.Check).Checked = data[id].Done
			o.(*fyne.Container).Objects[1].(*widget.Check).OnChanged = func(value bool) { data[id].Done = value }
			o.(*fyne.Container).Objects[0].(*widget.Label).SetText(data[id].Title)
			o.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				data = append(data[:id], data[id+1:]...)
				list.Refresh()
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		currentItemId = id
		card.SetTitle(data[id].Title)
		card.SetSubTitle(data[id].Ts.Format(time.RFC822))
		content.ParseMarkdown(data[id].Comment)
		card.Refresh()
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
			list.Refresh()
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
