package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
func (tab *Tab) ScrollDownUntilTheClosestHeader() { 
	if tab.currentPosition >= len(tab.currentPage.Text) {
		tab.currentPosition = len(tab.currentPage.Text) - 1
	}
	for {
		tab.currentPosition -= 1
		if tab.currentPosition <= 0 || strings.HasPrefix(tab.currentPage.Text[tab.currentPosition], "#"){
			if tab.currentPosition < 0 {
				tab.currentPosition = 0
			}
			return
		}
	}
}

// Up until the closest string starting with '#'
func (tab *Tab) ScrollUpUntilTheClosestHeader() { 
	if tab.currentPosition >= len(tab.currentPage.Text) {
		tab.currentPosition = len(tab.currentPage.Text) - 1
	}
	var MaxPosition = len(tab.currentPage.Text)
	for {
		tab.currentPosition += 1
		if tab.currentPosition == MaxPosition || strings.HasPrefix(tab.currentPage.Text[tab.currentPosition], "#") {
			return
		}
	}
}

func (tab *Tab) ScrollDownUntilTheClosestSpace() {
	if tab.currentPosition >= len(tab.currentPage.Text) {
		tab.currentPosition = len(tab.currentPage.Text) - 1
	}
	var MaxPosition = len(tab.currentPage.Text)
	for {
		tab.currentPosition += 1
		if len(strings.Trim(tab.currentPage.Text[tab.currentPosition], " ")) < 2 || tab.currentPosition == MaxPosition {
			tab.currentPosition += 1
			return
		}
	}
}

func (tab *Tab) ScrollUpUntilTheClosestSpace() {
	if tab.currentPosition >= len(tab.currentPage.Text) {
		tab.currentPosition = len(tab.currentPage.Text) - 1
	}
	tab.currentPosition -= 1 // Compenstaing to additional +1 on every return
	for {
		tab.currentPosition -= 1
		if tab.currentPosition < 0 {
			tab.currentPosition = 0
			return
		}
		if len(strings.Trim(tab.currentPage.Text[tab.currentPosition], " ")) < 2 {
			tab.currentPosition += 1
			return
		}
	}
}

// Returns true if the app should still run
func HandleCommand(command string, currentTab *Tab, requestChan chan RequestCommand, TerminationChan chan bool) (bool) {
	switch command {
		case "": // Rerendering the page
			return true
		case ":q": // Quitting the app
			return false
		case "..": // Going to the previous page
			currentTab.PopPage(requestChan)
			return true
		case ":r": // Reload current page
		requestChan <- RequestCommand{URL: currentTab.history[currentTab.historyLength - 1], MandatoryReload: true}
			return true
		case ":u": // Abort all current responces
			TerminationChan <- true
			currentTab.PendingRequests = 0
			return true

		// Movement
		case "/": // Scroll up by half a page
			currentTab.currentPosition += currentTab.screenInfo.Height / 2
			return true
		case "\\": // Scroll down by half a page
			currentTab.currentPosition -= currentTab.screenInfo.Height / 2
			if currentTab.currentPosition < 0 {
				currentTab.currentPosition = 0
			}
			return true
		case "}": // Scroll down until the closest space
			currentTab.ScrollDownUntilTheClosestSpace()
			return true
		case "{": // Scroll up until the closest space
			currentTab.ScrollUpUntilTheClosestSpace()
			return true
		case "[": // Scroll down until the closest header
			currentTab.ScrollDownUntilTheClosestHeader()
			return true
		case "]": // Scroll up until the closest header 
			currentTab.ScrollUpUntilTheClosestHeader()
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
		requestChan <- RequestCommand{ URL: command }
		currentTab.PendingRequests += 1
		return true
	}

	// Going to a page by relative link
	if slices.Contains(currentTab.currentPage.Links, command) {
		fmt.Print("Navigating to a next page...")
		var newURI = AppendToLink(currentTab.currentPage.URI, command)
		requestChan <- RequestCommand{ URL: newURI  }
		currentTab.PendingRequests += 1
		return true
	}

	return true
}
