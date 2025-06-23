package tuihandlers

import geminiprotocol "WillSmith/GeminiProtocol"

type Tab struct {
	History []string
	HistoryLength int
	CurrentResp geminiprotocol.Request
	CurrentPage Page
	ScreenInfo ScreenInfo
	CurrentPosition int
	PendingRequests int
}

func (tab *Tab) AddPage(newPage string) {
	if tab.HistoryLength > 0 && tab.History[tab.HistoryLength - 1] == newPage {
		return
	}
	if(len(tab.History) <= tab.HistoryLength) {
		tab.History = append(tab.History, newPage)
		tab.HistoryLength += 1
		return
	}
	tab.History[tab.HistoryLength] = newPage
	tab.HistoryLength += 1
}

func (tab *Tab) PopPage(requestChannel chan geminiprotocol.RequestCommand) {
	if tab.HistoryLength > 1 {
		tab.HistoryLength -= 2
	} else {
		tab.HistoryLength -= 1
	}
	requestChannel <- geminiprotocol.RequestCommand{ URL: tab.History[tab.HistoryLength] }
}

