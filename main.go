package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func saveNotes(notes []string, path string) {
	data, _ := json.MarshalIndent(notes, "", "  ")
	_ = ioutil.WriteFile(path, data, 0644)
}

func loadNotes(path string) []string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []string{}
	}
	var notes []string
	_ = json.Unmarshal(data, &notes)
	return notes
}

func main() {
	a := app.NewWithID("todo.notes.app")
	w := a.NewWindow("📝 Notes ToDo")
	w.Resize(fyne.NewSize(420, 600))

	storageDir := a.Storage().RootURI().Path()
	notesPath := filepath.Join(storageDir, "notes.json")

	var notes []string = loadNotes(notesPath)
	var list *widget.List

	list = widget.NewList(
		func() int { return len(notes) },
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			label.Wrapping = fyne.TextWrapWord
			label.Alignment = fyne.TextAlignLeading

			del := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			del.Importance = widget.DangerImportance
			del.Resize(fyne.NewSize(40, 40))

			row := container.NewBorder(nil, nil, nil, del, label)
			return container.NewVBox(row)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if i >= len(notes) {
				return
			}
			box := o.(*fyne.Container)
			row := box.Objects[0].(*fyne.Container)

			label := row.Objects[0].(*widget.Label)
			del := row.Objects[1].(*widget.Button)

			label.SetText(notes[i])
			del.OnTapped = func() {
				notes = append(notes[:i], notes[i+1:]...)
				saveNotes(notes, notesPath)
				list.Refresh()
			}
		},
	)

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Dodaj novu belešku...")

	addBtn := widget.NewButtonWithIcon("Dodaj", theme.ContentAddIcon(), func() {
		text := strings.TrimSpace(entry.Text)
		if text == "" {
			return
		}
		notes = append(notes, text)
		saveNotes(notes, notesPath)
		entry.SetText("")
		list.Refresh()
	})
	addBtn.Importance = widget.HighImportance

	inputRow := container.NewHBox(
		container.NewGridWrap(fyne.NewSize(280, 40), entry),
		container.NewGridWrap(fyne.NewSize(100, 40), addBtn),
	)

	title := canvas.NewText("📋 Vaše beleške", theme.PrimaryColor())
	title.TextSize = 22
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	divider := canvas.NewLine(theme.ShadowColor())
	divider.StrokeWidth = 1

	scroll := container.NewVScroll(list)
	scroll.SetMinSize(fyne.NewSize(400, 400))

	content := container.NewVBox(
		title,
		divider,
		inputRow,
		scroll,
	)

	w.SetContent(container.NewPadded(content))
	w.ShowAndRun()
}
