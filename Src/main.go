package main

import (
	geminiprotocol "WillSmith/GeminiProtocol"
	globalstate "WillSmith/GlobalState"
	localresources "WillSmith/LocalResources"
	logger "WillSmith/Logger"
	tuihandlers "WillSmith/TUIHandlers"
	"fmt"
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
	var RequestChan = make(chan geminiprotocol.RequestCommand, globalstate.State.ChannelLengths.RequestChannel)
	var TerminationChan = make(chan bool, globalstate.State.ChannelLengths.RequestChannel)
	
	// STARTING COROUTINES
	// var CommandsChannel = tuihandlers.CreateCommandChannel(&ControlChan)
	var CommandsChannel, EchoChannel = tuihandlers.CreateInputHandler()
	var ResponceChannel, DownloadChannel = geminiprotocol.CreateConnectionTask(&RequestChan, &TerminationChan)
	var ScreenInfoChannel = tuihandlers.GetScreenChannel()
	logger.CreateLoggingTask()
	geminiprotocol.InitCache()

	// Getting a start page
	RequestChan <- geminiprotocol.RequestCommand{URL: HomePageFile}

	// Handling events
	for {
		// logger.SendInfo("Listening")
		select {
		case responce := <- *ResponceChannel:
			// logger.SendInfo("Drawing")
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

		case NewSize := <- ScreenInfoChannel:
			// logger.SendInfo("Rescaling")
			CurrentTab.ScreenInfo = NewSize
			CurrentTab.CurrentPage = *tuihandlers.ParseRequest(&CurrentTab.CurrentResp, CurrentTab.ScreenInfo)
			tuihandlers.RenderPage(&CurrentTab)
			continue

		case command := <- CommandsChannel:
			logger.SendInfo(fmt.Sprintf("Comand: %s", command))
			if !(tuihandlers.HandleCommand(command, &CurrentTab, RequestChan, TerminationChan)) {
				tuihandlers.ClearConsole()
				return
			}
			tuihandlers.RenderPage(&CurrentTab)
			continue

		case resp := <- *DownloadChannel:
			localresources.Download(resp.URI, resp.Body)
			continue

		case char := <- EchoChannel:
			logger.SendInfo("Echo")
			fmt.Printf(string(char))
			continue

		}
	}
}
