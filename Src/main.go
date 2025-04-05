package main

// The main app flow

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const VersionName string = "0.4a"
const HomePageFile string = "file://../StaticPages/IndexPage"

const CTRL_CH_LEN int = 2
const REQ_CH_LEN int = 2

type Tab struct {
	history []string
	historyLength int
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

func (tab *Tab) PopPage(requestChannel chan string) {
	if tab.historyLength > 1 {
		tab.historyLength -= 2
	} else {
		tab.historyLength -= 1
	}
	requestChannel <- tab.history[tab.historyLength]
}

func main() {
	var CurrentTab = Tab{
		history: make([]string, 0),
		historyLength: 0,
	}
	// Initial size

	// CHANNELING
	var ControlChan = make(chan int, CTRL_CH_LEN)
	var RequestChan = make(chan string, REQ_CH_LEN)
	var TerminationChan = make(chan bool, REQ_CH_LEN)
	
	// STARTING COROUTINES
	var CommandsChannel = CreateCommandChannel(&ControlChan)
	var ResponceChannel = CreateConnectionTask(&RequestChan, &TerminationChan, &ControlChan)
	var ScreenInfoChannel = GetScreenChannel(&ControlChan)

	// Getting a start page
	RequestChan <- HomePageFile
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
			CurrentTab.currentPage = *ReadRequest(responce)
			RenderPage(&CurrentTab)
			continue
		case SCR_CHN_ID:
			var NewSize = <- ScreenInfoChannel
			CurrentTab.screenInfo = NewSize
			RenderPage(&CurrentTab)
			continue
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
	//		case "}": // Go until the next white space
	//			if(int(currentPage.ScrollOffser) >= len(currentPage.Text)) {
	//				currentPage.ScrollOffser = uint(len(currentPage.Text) - 1)
	//			}
	//			currentPage.ScrollOffser += 1;
	//			for(int(currentPage.ScrollOffser) < len(currentPage.Text) && currentPage.Text[currentPage.ScrollOffser] != "") {
	//				currentPage.ScrollOffser += 1;
	//			}
	//			currentPage.ScrollOffser += 1;
	//			continue
	//		case "{": // Go until the pervious white space
	//			if(currentPage.ScrollOffser < 2) {
	//				continue
	//			}
	//			if(int(currentPage.ScrollOffser) >= len(currentPage.Text)) {
	//				currentPage.ScrollOffser = uint(len(currentPage.Text) - 1)
	//			}
	//			currentPage.ScrollOffser -= 2;
	//			for(int(currentPage.ScrollOffser) > 0 && currentPage.Text[currentPage.ScrollOffser] != "") {
	//				currentPage.ScrollOffser -= 1;
	//			}
	//			continue
}

func RenderPage(currentTab *Tab) {
	ClearConsole()
	fmt.Println(GetStatusBar(currentTab))
	WriteLine(currentTab.screenInfo.Width)
	currentTab.currentPage.ScrollOffser = uint(currentTab.currentPosition)
	DisplayPage(&currentTab.currentPage)
	WriteLine(currentTab.screenInfo.Width)
	fmt.Print("Enter command: >")
}

// Returns true if the app should still run
func HandleCommand(command string, currentTab *Tab, requestChan chan string, TerminationChan chan bool) (bool) {
	switch command {
		case "": // Rerendering the page
			return true
		case ":q": // Quitting the app
			return false
		case "/": // Scroll up by half a page
			currentTab.currentPosition += currentTab.screenInfo.Height / 2
			return true
		case ":r": // Scroll up by half a page
			requestChan <- currentTab.history[currentTab.historyLength - 1]
			return true
		case ":u": // Scroll up by half a page
			TerminationChan <- true
			currentTab.PendingRequests = 0
			return true
		case "\\": // Scroll down by half a page
			currentTab.currentPosition -= currentTab.screenInfo.Height / 2
			if currentTab.currentPosition < 0 {
				currentTab.currentPosition = 0
			}
			return true
		case "..":
			currentTab.PopPage(requestChan)
			return true
	}

	// Going to a link by its index
	if strings.HasPrefix(command, ":") {
		var LinkIndex, err = strconv.Atoi(strings.ReplaceAll(command, ":", ""))
		if err != nil || LinkIndex >= len(currentTab.currentPage.Links) {
			return true
		}
		command = currentTab.currentPage.Links[LinkIndex]
	}

	// Going to a page by full link
	if strings.HasPrefix(command, "gemini://") || strings.HasPrefix(command, "file://") {
		fmt.Print("Navigating to a specified page...")
		if !strings.HasSuffix(command, "/") && !IsAnEndpoint(command) {
			command = strings.Join([]string{command,"/"}, "")
		}
		requestChan <- command
		currentTab.PendingRequests += 1
		return true
	}

	// Going to a page by relative link
	if slices.Contains(currentTab.currentPage.Links, command) {
		fmt.Print("Navigating to a next page...")
		var newURI = AppendToLink(currentTab.currentPage.URI, command)
		requestChan <- newURI 
		currentTab.PendingRequests += 1
		return true
	}

	return true
}

// Returns new index and history
func DirectToANewPage(NewPageURI string, history []string, currntIndex int, currentPage *Page) (int, []string) {
	if len(history) <= currntIndex + 1 {
		history = append(history, NewPageURI)
	} else {
		history[currntIndex + 1] = NewPageURI
	}
	currntIndex += 1
	var resp = SendRequest(history[currntIndex], DEFAULT_PORT)
	currentPage = ReadRequest(resp)
	return currntIndex, history
}

