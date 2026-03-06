package geminiprotocol

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

	globalstate "WillSmith/GlobalState"
	logger "WillSmith/Logger"
	renders "WillSmith/Renderers"
)

type TargetActionType int

const (
	RENDER TargetActionType = iota 
	DOWNLOAD
)

type Request struct {
	URI string
	ResultCode byte
	Body []byte
	target TargetActionType
}

type RequestCommand struct {
	URL string
	MandatoryReload bool
	TargetAction TargetActionType
}

type SpecialPage struct {
	Name string
	RenderFunc func()[]byte
}

const CON_CHAN_ID int = 2
const DOWNLOAD_CHAN_ID int = 21

const ERR_HOST_NOT_FOUND string = "file://../StaticPages/Errors/NotFound"
const ERR_BODY_READ string = "file://../StaticPages/Errors/BodyErr"
const ERR_TEMP_FALIURE string = "file://../StaticPages/Errors/TempFaliure"
const ERR_PERMA_ERROR string = "file://../StaticPages/Errors/PermaError"
const ERR_CLIENT_CERTS string = "file://../StaticPages/Errors/ClientCerts"
const ERR_INPUT_EXPECTED string = "file://../StaticPages/Errors/InputExpected"

const DEFAULT_PORT int = 1965

var SpecialPages = []SpecialPage {
	SpecialPage {
		Name: "../StaticPages/IndexPage",
		RenderFunc: renders.GetIndexPage,
	},
}

var conf = &tls.Config{
	InsecureSkipVerify: true,
}

var dialer = tls.Dialer{
	Config: conf,
}

var Cache PagesCache = PagesCache{CachedPages: make(map[string]CachedPage)}

// Sends a request to the server and returns a responce
func SendRequest(URI string, port int, verify bool) *Request{
	if(strings.HasPrefix(URI, "file")) {
		return ServeFile(URI, URI)
	}
	var url_parsed, urlerr = url.Parse(URI)
	if urlerr != nil {
		return ServeFile(ERR_HOST_NOT_FOUND, URI)
	}
	// The provider specified a custom port
	if strings.Contains(url_parsed.Host, ":") {
		var parseErr error
		url_parsed.Host, port, parseErr = splitHostname(url_parsed.Host)
		if parseErr != nil {
			return generateErrorResponce(URI, parseErr)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(globalstate.CurrentSettings.ConnectionTimeout) * time.Second)
	var conn, err = dialer.DialContext(ctx, "tcp", url_parsed.Host+":"+strconv.Itoa(port))
	cancel()
	if err != nil {
		return generateErrorResponce(URI, err)
	}

	conn.Write([]byte(URI + "\r\n"))
	defer conn.Close()
	
	// Reading the header
	conn.Write([]byte(URI))
	var reader = bufio.NewReader(conn)
	var header, _ = reader.ReadString('\n')
	var RespCode byte
	if verify {
		RespCode, HeaderParsingErr := ParseResponceHeader(header)
		if(RespCode < 20 || RespCode > 29) {
			return GetErrorMessage(int(RespCode), URI)
		}
		if HeaderParsingErr != nil {
			return generateErrorResponce(URI, HeaderParsingErr)
		}
	}
	var body, bodyReadingErr = io.ReadAll(reader)
	logger.SendInfo(fmt.Sprintf("Got %d bytes from gemispace (verified: %v)", len(body), verify))
	if bodyReadingErr != nil {
		return generateErrorResponce(URI, bodyReadingErr)
	} 

	var outp = Request{URI: URI, ResultCode: RespCode, Body: body}
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
func ServeErrorMessage(errorPage string, link string) *Request {
	var FilePath = strings.TrimPrefix(errorPage, "file://")
	var file, fopenerr = os.ReadFile(FilePath)
	if fopenerr != nil {
		return &Request{ResultCode: 40}
	}
	var outp = Request{ResultCode: 20, Body: file, URI: link}
	return &outp
}

func generateErrorResponce(link string, err error) *Request {
	return &Request {
		ResultCode: 100,
		URI: link,
		Body: renders.CreateErrorWrapper(err)(),
	}
}

// Serves the file as a responce. Should be invoked when it starts with file://
func ServeFile(link string, baseURI string) *Request {
	var FilePath = strings.TrimPrefix(link, "file://")
	var isSpecial, renderFunc = GetSpecificRenderer(FilePath)
	if isSpecial {
		var file = renderFunc()
		var outp = Request{ResultCode: 20, Body: file, URI: link}
		return &outp
	}
	var file, fopenerr = os.ReadFile(FilePath)
	if fopenerr != nil {
		return &Request{ResultCode: 40}
	}
	var outp = Request{ResultCode: 20, Body: file, URI: baseURI}
	return &outp
}

// Returns the error message responce that corresponds to a responce code
func GetErrorMessage(errorCode int, connectedURL string) *Request {
	if(errorCode < 20) {
		return ServeErrorMessage(ERR_INPUT_EXPECTED, connectedURL)
	}
	if(errorCode >= 40 && errorCode < 40) {
		return ServeErrorMessage(ERR_TEMP_FALIURE, connectedURL)
	}
	if(errorCode >= 50 && errorCode < 60) {
		return ServeErrorMessage(ERR_PERMA_ERROR, connectedURL)
	}
	if(errorCode >= 60 && errorCode < 70) {
		return ServeErrorMessage(ERR_CLIENT_CERTS, connectedURL)
	}
	return ServeErrorMessage(ERR_BODY_READ, connectedURL)
}

func GetSpecificRenderer(fileName string) (bool, func()[]byte) {
	for _, v := range SpecialPages {
		logger.SendInfo(fmt.Sprintf("Comparing file links %v and %v\n", fileName, v.Name))
		if v.Name == fileName {
			return true, v.RenderFunc
		}
	}
	return false, nil
}

func ShouldVerify(requestType TargetActionType) bool {
	switch requestType {
	case DOWNLOAD:
		return false
	}
	return true
}

// Assume that the hostname has a structure of host:port
func splitHostname(hostname string) (string, int, error) {
	var parts = strings.Split(hostname, ":")
	if len(parts) != 2 {
		return hostname, DEFAULT_PORT, fmt.Errorf("The hostname structure is invalid")
	}
	var port, err = strconv.Atoi(parts[1])
	return parts[0], port, err
}
