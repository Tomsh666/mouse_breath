package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Обработка результатов")

	var selectedFile string

	fileLabel := widget.NewLabel("Файл не выбрана")

	btnSelect := widget.NewButton("Выбрать файл", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader != nil {
				selectedFile = reader.URI().Path()
				fileLabel.SetText("Выбран файл:\n" + selectedFile)
			}
		}, myWindow)
	})

	btnConvertExcel := widget.NewButton("Преобразовать в Excel", func() {
		if selectedFile == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите текстовый файл", myWindow)
			return
		}
		err := ConvertTxtToExcel(selectedFile, "result.xlsx")
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		dialog.ShowInformation("Готово", "Сохранено в файл: result.xlsx", myWindow)
	})

	intervalMinEntry := widget.NewEntry()
	intervalMaxEntry := widget.NewEntry()
	intervalMinEntry.SetPlaceHolder("Начальное значение")
	intervalMaxEntry.SetPlaceHolder("Конечное значение")

	resultLabel := widget.NewLabel("")

	calcFreqBtn := widget.NewButton("Посчитать частоту", func() {
		intervalMinStr := intervalMinEntry.Text
		intervalMaxStr := intervalMaxEntry.Text
		intervalMin, err := strconv.ParseFloat(intervalMinStr, 64)
		if err != nil || intervalMin < 0 {
			resultLabel.SetText("Введите корректный интервал")
			return
		}

		intervalMax, err := strconv.ParseFloat(intervalMaxStr, 64)
		if err != nil || intervalMax <= 0 || intervalMax <= intervalMin {
			resultLabel.SetText("Введите корректный интервал")
			return
		}

		data, err := ParseTimestamps(selectedFile)
		if err != nil {
			resultLabel.SetText("Ошибка чтения файла")
			return
		}
		count, freq := CountFrequency(data, intervalMin, intervalMax)

		resultLabel.SetText(fmt.Sprintf("Нажатий: %d\nЧастота: %.2f в сек", count, freq))
	})

	content := container.NewVBox(
		fileLabel,
		btnSelect,
		btnConvertExcel,
		intervalMinEntry,
		intervalMaxEntry,
		calcFreqBtn,
		resultLabel,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1000, 600))
	myWindow.ShowAndRun()

}
