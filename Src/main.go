package main

import (
	geminiprotocol "WillSmith/GeminiProtocol"
	globalstate "WillSmith/GlobalState"
	localresources "WillSmith/LocalResources"
	logger "WillSmith/Logger"
	tuihandlers "WillSmith/TUIHandlers"
)

// The main app flow

const HomePageFile string = "file://../StaticPages/IndexPage"

func main() {
	globalstate.ReadSettings()

	var CurrentTab = tuihandlers.Tab {
		History: make([]string, 0),
		HistoryLength: 0,
	}
	// Initial size

	// CHANNELING
	var ControlChan = make(chan int, globalstate.State.ChannelLengths.ControlChannel)
	var RequestChan = make(chan geminiprotocol.RequestCommand, globalstate.State.ChannelLengths.RequestChannel)
	var TerminationChan = make(chan bool, globalstate.State.ChannelLengths.RequestChannel)
	
	// STARTING COROUTINES
	var CommandsChannel = tuihandlers.CreateCommandChannel(&ControlChan)
	var ResponceChannel, DownloadChannel = geminiprotocol.CreateConnectionTask(&RequestChan, &TerminationChan, &ControlChan)
	var ScreenInfoChannel = tuihandlers.GetScreenChannel(&ControlChan)
	logger.CreateLoggingTask()
	geminiprotocol.InitCache()

	// Getting a start page
	RequestChan <- geminiprotocol.RequestCommand{URL: HomePageFile}

	// Handling events
	for {
		CommandType := <- ControlChan
		switch CommandType{

		case geminiprotocol.CON_CHAN_ID:
			var responce = <- *ResponceChannel
			CurrentTab.PendingRequests -= 1
			if CurrentTab.PendingRequests < 0 {
				CurrentTab.PendingRequests = 0
			}
			CurrentTab.AddPage(responce.URI)
			CurrentTab.CurrentResp = *responce
			CurrentTab.CurrentPage = *tuihandlers.ParseRequest(responce, CurrentTab.ScreenInfo)
			CurrentTab.CurrentPosition = 0
			tuihandlers.RenderPage(&CurrentTab)
			continue

		case tuihandlers.SCR_CHN_ID:
			var NewSize = <- ScreenInfoChannel
			CurrentTab.ScreenInfo = NewSize
			CurrentTab.CurrentPage = *tuihandlers.ParseRequest(&CurrentTab.CurrentResp, CurrentTab.ScreenInfo)
			tuihandlers.RenderPage(&CurrentTab)

		case tuihandlers.CMD_CHAN_ID:
			var command = <- CommandsChannel
			if !(tuihandlers.HandleCommand(command, &CurrentTab, RequestChan, TerminationChan)) {
				tuihandlers.ClearConsole()
				return
			}
			tuihandlers.RenderPage(&CurrentTab)
			continue 

		case geminiprotocol.DOWNLOAD_CHAN_ID :
			var resp = <- *DownloadChannel
			localresources.Download(resp.URI, resp.Body)

		}
	}
}
