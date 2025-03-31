package main

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"strconv"
)

var measurementData [][]interface{}

func startMeasuring(entry *widget.Entry) {
	secondsStr := entry.Text
	seconds, err := strconv.Atoi(secondsStr)
	if err != nil || seconds <= 0 {
		fmt.Println("Error: incorrect number")
		return
	}
	fmt.Println("Seconds:", seconds)

	measurementData = make([][]interface{}, 2)
	for i := range measurementData {
		measurementData[i] = make([]interface{}, seconds)
	}

	tmpFunc(seconds) //TODO: поменять заглушку, как-то передавать занчения с контроллера в массив
	fmt.Println("measuring")
}

func tmpFunc(seconds int) { //Заглушка
	for i := 0; i < seconds; i++ {
		measurementData[0][i] = i + 1
		measurementData[1][i] = rand.Intn(100)
	}
}
