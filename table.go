package main

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
	"time"
)

func safeSetCell(f *excelize.File, sheet, cell string, value any) {
	if err := f.SetCellValue(sheet, cell, value); err != nil {
		fmt.Printf("Ошибка записи в %s: %v\n", cell, err)
	}
}

func uploadValuesToTable(nameEntry *widget.Entry, durationEntry *widget.Entry) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := nameEntry.Text
	if len(sheet) == 0 || len(sheet) > 31 {
		fmt.Println("Ошибка: некорректное имя листа")
		return
	}
	if err := f.SetSheetName(f.GetSheetName(0), sheet); err != nil {
		fmt.Println("Ошибка переименования листа:", err)
		return
	}

	data := map[string]any{
		"A1": "Время измерения, сек",
		"B1": durationEntry.Text,
		"A2": "Начало",
		"A3": "Конец",
		"A4": "Время, сек",
		"B4": "Частота, вдох./сек.",
	}

	for cell, value := range data {
		safeSetCell(f, sheet, cell, value)
	}

	currentTime := time.Now()
	timeStr := currentTime.Format("15:04:05")
	fmt.Println("Начальное время:", timeStr)
	safeSetCell(f, sheet, "B2", timeStr)

	if len(measurementData) == 0 {
		fmt.Println("Ошибка: нет данных для записи")
		return
	}

	for i, row := range measurementData[0] {
		safeSetCell(f, sheet, fmt.Sprintf("A%d", 5+i), row)
		safeSetCell(f, sheet, fmt.Sprintf("B%d", 5+i), measurementData[1][i])
	}

	timeStr = currentTime.Format("15:04:05")
	fmt.Println("Конечное время:", timeStr)
	safeSetCell(f, sheet, "B3", timeStr)

	fmt.Println("Uploading")

	filename := sheet + ".xlsx"
	if err := f.SaveAs(filename); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	// TODO: Изменит - Время конца определяется по реальному времени, когда отработает прога => сейчас время начала и коца одинаковые
}
