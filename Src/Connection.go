package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Request struct {
	ResultCode byte
	Body []byte
}

const DEFAULT_PORT int = 1965

const ERR_HOST_NOT_FOUND string = "file://../StaticPages/Errors/NotFound"
const ERR_BODY_READ string = "file://../StaticPages/Errors/BodyErr"

// Sends a request to the server and returns a responce
func SendRequest(URI string, port int ) *Request{
	if(strings.HasPrefix(URI, "file")) {
		return ServeFile(URI)
	}
	var url_parsed, urlerr = url.Parse(URI)
	if urlerr != nil {
		fmt.Printf("Invalid URL")
		return ServeFile(ERR_HOST_NOT_FOUND)
	}
	var conn, err = tls.Dial("tcp", url_parsed.Host+":"+strconv.Itoa(port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return ServeFile(ERR_HOST_NOT_FOUND)
	}
	conn.Write([]byte(URI + "\r\n"))
	defer conn.Close()
	
	// Reading the header
	conn.Write([]byte(URI))
	var reader = bufio.NewReader(conn)
	var header, _ = reader.ReadString('\n')
	var RespCode, HeaderParsingErr = ParseResponceHeader(header)
	if HeaderParsingErr != nil {
		return ServeFile(ERR_BODY_READ)
	}
	var body, bodyReadingErr = io.ReadAll(reader)
	if bodyReadingErr != nil {
		return ServeFile(ERR_BODY_READ)
	} 

	var outp = Request{ResultCode: RespCode, Body: body}
	return &outp
}

func ParseResponceHeader(inp string) (byte, error) {
	var parts = strings.Split(inp, " ")
	if(len(parts) < 2) {
		return 0, fmt.Errorf("Invalid header")
	}
	var ResponceCode,CodeParsingErr = strconv.Atoi(parts[0])
	if CodeParsingErr != nil {
		return 0, fmt.Errorf("Error while parsing the request header")
	}
	return byte(ResponceCode), nil;
}

// Serves the file as a responce. Invoked when it starts with file://
func ServeFile(link string) *Request {
	var FilePath = strings.TrimPrefix(link, "file://")
	var file, fopenerr = os.ReadFile(FilePath)
	if fopenerr != nil {
		return &Request{ResultCode: 40}
	}
	var outp = Request{ResultCode: 20, Body: file}
	return &outp
}

func CompactAllBackwardsMotions(inp string) string {
	var outp = strings.Clone(inp)
	var r = regexp.MustCompile(`\/[^\/:]*\/\.\.\/`)
	for len(r.FindAllString(outp, 1)) > 0 {
		outp = r.ReplaceAllString(outp, "/")
	}
	return outp
}

func GoBackOneLayer(inp string) string {
	var r = regexp.MustCompile(`\/[^\/:]*\/?$`)
	return r.ReplaceAllString(inp, "")
}
