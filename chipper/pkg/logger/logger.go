package logger

import (
	"fmt"
	"time"
)

var LogList string
var LogArray []string

func Println(a ...any) {
	fmt.Println(a...)
	LogArray = append(LogArray, time.Now().Format("2006.01.02 15:04:05")+": "+fmt.Sprint(a...)+"\n")
	if len(LogArray) >= 30 {
		LogArray = LogArray[1:]
	}
	LogList = ""
	for _, b := range LogArray {
		LogList = LogList + b
	}
}
