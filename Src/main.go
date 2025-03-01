package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

const HomePage string = "gemini://geminiprotocol.net/"

func main() {
	var history = make([]string, 1)
	var currntIndex = 0
	history[0] = HomePage
	var reader = bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Trying to access %v\n", history[currntIndex])
		var resp = SendRequest(history[currntIndex], DEFAULT_PORT)
		var currentPage = ReadRequest(resp)
		ClearConsole()
		fmt.Printf("Current history lenth: %d\n", len(history))
		DisplayPage(currentPage)
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
		if slices.Contains(currentPage.Links, TrimmedCommand) {
			fmt.Println("Detected link")
			if len(history) <= currntIndex + 1 {
				history = append(history, strings.Join([]string{history[currntIndex], TrimmedCommand, "/"}, ""))
			} else {
				history[currntIndex + 1] = strings.Join([]string{history[currntIndex], TrimmedCommand, "/"}, "")
			}
			currntIndex += 1
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
