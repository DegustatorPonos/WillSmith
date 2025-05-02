package main

import (
	"fmt"
	"net/http"
	"strings"
)

const DebugToolPort int = 6969
const LOG_CH_LEN int = 2

const InfoMessage string    = "Info"
const WarningMessage string = "Warn"
const ErrorMessage string   = "Err"

const IsLoggingEnabled bool = true

var logChan = make(chan LogMessage, LOG_CH_LEN)

type LogMessage struct {
	Message string
	Type string
}

func SendInfo(inp string) {
	if IsLoggingEnabled {
		logChan <- LogMessage{Type: InfoMessage, Message: inp}
	}
}

func SendWarning(inp string) {
	if IsLoggingEnabled {
		logChan <- LogMessage{Type: WarningMessage, Message: inp}
	}
}

func SendError(inp string) {
	if IsLoggingEnabled {
		logChan <- LogMessage{Type: ErrorMessage, Message: inp}
	}
}

func CreateLoggingTask() {
	if IsLoggingEnabled {
		go loggingTask(logChan)
	}
}

func loggingTask(inpChan chan LogMessage) {
	var baseUrl = fmt.Sprintf("http://localhost:%d/", DebugToolPort)
	for {
		var msg = <- inpChan
		http.Post(strings.Join([]string{baseUrl, msg.Type}, ""), "", strings.NewReader(msg.Message))
	}
}
