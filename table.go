package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ConvertTxtToExcel(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	f := excelize.NewFile()
	sheet := "Sheet1"

	scanner := bufio.NewScanner(file)

	var totalSeconds int
	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		totalSeconds, _ = strconv.Atoi(line)
	}

	f.SetCellValue(sheet, "A1", "№")
	f.SetCellValue(sheet, "B1", "Метка времени")
	f.SetCellValue(sheet, "C1", "Время измерения, в секундах")
	f.SetCellValue(sheet, "C2", totalSeconds)

	row := 1
	index := 0
	var timestamps []float64

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		timeVal, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		timestamps = append(timestamps, timeVal)
		row++
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), index)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), timeVal)
		index++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(timestamps) > 0 {
		sort.Float64s(timestamps)

		var deltas []float64
		for i := 1; i < len(timestamps); i++ {
			deltas = append(deltas, timestamps[i]-timestamps[i-1])
		}

		avgDelta := 0.0
		minDelta := 0.0
		maxDelta := 0.0

		if len(deltas) > 0 {
			sum := 0.0
			minDelta = deltas[0]
			maxDelta = deltas[0]
			for _, d := range deltas {
				sum += d
				if d < minDelta {
					minDelta = d
				}
				if d > maxDelta {
					maxDelta = d
				}
			}
			avgDelta = sum / float64(len(deltas))
		}

		f.SetCellValue(sheet, "D1", "Статистика")
		f.SetCellValue(sheet, "D2", "Всего нажатий")
		f.SetCellValue(sheet, "E2", len(timestamps))

		f.SetCellValue(sheet, "D3", "Минимальный интервал (сек)")
		f.SetCellValue(sheet, "E3", minDelta)

		f.SetCellValue(sheet, "D4", "Максимальный интервал (сек)")
		f.SetCellValue(sheet, "E4", maxDelta)

		f.SetCellValue(sheet, "D5", "Средний интервал между нажатиями (сек)")
		f.SetCellValue(sheet, "E5", avgDelta)
	}

	return f.SaveAs(outputPath)
}

func ParseTimestamps(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var timestamps []float64
	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		_ = scanner.Text()
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		sec, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		timestamps = append(timestamps, sec)
	}
	return timestamps, scanner.Err()
}

func CountFrequency(data []float64, from, to float64) (int, float64) {
	startIdx := -1
	endIdx := -1

	for i, v := range data {
		if v >= from {
			startIdx = i
			break
		}
	}

	for i := len(data) - 1; i >= 0; i-- {
		if data[i] <= to {
			endIdx = i
			break
		}
	}

	if startIdx == -1 || endIdx == -1 || endIdx < startIdx {
		return 0, 0.0
	}

	count := float64(endIdx - startIdx + 1)
	duration := data[endIdx] - data[startIdx]

	if duration <= 0 {
		return int(count), 0.0
	}

	freq := count / duration

	return int(count), freq
}

func ConvertTxtToGoogle(path string) error {

	google_token_path := "C:\\x\\x..."
	ctx := context.Background()
	creds, err := os.ReadFile(google_token_path)
	if err != nil {
		return fmt.Errorf("не удалось прочитать credentials.json: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("не удалось создать конфигурацию: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("не удалось создать сервис Google Sheets: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("не удалось открыть текстовый файл: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var values [][]interface{}
	var timestamps []float64
	var totalSeconds int

	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		totalSeconds, err = strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("не удалось прочитать время измерения: %v", err)
		}
	}

	values = append(values, []interface{}{"№", "Метка времени", "Время измерения, в секундах"})
	values = append(values, []interface{}{"", "", totalSeconds})

	index := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		timeVal, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		timestamps = append(timestamps, timeVal)
		values = append(values, []interface{}{index, timeVal, ""})
		index++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	if len(timestamps) > 0 {
		sort.Float64s(timestamps)
		var deltas []float64
		for i := 1; i < len(timestamps); i++ {
			deltas = append(deltas, timestamps[i]-timestamps[i-1])
		}

		avgDelta := 0.0
		minDelta := 0.0
		maxDelta := 0.0

		if len(deltas) > 0 {
			sum := 0.0
			minDelta = deltas[0]
			maxDelta = deltas[0]
			for _, d := range deltas {
				sum += d
				if d < minDelta {
					minDelta = d
				}
				if d > maxDelta {
					maxDelta = d
				}
			}
			avgDelta = sum / float64(len(deltas))
		}

		values = append(values, []interface{}{"", "", ""})
		values = append(values, []interface{}{"Статистика", "", ""})
		values = append(values, []interface{}{"Всего нажатий", len(timestamps), ""})
		values = append(values, []interface{}{"Минимальный интервал (сек)", minDelta, ""})
		values = append(values, []interface{}{"Максимальный интервал (сек)", maxDelta, ""})
		values = append(values, []interface{}{"Средний интервал между нажатиями (сек)", avgDelta, ""})
	}

	spreadsheetID := "xxxxxxxx"
	rangeData := "Лист1!A1"

	_, err = srv.Spreadsheets.Values.Clear(spreadsheetID, rangeData, &sheets.ClearValuesRequest{}).Do()
	if err != nil {
		return fmt.Errorf("не удалось очистить таблицу: %v", err)
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, rangeData, valueRange).
		ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("не удалось записать данные в таблицу: %v", err)
	}

	fmt.Println("Данные успешно импортированы в Google Таблицы!")
	return nil
}
