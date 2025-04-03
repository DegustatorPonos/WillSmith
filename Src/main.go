package main

// The main app flow

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const VersionName string = "0.3.7a"
const HomePageFile string = "file://../StaticPages/IndexPage"

const CTRL_CH_LEN int = 2
const REQ_CH_LEN int = 2

type Tab struct {
	history []string
	historyLength int
	currentPage Page
	screenInfo ScreenInfo
	currentPosition int
}

func (tab *Tab) AddPage(newPage string) {
	if(len(tab.history) >= tab.historyLength) {
		tab.history = append(tab.history, newPage)
		tab.historyLength += 1
		return
	}
	tab.history[tab.historyLength] = newPage
	tab.historyLength += 1

}

func main() {

	// TODO: Delete
	var history = make([]string, 1)
	history[0] = HomePageFile
	var currntIndex = 0
	var PrevCurrentIndex = 0
	var reader = bufio.NewReader(os.Stdin)
	//	var resp = SendRequest(history[currntIndex], DEFAULT_PORT)
	var width, height, _ = term.GetSize(0)
	// END OF TODO

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
		fmt.Printf("Command type: %v\n", CommandType)
		switch CommandType{
		case CON_CHAN_ID:
			var responce = <- *ResponceChannel
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

	for {
		if PrevCurrentIndex != currntIndex {
//			resp = SendRequest(history[currntIndex], DEFAULT_PORT)
			PrevCurrentIndex = currntIndex
		}
		var newwidth, newheight, _ = term.GetSize(0)
		if(newheight != height || newwidth != width) {
			// Resizing
			height = newheight
			width = newwidth
		}

		// Rendering the screen and reading the command
		ClearConsole()
		WriteLine(width)
		WriteLine(width)
		fmt.Print("Enter command: >")
		var command, _ = reader.ReadString('\n')
		var TrimmedCommand = strings.TrimRight(command, "\n")
		TrimmedCommand = strings.Trim(TrimmedCommand, " ")

		// Handling commands
		switch(TrimmedCommand) {
		case "..": // Go to previous page
			if(currntIndex >= 1) {
				currntIndex -= 1
			}
			continue
//		case "/": // Scroll up by half a page
//			currentPage.ScrollOffser += uint(height / 2)
//			continue
//		case "\\": // Scroll up by half a page
//			if(currentPage.ScrollOffser > uint(height / 2)) {
//				currentPage.ScrollOffser -= uint(height / 2)
//			} else {
//				currentPage.ScrollOffser = 0
//			}
//			continue
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
case ":r": // Reload the current page
//			resp = SendRequest(history[currntIndex], DEFAULT_PORT)
continue
		}

		// Going to a page by full link
		if strings.HasPrefix(TrimmedCommand, "gemini://") || strings.HasPrefix(TrimmedCommand, "file://") {
		fmt.Print("Navigating to a specified page...")
		if !strings.HasSuffix(TrimmedCommand, "/") && !IsAnEndpoint(TrimmedCommand) {
			TrimmedCommand = strings.Join([]string{TrimmedCommand ,"/"}, "")
		}
		//			currntIndex, history = DirectToANewPage(TrimmedCommand, history, currntIndex, currentPage)
		continue
	}

	// Going to a page by relative link
	//		if slices.Contains(currentPage.Links, TrimmedCommand) {
	//			fmt.Print("Navigating to a next page...")
	//			var newURI = AppendToLink(history[currntIndex], TrimmedCommand)
	//			currntIndex, history = DirectToANewPage(newURI, history, currntIndex, currentPage)
	//			continue
	//		}

}
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
		case "\\": // Scroll down by half a page
			currentTab.currentPosition -= currentTab.screenInfo.Height / 2
			if currentTab.currentPosition < 0 {
				currentTab.currentPosition = 0
			}
			return true
		case "..":
			if currentTab.historyLength <= 1 {
				return true
			}
			currentTab.historyLength -= 1
			requestChan <- currentTab.history[currentTab.historyLength - 1]
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
		return true
	}

	// Going to a page by relative link
	if slices.Contains(currentTab.currentPage.Links, command) {
		fmt.Print("Navigating to a next page...")
		var newURI = AppendToLink(currentTab.currentPage.URI, command)
		requestChan <- newURI 
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

