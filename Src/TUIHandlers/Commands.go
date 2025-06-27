package tuihandlers

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	geminiprotocol "WillSmith/GeminiProtocol"
	localresources "WillSmith/LocalResources"
	logger "WillSmith/Logger"
)

const CMD_CHAN_BUFF_SIZE int = 1
const CMD_CHAN_ID int = 1

func NavigationTask(output chan string, controlChannel *chan int) {
	var reader = bufio.NewReader(os.Stdin)
	for {
		var command, readErr = reader.ReadString('\n')
		if readErr != nil {
			panic("Error in the command reading coroutine.")
		}
		command = strings.Trim(command, "\n")
		// fmt.Printf("Sending a command \"%v\"\n", command)
		output <- command
		*controlChannel <- CMD_CHAN_ID
	}
}

func CreateCommandChannel(controlChannel *chan int) chan string {
	var outpChannel = make(chan string, CMD_CHAN_BUFF_SIZE);
	go NavigationTask(outpChannel, controlChannel)
	return outpChannel
}

// Down until the closest string starting with '#'
func (tab * Tab) ScrollDownUntilTheClosestHeader() { 
	if tab.CurrentPosition >= len(tab.CurrentPage.Text) {
		tab.CurrentPosition = len(tab.CurrentPage.Text) - 1
	}
	for {
		tab.CurrentPosition -= 1
		if tab.CurrentPosition <= 0 || strings.HasPrefix(tab.CurrentPage.Text[tab.CurrentPosition], "#"){
			if tab.CurrentPosition < 0 {
				tab.CurrentPosition = 0
			}
			return
		}
	}
}

// Up until the closest string starting with '#'
func (tab *Tab) ScrollUpUntilTheClosestHeader() { 
	if tab.CurrentPosition >= len(tab.CurrentPage.Text) {
		tab.CurrentPosition = len(tab.CurrentPage.Text) - 1
	}
	var MaxPosition = len(tab.CurrentPage.Text)
	for {
		tab.CurrentPosition += 1
		if tab.CurrentPosition == MaxPosition || strings.HasPrefix(tab.CurrentPage.Text[tab.CurrentPosition], "#") {
			return
		}
	}
}

func (tab *Tab) ScrollDownUntilTheClosestSpace() {
	if tab.CurrentPosition >= len(tab.CurrentPage.Text) {
		tab.CurrentPosition = len(tab.CurrentPage.Text) - 1
	}
	var MaxPosition = len(tab.CurrentPage.Text)
	for {
		tab.CurrentPosition += 1
		if len(strings.Trim(tab.CurrentPage.Text[tab.CurrentPosition], " ")) < 2 || tab.CurrentPosition == MaxPosition {
			tab.CurrentPosition += 1
			return
		}
	}
}

func (tab *Tab) ScrollUpUntilTheClosestSpace() {
	if tab.CurrentPosition >= len(tab.CurrentPage.Text) {
		tab.CurrentPosition = len(tab.CurrentPage.Text) - 1
	}
	tab.CurrentPosition -= 1 // Compenstaing to additional +1 on every return
	for {
		tab.CurrentPosition -= 1
		if tab.CurrentPosition < 0 {
			tab.CurrentPosition = 0
			return
		}
		if len(strings.Trim(tab.CurrentPage.Text[tab.CurrentPosition], " ")) < 2 {
			tab.CurrentPosition += 1
			return
		}
	}
}

// Returns true if the app should still run
func HandleCommand(command string, CurrentTab *Tab, requestChan chan geminiprotocol.RequestCommand, TerminationChan chan bool) (bool) {
	switch command {
		case "": // Rerendering the page
			return true
		case ":q": // Quitting the app
			logger.SendInfo("=========== END OF SESSION ===========")
			logger.SendInfo("")
			return false
		case "..": // Going to the previous page
			CurrentTab.PopPage(requestChan)
			return true
		case ":r": // Reload Current page
		requestChan <- geminiprotocol.RequestCommand{URL: CurrentTab.History[CurrentTab.HistoryLength - 1], MandatoryReload: true}
			return true
		case ":u": // Abort all Current responces
			TerminationChan <- true
			CurrentTab.PendingRequests = 0
			return true

		// Movement
		case "/": // Scroll up by half a page
			CurrentTab.CurrentPosition += CurrentTab.ScreenInfo.Height / 2
			return true
		case "\\": // Scroll down by half a page
			CurrentTab.CurrentPosition -= CurrentTab.ScreenInfo.Height / 2
			if CurrentTab.CurrentPosition < 0 {
				CurrentTab.CurrentPosition = 0
			}
			return true
		case "}": // Scroll down until the closest space
			CurrentTab.ScrollDownUntilTheClosestSpace()
			return true
		case "{": // Scroll up until the closest space
			CurrentTab.ScrollUpUntilTheClosestSpace()
			return true
		case "[": // Scroll down until the closest header
			CurrentTab.ScrollDownUntilTheClosestHeader()
			return true
		case "]": // Scroll up until the closest header 
			CurrentTab.ScrollUpUntilTheClosestHeader()
		CurrentTab.ScrollUpUntilTheClosestHeader()
		return true
	}

	// Bookmarks
	if strings.HasPrefix(command, ":b ") {
		var args = strings.Split(command, " ")
		if len(args) < 2 {
			return true
		}
		var description = args[1]
		localresources.AddBookmark(localresources.Bookmark{
			URL: CurrentTab.CurrentPage.URI, 
			Description: description,
		})
		return true
	}

	// Delete command
	if strings.HasPrefix(command, ":delb") {
		var args = strings.Split(command, " ")
		if len(args) < 2 {
			localresources.DeleteBookmark(CurrentTab.CurrentPage.URI)
			return true
		}
		var URLarg = args[1]
		localresources.DeleteBookmark(getURIFromArgument(URLarg, CurrentTab))
		return true
	}

	// Downloading
	if strings.HasPrefix(command, ":d") {
		var args = strings.Split(command, " ")
		if len(args) < 2 {
			requestChan <- geminiprotocol.RequestCommand{ URL: CurrentTab.CurrentPage.URI, TargetAction: geminiprotocol.DOWNLOAD }
			return true
		}
		var URLarg = args[1]
		var targetURI = getURIFromArgument(URLarg, CurrentTab)
		requestChan <- geminiprotocol.RequestCommand{ URL: targetURI, TargetAction: geminiprotocol.DOWNLOAD }
		return true
	}

	// Going to a link by its index
	if strings.HasPrefix(command, ":") {
		var LinkIndex, err = strconv.Atoi(strings.ReplaceAll(command, ":", ""))
		if err != nil || LinkIndex >= len(CurrentTab.CurrentPage.Links) {
			return true
		}
		command = CurrentTab.CurrentPage.Links[LinkIndex]
	}

	// Going to a page by full link
	if strings.HasPrefix(command, "gemini://") || strings.HasPrefix(command, "file://") {
		fmt.Print("Navigating to a specified page...")
		if !strings.HasSuffix(command, "/") && !geminiprotocol.IsAnEndpoint(command) {
			command = strings.Join([]string{command,"/"}, "")
		}
		requestChan <- geminiprotocol.RequestCommand{ URL: command }
		CurrentTab.PendingRequests += 1
		return true
	}

	// Going to a page by relative link
	if slices.Contains(CurrentTab.CurrentPage.Links, command) {
		fmt.Print("Navigating to a next page...")
		var newURI = geminiprotocol.AppendToLink(CurrentTab.CurrentPage.URI, command)
		requestChan <- geminiprotocol.RequestCommand{ URL: newURI  }
		CurrentTab.PendingRequests += 1
		return true
	}

	return true
}

func getURIFromArgument(arg string, CurrentTab *Tab) string {
	var trg = arg
	var LinkIndex, err = strconv.Atoi(arg)
	if err == nil && LinkIndex < len(CurrentTab.CurrentPage.Links) {
		trg = CurrentTab.CurrentPage.Links[LinkIndex]
	}
	if strings.HasPrefix(trg, "gemini") || strings.HasPrefix(trg, "file") {
		if !strings.HasSuffix(trg, "/") && !geminiprotocol.IsAnEndpoint(trg) {
			trg = strings.Join([]string{trg ,"/"}, "")
		}
		return trg
	}
	if slices.Contains(CurrentTab.CurrentPage.Links, trg) {
		var newURI = geminiprotocol.AppendToLink(CurrentTab.CurrentPage.URI, trg)
		return newURI
	}
	return trg
}
