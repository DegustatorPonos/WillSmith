package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const VersionName string = "0.3.3a"
const HomePage string = "gemini://geminiprotocol.net/"
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
			height = newheight
			width = newwidth
			currentPage = ReadRequest(resp)
		}
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
		}

		if strings.HasPrefix(TrimmedCommand, "gemini://") {
			fmt.Print("Navigating to a specified page...")
			currntIndex, history = DirectToANewPage(TrimmedCommand, history, currntIndex, currentPage)
			continue
		}

		if slices.Contains(currentPage.Links, TrimmedCommand) {
			fmt.Print("Navigating to a next page...")
			var newURI = AppendToLink(history[currntIndex], TrimmedCommand)
			currntIndex, history = DirectToANewPage(newURI, history, currntIndex, currentPage)
			continue
		}

	}
}

func AppendToLink(BaseURI string, ToAppend string) string { 
	if(len(ToAppend) > 0 && ToAppend[0] == '/') {
		return strings.Join([]string{StripToHost(BaseURI), strings.Replace(ToAppend, "/", "", 1), "/"}, "")
	}
	var newURI = ""
	if IsAnEndpoint(BaseURI) {
		BaseURI = GoBackOneLayer(BaseURI)
	}
	if strings.HasSuffix(BaseURI, "/") || strings.HasPrefix(ToAppend, "/") {
		newURI = strings.Join([]string{BaseURI, ToAppend, "/"}, "")
	} else {
		newURI = strings.Join([]string{BaseURI, "/", ToAppend, "/"}, "")
	}
	return CompactAllBackwardsMotions(newURI)
}

func StripToHost(URI string) string {
	// This regex returns the mase address. For example, 
	// gemini://gemini.circumlunar.space/capcom/ returns gemini://gemini.circumlunar.space/
	var r = regexp.MustCompile(`gemini:\/\/[^\/:]*\/`)
	return r.FindString(URI)
}

func IsAnEndpoint(inp string) bool {
	var r = regexp.MustCompile(`\/[^\/:]*\.gmi\/?$`)
	return r.FindString(inp) != ""
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

func ClearConsole() {
	var cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func GetStatusBar(ScreenWidth int, ScreenHeight int, URI string, HistoryLength int, ScrollOffset int, PageLength int) string {
	var sb = strings.Builder{}
	sb.WriteString(URI)
	sb.WriteString(" | ")
	// Page position
	sb.WriteString("Position: ")
	sb.WriteString(strconv.Itoa(ScrollOffset))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(ScrollOffset - 5 + ScreenHeight))
	sb.WriteString("/")
	sb.WriteString(strconv.Itoa(PageLength))
	sb.WriteString(" | ")
	sb.WriteString("History: ")
	sb.WriteString(strconv.Itoa(HistoryLength))
	sb.WriteString(" | ")
	sb.WriteString("Window size: ")
	sb.WriteString(strconv.Itoa(ScreenWidth))
	sb.WriteString(" x ")
	sb.WriteString(strconv.Itoa(ScreenHeight))
	sb.WriteString(" | ")
	sb.WriteString("WillSmith v.")
	sb.WriteString(VersionName)
	sb.WriteString(" | ")
	if(sb.Len() >= ScreenWidth) {
		return sb.String()[0:ScreenWidth-1]
	}
	return sb.String()
}

func WriteLine(Width int) {
	for range Width {
		fmt.Print("=")
	}
	fmt.Println()
}
