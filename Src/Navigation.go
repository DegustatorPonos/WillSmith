package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const CMD_CHAN_BUFF_SIZE int = 1
const CMD_CHAN_ID int = 1

// Returns the path that does not contain the top layer
// Example: gemini://a/b/c/ -> gemini://a/b
func GoBackOneLayer(inp string) string {
	var r = regexp.MustCompile(`\/[^\/:]*\/?$`)
	return r.ReplaceAllString(inp, "")
}

// Returns a link that is a relative of ToAppend to BaseURI
func AppendToLink(BaseURI string, ToAppend string) string { 
	if(len(ToAppend) > 0 && ToAppend[0] == '/') {
		return strings.Join([]string{GetHostURI(BaseURI), strings.Replace(ToAppend, "/", "", 1), "/"}, "")
	}
	var newURI = ""
	if IsAnEndpoint(BaseURI) {
		BaseURI = GoBackOneLayer(BaseURI)
	}
	if strings.HasSuffix(BaseURI, "/") || strings.HasPrefix(ToAppend, "/") {
		newURI = strings.Join([]string{BaseURI, ToAppend}, "")
	} else {
		newURI = strings.Join([]string{BaseURI, "/", ToAppend}, "")
	}
	if !strings.HasSuffix(newURI, ".gmi") {
		newURI = strings.Join([]string{newURI, "/"}, "")
	}
	return CompactAllBackwardsMotions(newURI)
}

// Returns the host name of the URI
func GetHostURI(URI string) string {
	// This regex returns the mase address. For example, 
	// gemini://gemini.circumlunar.space/capcom/ returns gemini://gemini.circumlunar.space/
	var r = regexp.MustCompile(`gemini:\/\/[^\/:]*\/`)
	return r.FindString(URI)
}

// Returns true if the link is a gmi file reference
func IsAnEndpoint(inp string) bool {
	var r = regexp.MustCompile(`\/[^\/:]*\.gmi\/?$`)
	return r.FindString(inp) != ""
}

// If the relative path contains "../" then this function will delete it with the coresponding path part
// Example: gemini://a/b/../c/ -> gemini://a/c/
func CompactAllBackwardsMotions(inp string) string {
	var outp = strings.Clone(inp)
	var r = regexp.MustCompile(`\/[^\/:]+\/\.\.\/`)
	for len(r.FindAllString(outp, 1)) > 0 {
		outp = r.ReplaceAllString(outp, "/")
	}
	return outp
}

func NavigationTask(output chan string, controlChannel *chan int) {
	var reader = bufio.NewReader(os.Stdin)
	for {
		var command, readErr = reader.ReadString('\n')
		if readErr != nil {
			panic("Error in the command reading coroutine.")
		}
		command = strings.Trim(command, "\n")
		fmt.Printf("Sending a command \"%v\"\n", command)
		output <- command
		*controlChannel <- CMD_CHAN_ID
	}
}

func CreateCommandChannel(controlChannel *chan int) chan string {
	var outpChannel = make(chan string, CMD_CHAN_BUFF_SIZE);
	go NavigationTask(outpChannel, controlChannel)
	return outpChannel
}
