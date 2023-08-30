package logs

import (
	"fmt"
	"log"
	"os"
	"time"
)

func SetupLogging() {
	currentTime := time.Now()
	filename := fmt.Sprintf("logs/%s.log", currentTime.Format("2006-01-02"))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	log.SetOutput(file)
}

func AutoRotateLogs() {
	for {

		now := time.Now()
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		time.Sleep(time.Until(nextDay))

		SetupLogging()
	}
}
