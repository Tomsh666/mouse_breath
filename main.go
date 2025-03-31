package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
	"unicode"
)

func main() {
	a := app.New()
	w := a.NewWindow("App name") //TODO: изменить название название

	nameOfTheExperimentLabel := widget.NewLabel("Name of experiment")
	nameOfTheExperimentEntry := widget.NewEntry()
	nameOfTheExperimentEntry.SetPlaceHolder("Enter name of experiment")

	durationOfTheExperimentLabel := widget.NewLabel("Time of measuring (sec):")
	durationOfTheExperimentEntry := widget.NewEntry()
	durationOfTheExperimentEntry.SetPlaceHolder("Enter time in seconds")

	durationOfTheExperimentEntry.OnChanged = func(s string) {
		var filtered strings.Builder
		for _, r := range s {
			if unicode.IsDigit(r) {
				filtered.WriteRune(r)
			}
		}
		if s != filtered.String() {
			durationOfTheExperimentEntry.SetText(filtered.String())
		}
	}

	buttonMeasure := widget.NewButton("Start Measuring", func() {
		startMeasuring(durationOfTheExperimentEntry)
	})

	buttonUploadToTable := widget.NewButton("Upload values to table", func() {
		uploadValuesToTable(nameOfTheExperimentEntry, durationOfTheExperimentEntry)
	})

	content := container.NewVBox(
		nameOfTheExperimentLabel,
		nameOfTheExperimentEntry,
		durationOfTheExperimentLabel,
		durationOfTheExperimentEntry,
		buttonMeasure,
		buttonUploadToTable,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 200))
	w.ShowAndRun()
}
