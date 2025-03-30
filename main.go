package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("App name")
	w.Resize(fyne.NewSize(300, 200))

	buttonMeasure := widget.NewButton("Start Measuring", StartMeasuring)

	buttonUploadToTable := widget.NewButton("Upload values to table", UploadValuesToTable)

	content := container.NewVBox(
		buttonMeasure,
		buttonUploadToTable,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
