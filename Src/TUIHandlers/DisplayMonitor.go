package tuihandlers

import (
	"time"

	"golang.org/x/term"
)

type ScreenInfo struct {
	Height int
	Width int
}

const SCR_CHN_ID int = 3
const SCR_CHN_BUF_LEN int = 1
// Refresh rate in ms
const SCR_ROUTINE_REFRESH_DELAY time.Duration = 250

func ScreenMonitorRoutine(outp chan ScreenInfo) {
	var buffer = ScreenInfo{}
	for {
		var newwidth, newheight, _ = term.GetSize(0)
		if newheight != buffer.Height || newwidth != buffer.Width {
			buffer.Height = newheight
			buffer.Width = newwidth
			outp <- buffer
		}
		time.Sleep(SCR_ROUTINE_REFRESH_DELAY * time.Millisecond)
	}
}

func GetScreenChannel() chan ScreenInfo {
	var outp = make(chan ScreenInfo, SCR_CHN_BUF_LEN)
	go ScreenMonitorRoutine(outp)
	return outp
}
