package main

import (
	"fmt"
	"strings"
)

type Page struct {
	Text []string
	Links []string
	L1Headers []string
}

func ReadRequest(r *Request) *Page {
	var outp = Page{}
	outp.Text = make([]string, 0)
	outp.Links = make([]string, 0)
	for _, str := range strings.Split(string(r.Body), "\n") {
		if strings.HasPrefix(str, "=> ") {
			outp.Links = append(outp.Links, ParseLink(str))
		}
		outp.Text = append(outp.Text, str)
	}
	return &outp
}

func ParseLink(inp string) string {
	var prefixless, _ = strings.CutPrefix(inp, "=> ")
	var pureLink = strings.Split(prefixless, "	")[0]
	var outp, _ = strings.CutSuffix(pureLink, "/")
	return outp
}

func DisplayPage(page *Page) {
	for _, str := range page.Text {
		fmt.Println(str)
	}
	return
	
	fmt.Println("LINKS: ")
	for _, str := range page.Links{
		fmt.Println(str)
	}
}
