package main

import (
	"fmt"
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
	var prefixless, _ = strings.CutPrefix(inp, "=> ")
	var pureLink = strings.Split(prefixless, "	")[0]
//	pureLink = strings.Split(prefixless, " ")[0]
	var outp, _ = strings.CutSuffix(pureLink, "/")
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
