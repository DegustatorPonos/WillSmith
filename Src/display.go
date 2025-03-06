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
	Text []string
	Links []string
	L1Headers []string
	ScrollOffser uint
}

func ReadRequest(r *Request) *Page {
	var outp = Page{}
	outp.Text = make([]string, 0)
	outp.Links = make([]string, 0)
	var width, _, _ = term.GetSize(0)
	for _, str := range strings.Split(string(r.Body), "\n") {
		if strings.HasPrefix(str, "=>") {
			outp.Links = append(outp.Links, ParseLink(str))
		}
		if len(str) < width {
			outp.Text = append(outp.Text, str)
			continue
		}
		var rightSide = 0
		for i :=0; i < len(str); i = rightSide {
			rightSide = i + width
			if rightSide >= len(str) {
				rightSide = len(str)
				outp.Text = append(outp.Text, strings.Trim(str[i:rightSide], " "))
				continue
			}
			for(rightSide > 0 && str[rightSide] != '\n' && str[rightSide] != ' ') {
				rightSide--
			}
			outp.Text = append(outp.Text, strings.Trim(str[i:rightSide], " "))
		}
	}
	return &outp
}

func ParseLink(inp string) string {
	var prefixless, _ = strings.CutPrefix(inp, "=>")
	var pureLink = strings.Trim(prefixless, " ")
	if strings.Contains(prefixless, "	") {
		pureLink = strings.Split(prefixless, "	")[0]
	}
	if strings.Contains(pureLink, " ") {
		pureLink = strings.Split(prefixless, " ")[0]
	}
	var outp, _ = strings.CutSuffix(pureLink, "/")
	switch outp {
	case ".." :
		return "../"
	case "/" :
		return "//"
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
