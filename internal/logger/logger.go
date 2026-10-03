package logger

import (
	"TaskForge/internal/task"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func LogSuccess(curTask task.Task) {
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer f.Close()

	payloadStr, err := json.Marshal(curTask.Payload)
	if err != nil {
		payloadStr = []byte{}
	}

	text := fmt.Sprintf("\nSUCCESS: Task type: %s | Payload: %s | Retries left: %d",
		curTask.Type, string(payloadStr), curTask.Retries)

	if _, err := f.WriteString(text); err != nil {
		log.Println("Error writing to log file:", err)
		return
	}
}

func LogFailure(curTask task.Task, curErr error) {
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer f.Close()

	payloadStr, err := json.Marshal(curTask.Payload)
	if err != nil {
		payloadStr = []byte{}
	}

	text := fmt.Sprintf("\nFAILURE: Task type: %s | Payload: %s | Retries left: %d | Error: %s",
		curTask.Type, string(payloadStr), curTask.Retries, curErr.Error())

	if _, err := f.WriteString(text); err != nil {
		log.Println("Error writing to log file:", err)
		return
	}
}