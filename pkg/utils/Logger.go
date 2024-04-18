package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

func WriteLogToFile(logMessage string) {
	startDate := time.Now().Format("2006-01-02")
	logFolderPath := "./log"
	logFilePath := fmt.Sprintf("%s/logFile-%s.log",
		logFolderPath, startDate)

	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		os.MkdirAll(logFolderPath, 0777)
	}

	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		os.Create(logFilePath)
	}
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println(logMessage)
}
