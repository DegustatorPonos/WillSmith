package main

// Everything related to server connections, routing and data retrival

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	ResultCode byte
	Body []byte
}

const DEFAULT_PORT int = 1965

const ERR_HOST_NOT_FOUND string = "file://../StaticPages/Errors/NotFound"
const ERR_BODY_READ string = "file://../StaticPages/Errors/BodyErr"
const ERR_TEMP_FALIURE string = "file://../StaticPages/Errors/TempFaliure"
const ERR_PERMA_ERROR string = "file://../StaticPages/Errors/PermaError"
const ERR_CLIENT_CERTS string = "file://../StaticPages/Errors/ClientCerts"
const ERR_INPUT_EXPECTED string = "file://../StaticPages/Errors/InputExpected"

var conf = &tls.Config{
	InsecureSkipVerify: true,
}

var dialer = tls.Dialer{
	Config: conf,
}

// Sends a request to the server and returns a responce
func SendRequest(URI string, port int) *Request{
	if(strings.HasPrefix(URI, "file")) {
		return ServeFile(URI)
	}
	var url_parsed, urlerr = url.Parse(URI)
	if urlerr != nil {
		return ServeFile(ERR_HOST_NOT_FOUND)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	var conn, err = dialer.DialContext(ctx, "tcp", url_parsed.Host+":"+strconv.Itoa(port))
	cancel()
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
	if(RespCode < 20 || RespCode > 29) {
		return GetErrorMessage(int(RespCode))
	}
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

// Returns the responce code and error (if any)
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

// Serves the file as a responce. Should be invoked when it starts with file://
func ServeFile(link string) *Request {
	var FilePath = strings.TrimPrefix(link, "file://")
	var file, fopenerr = os.ReadFile(FilePath)
	if fopenerr != nil {
		return &Request{ResultCode: 40}
	}
	var outp = Request{ResultCode: 20, Body: file}
	return &outp
}

// Returns the error message responce that corresponds to a responce code
func GetErrorMessage(errorCode int) *Request {
	if(errorCode < 20) {
		return ServeFile(ERR_INPUT_EXPECTED)
	}
	if(errorCode >= 40 && errorCode < 40) {
		return ServeFile(ERR_TEMP_FALIURE)
	}
	if(errorCode >= 50 && errorCode < 60) {
		return ServeFile(ERR_PERMA_ERROR)
	}
	if(errorCode >= 60 && errorCode < 70) {
		return ServeFile(ERR_CLIENT_CERTS)
	}
	return ServeFile(ERR_BODY_READ)
}

