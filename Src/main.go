package main

// The main app flow

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"golang.org/x/term"
)

const VersionName string = "0.3.5a"
const HomePageFile string = "file://../StaticPages/IndexPage"

func main() {
	var history = make([]string, 1)
	history[0] = HomePageFile
	var currntIndex = 0
	var PrevCurrentIndex = 0
	var reader = bufio.NewReader(os.Stdin)

	var resp = SendRequest(history[currntIndex], DEFAULT_PORT)
	var currentPage = ReadRequest(resp)
	var width, height, _ = term.GetSize(0)

	for {
		if PrevCurrentIndex != currntIndex {
			resp = SendRequest(history[currntIndex], DEFAULT_PORT)
			currentPage = ReadRequest(resp)
			PrevCurrentIndex = currntIndex
		}
		var newwidth, newheight, _ = term.GetSize(0)
		if(newheight != height || newwidth != width) {
			// Resizing
			currentPage = ReadRequest(resp)
			height = newheight
			width = newwidth
		}

		// Rendering the screen and reading the command
		ClearConsole()
		fmt.Println(GetStatusBar(width, height, history[currntIndex], currntIndex, int(currentPage.ScrollOffser), len(currentPage.Text)))
		WriteLine(width)
		DisplayPage(currentPage)
		WriteLine(width)
		fmt.Print("Enter command: >")
		var command, _ = reader.ReadString('\n')
		var TrimmedCommand = strings.TrimRight(command, "\n")
		TrimmedCommand = strings.Trim(TrimmedCommand, " ")

		// Handling commands
		switch(TrimmedCommand) {
		case "": // Rerendering the page
			continue
		case "..": // Go to previous page
			if(currntIndex >= 1) {
				currntIndex -= 1
			}
			continue
		case "/": // Scroll up by half a page
			currentPage.ScrollOffser += uint(height / 2)
			continue
		case "\\": // Scroll up by half a page
			if(currentPage.ScrollOffser > uint(height / 2)) {
				currentPage.ScrollOffser -= uint(height / 2)
			} else {
				currentPage.ScrollOffser = 0
			}
			continue
		case "}": // Go until the next white space
			if(int(currentPage.ScrollOffser) >= len(currentPage.Text)) {
				currentPage.ScrollOffser = uint(len(currentPage.Text) - 1)
			}
			currentPage.ScrollOffser += 1;
			for(int(currentPage.ScrollOffser) < len(currentPage.Text) && currentPage.Text[currentPage.ScrollOffser] != "") {
				currentPage.ScrollOffser += 1;
			}
			currentPage.ScrollOffser += 1;
			continue
		case "{": // Go until the pervious white space
			if(currentPage.ScrollOffser < 2) {
				continue
			}
			if(int(currentPage.ScrollOffser) >= len(currentPage.Text)) {
				currentPage.ScrollOffser = uint(len(currentPage.Text) - 1)
			}
			currentPage.ScrollOffser -= 2;
			for(int(currentPage.ScrollOffser) > 0 && currentPage.Text[currentPage.ScrollOffser] != "") {
				currentPage.ScrollOffser -= 1;
			}
			continue
		case ":q": // Exit the app
			ClearConsole()
			return
		case ":r": // Reload the current page
			resp = SendRequest(history[currntIndex], DEFAULT_PORT)
			currentPage = ReadRequest(resp)
			continue
		}

		// Going to a page by full link
		if strings.HasPrefix(TrimmedCommand, "gemini://") {
			fmt.Print("Navigating to a specified page...")
			currntIndex, history = DirectToANewPage(TrimmedCommand, history, currntIndex, currentPage)
			continue
		}

		// Going to a page by relative link
		if slices.Contains(currentPage.Links, TrimmedCommand) {
			fmt.Print("Navigating to a next page...")
			var newURI = AppendToLink(history[currntIndex], TrimmedCommand)
			currntIndex, history = DirectToANewPage(newURI, history, currntIndex, currentPage)
			continue
		}

	}
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
