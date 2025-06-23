package logger

import (
	globalstate "WillSmith/GlobalState"
	"fmt"
	"net/http"
	"strings"
)

const DebugToolPort int = 6969

const InfoMessage string    = "Info"
const WarningMessage string = "Warn"
const ErrorMessage string   = "Err"

var logChan = make(chan LogMessage, globalstate.State.ChannelLengths.LogChannel)

type LogMessage struct {
	Message string
	Type string
}

func SendInfo(inp string) {
	if globalstate.CurrentSettings.EnableLogging {
		logChan <- LogMessage{Type: InfoMessage, Message: inp}
	}
}

func SendWarning(inp string) {
	if globalstate.CurrentSettings.EnableLogging {
		logChan <- LogMessage{Type: WarningMessage, Message: inp}
	}
}

func SendError(inp string) {
	if globalstate.CurrentSettings.EnableLogging {
		logChan <- LogMessage{Type: ErrorMessage, Message: inp}
	}
}

func CreateLoggingTask() {
	if globalstate.CurrentSettings.EnableLogging {
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
