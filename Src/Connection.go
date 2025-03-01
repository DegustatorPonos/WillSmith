package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	ResultCode byte
	Body []byte
}

const DEFAULT_PORT int = 1965

// Sends a request to the server and returns a responce
func SendRequest(URI string, port int ) *Request{
	if(strings.HasPrefix(URI, "file")) {
		return ServeFile(URI)
	}
	var url_parsed, urlerr = url.Parse(URI)
	if urlerr != nil {
		fmt.Printf("Invalid URL")
		return nil
	}
	var conn, err = tls.Dial("tcp", url_parsed.Host+":"+strconv.Itoa(port), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		fmt.Printf("An error occured while sending a request. \n| Original error message: %v\n", err.Error())
		return nil
	}
	conn.Write([]byte(URI + "\r\n"))
	defer conn.Close()
	
	// Reading the header
	conn.Write([]byte(URI))
	var reader = bufio.NewReader(conn)
	var header, _ = reader.ReadString('\n')
	var RespCode, HeaderParsingErr = ParseResponceHeader(header)
	if HeaderParsingErr != nil {
		fmt.Printf("An error occured while parsing a responce header. \n| Original error message: %v\n", HeaderParsingErr.Error())
		return nil
	}
	var body, bodyReadingErr = io.ReadAll(reader)
	if bodyReadingErr != nil {
		fmt.Printf("An error occured while reading a responce body. \n| Original error message: %v\n", bodyReadingErr.Error())
		return nil
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
