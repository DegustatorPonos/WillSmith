package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const VersionName string = "0.2a"
const HomePage string = "gemini://geminiprotocol.net/"

func main() {
	var history = make([]string, 1)
	var currntIndex = 0
	var PrevCurrentIndex = 0
	history[0] = HomePage
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
		fmt.Printf("Trying to access %v\n", history[currntIndex])
		var newwidth, newheight, _ = term.GetSize(0)
		if(newheight != height || newwidth != width) {
			height = newheight
			width = newwidth
			currentPage = ReadRequest(resp)
		}
		ClearConsole()
		fmt.Println(GetStatusBar(width, height, history[currntIndex], currntIndex))
		WriteLine(width)
		DisplayPage(currentPage)
		WriteLine(width)
		fmt.Print("Enter command: >")
		var command, _ = reader.ReadString('\n')
		var TrimmedCommand = strings.TrimRight(command, "\n")
		if (TrimmedCommand == "..") {
			fmt.Println("Detected backward")
			if(currntIndex >= 1) {
				currntIndex -= 1
			}
			continue
		}
		if strings.HasPrefix(TrimmedCommand, "gemini://") {
			fmt.Println("Detected link")
			if len(history) <= currntIndex + 1 {
				history = append(history, TrimmedCommand)
			} else {
				history[currntIndex + 1] = TrimmedCommand
			}
			currntIndex += 1
			resp = SendRequest(history[currntIndex], DEFAULT_PORT)
			currentPage = ReadRequest(resp)
			continue
		}
		if slices.Contains(currentPage.Links, TrimmedCommand) {
			fmt.Println("Detected link")
			if len(history) <= currntIndex + 1 {
				history = append(history, strings.Join([]string{history[currntIndex], TrimmedCommand, "/"}, ""))
			} else {
				history[currntIndex + 1] = strings.Join([]string{history[currntIndex], TrimmedCommand, "/"}, "")
			}
			currntIndex += 1
			resp = SendRequest(history[currntIndex], DEFAULT_PORT)
			currentPage = ReadRequest(resp)
			continue
		}
		fmt.Println(command)
	}
}

func ClearConsole() {
	var cmd = exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func GetStatusBar(ScreenWidth int, ScreenHeight int, URI string, HistoryLength int) string {
	var sb = strings.Builder{}
	sb.WriteString(URI)
	sb.WriteString(" | ")
	sb.WriteString("History: ")
	sb.WriteString(strconv.Itoa(HistoryLength))
	sb.WriteString(" | ")
	sb.WriteString("WillSmith v. ")
	sb.WriteString(VersionName)
	sb.WriteString(" | ")
	sb.WriteString("Window size: ")
	sb.WriteString(strconv.Itoa(ScreenWidth))
	sb.WriteString(" x ")
	sb.WriteString(strconv.Itoa(ScreenHeight))
	sb.WriteString(" | ")
	return sb.String()
}

func WriteLine(Width int) {
	for range Width {
		fmt.Print("=")
	}
	fmt.Println()
}
