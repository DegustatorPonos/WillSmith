package main

// The main app flow

const VersionName string = "0.4.1a"
const HomePageFile string = "file://../StaticPages/IndexPage"

const CTRL_CH_LEN int = 2
const REQ_CH_LEN int = 2

type Tab struct {
	history []string
	historyLength int
	currentResp Request
	currentPage Page
	screenInfo ScreenInfo
	currentPosition int
	PendingRequests int
}

func (tab *Tab) AddPage(newPage string) {
	if tab.historyLength > 0 && tab.history[tab.historyLength - 1] == newPage {
		return
	}
	if(len(tab.history) <= tab.historyLength) {
		tab.history = append(tab.history, newPage)
		tab.historyLength += 1
		return
	}
	tab.history[tab.historyLength] = newPage
	tab.historyLength += 1
}

func (tab *Tab) PopPage(requestChannel chan RequestCommand) {
	if tab.historyLength > 1 {
		tab.historyLength -= 2
	} else {
		tab.historyLength -= 1
	}
	requestChannel <- RequestCommand{ URL: tab.history[tab.historyLength] }
}

func main() {
	var CurrentTab = Tab{
		history: make([]string, 0),
		historyLength: 0,
	}
	// Initial size

	// CHANNELING
	var ControlChan = make(chan int, CTRL_CH_LEN)
	var RequestChan = make(chan RequestCommand, REQ_CH_LEN)
	var TerminationChan = make(chan bool, REQ_CH_LEN)
	
	// STARTING COROUTINES
	var CommandsChannel = CreateCommandChannel(&ControlChan)
	var ResponceChannel = CreateConnectionTask(&RequestChan, &TerminationChan, &ControlChan)
	var ScreenInfoChannel = GetScreenChannel(&ControlChan)

	// Getting a start page
	RequestChan <- RequestCommand{URL: HomePageFile}
	// Handling events
	for {
		CommandType := <- ControlChan
		switch CommandType{

		case CON_CHAN_ID:
			var responce = <- *ResponceChannel
			CurrentTab.PendingRequests -= 1
			if CurrentTab.PendingRequests < 0 {
				CurrentTab.PendingRequests = 0
			}
			CurrentTab.AddPage(responce.URI)
			CurrentTab.currentResp = *responce
			CurrentTab.currentPage = *ParseRequest(responce, CurrentTab.screenInfo)
			CurrentTab.currentPosition = 0
			RenderPage(&CurrentTab)
			continue

		case SCR_CHN_ID:
			var NewSize = <- ScreenInfoChannel
			CurrentTab.screenInfo = NewSize
			CurrentTab.currentPage = *ParseRequest(&CurrentTab.currentResp, CurrentTab.screenInfo)
			RenderPage(&CurrentTab)

		case CMD_CHAN_ID:
			var command = <- CommandsChannel
			if !(HandleCommand(command, &CurrentTab, RequestChan, TerminationChan)) {
				ClearConsole()
				return
			}
			RenderPage(&CurrentTab)
			continue 

		}
	}
}
