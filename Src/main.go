package main

import "fmt"

func main() {
	println("Hello, world!")
	var resp = SendRequest("gemini://geminiprotocol.net/", DEFAULT_PORT)
	fmt.Print(string(resp.Body))
}
