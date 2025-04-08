package main

// Everything related to rendering

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type Page struct {
	URI string
	Text []string
	Links []string
	L1Headers []string
	ScrollOffser uint
}

func ParseRequest(r *Request, Screen ScreenInfo) *Page {
	var outp = Page{URI: r.URI}
	outp.Text = make([]string, 0)
	outp.Links = make([]string, 0)
	var width = Screen.Width
	for _, rawStr := range strings.Split(string(r.Body), "\n") {
		var line = rawStr
		if strings.HasPrefix(rawStr, "=>") {
			outp.Links = append(outp.Links, ParseLink(line))
			var Prefixless, _ = strings.CutPrefix(line, "=>")
			line = strings.Join([]string{"=> [", strconv.Itoa(len(outp.Links) - 1), "]", Prefixless}, "")
		}
		if len(line) < width {
			outp.Text = append(outp.Text, line)
			continue
		}
		var rightSide = 0
		for i :=0; i < len(line); i = rightSide {
			rightSide = i + width
			if rightSide >= len(line) {
				rightSide = len(line)
				outp.Text = append(outp.Text, strings.Trim(line[i:rightSide], " "))
				continue
			}
			for(rightSide > 0 && line[rightSide] != '\n' && line[rightSide] != ' ') {
				rightSide--
			}
			outp.Text = append(outp.Text, strings.Trim(line[i:rightSide], " "))
		}
	}
	return &outp
}

func ParseLink(inp string) string {
	var prefixless, _ = strings.CutPrefix(inp, "=>")
	var pureLink = strings.Trim(prefixless, " ")
	if strings.Contains(pureLink, "	") {
		pureLink = strings.Split(pureLink, "	")[0]
	}
	if strings.Contains(pureLink, " ") {
		pureLink = strings.Split(pureLink, " ")[0]
	}
	var outp, _ = strings.CutSuffix(pureLink, "/")
	switch outp {
	case ".." :
		return "../"
	case "/" :
		return "//"
	default:
		break
	}
	return outp
}

func DisplayPage(page *Page) {
	var _, height, _ = term.GetSize(0)
	height -= 5 // Subtract 2 lines, status line and command line
	for i := range height {
		if(uint(len(page.Text)) > uint(i) + page.ScrollOffser) {
			fmt.Println(page.Text[i + int(page.ScrollOffser)])
		} else {
			fmt.Println("")
		}
	}
}

func GetStatusBar(currentTab *Tab) string {
	var ScrollOffset = currentTab.currentPosition
	var sb = strings.Builder{}
	sb.WriteString(currentTab.currentPage.URI)
	sb.WriteString(" | ")
	// Page position
	sb.WriteString(fmt.Sprintf("Position: %v-%v/%v | ", ScrollOffset, ScrollOffset - 5 + currentTab.screenInfo.Height, len(currentTab.currentPage.Text)))
	sb.WriteString(fmt.Sprintf("History: %v | ", currentTab.historyLength))
	sb.WriteString(fmt.Sprintf("Window size: %v x %v | ", currentTab.screenInfo.Width, currentTab.screenInfo.Height))
	sb.WriteString("WillSmith v.")
	sb.WriteString(VersionName)
	sb.WriteString(" | ")
	if currentTab.PendingRequests > 0 {
		sb.WriteString(fmt.Sprintf("Pending requests: %v", currentTab.PendingRequests))
	}
	if(sb.Len() >= currentTab.screenInfo.Width) {
		return sb.String()[0:currentTab.screenInfo.Width-1]
	}
	return sb.String()
}

func ClearConsole() {
	var cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Writes a line of '=' to the width of the screen
func WriteLine(Width int) {
	for range Width {
		fmt.Print("=")
	}
	fmt.Println()
}
