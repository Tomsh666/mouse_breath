package main

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func startMeasuring(entry *widget.Entry) {
	secondsStr := entry.Text
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil || seconds <= 0 {
		fmt.Println("Error: incorrect number")
		return
	}
	fmt.Println("Seconds:", seconds)
	// TODO: Добавить логику для старта измерения, для отправки seconds на контроллер
	fmt.Println("measuring")
}
