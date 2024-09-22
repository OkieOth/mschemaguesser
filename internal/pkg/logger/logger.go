package logger

import (
	"log"
	"os"
)

func Init(logFile string) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
		panic("Failed to open log file")
	}
	log.SetOutput(file)
}
