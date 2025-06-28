package geminiprotocol

import (
	globalstate "WillSmith/GlobalState"
	logger "WillSmith/Logger"
	"fmt"
	"strings"
)


func ConnectionTask(RequestChan *chan RequestCommand, ResponceChan *chan *Request, DownloadChan *chan *Request, TerminationChan *chan bool, controlChannel *chan int) {
	defer close(*ResponceChan)
	var PendingRequests = make([]string, 0)
	var PendngRequestsChan = make(chan *Request, globalstate.State.ChannelLengths.ConnectionBuffer)
	for {
		select {
		case req := <-*RequestChan:
			logger.SendInfo(fmt.Sprintf("Requesting \"%v\"", req.URL))

			PendingRequests = append(PendingRequests, req.URL)
			// Checking for a page in cashe
			if req.MandatoryReload {
				Cache.InvalidatePage(req.URL)
			}
			var CachedPage = Cache.GetPageFromCache(req.URL)
			if CachedPage != nil {
				logger.SendInfo(fmt.Sprintf("Retrived \"%v\" from cashe", req.URL))
				CachedPage.target = req.TargetAction
				PendngRequestsChan <- CachedPage
				continue
			}

			// Sending a request here
			go GetPageTask(req.URL, req.TargetAction, &PendngRequestsChan)
			continue
		
		case <-*TerminationChan:
			// Clearing all pending requests
			PendingRequests = make([]string, 0)
			continue

		case resp := <- PendngRequestsChan:
			// Checking if the page we recived was requested or we got an error page
			logger.SendInfo(fmt.Sprintf("Retrived \"%v\"", resp.URI))
			if resp.target == DOWNLOAD {
				*DownloadChan <- resp
				*controlChannel <- DOWNLOAD_CHAN_ID
				continue
			}
			if len(PendingRequests) != 0 && strings.HasPrefix(resp.URI, "file://../StaticPages/Errors/") {
				*ResponceChan <- resp
				*controlChannel <- CON_CHAN_ID
			}
			for i, val := range PendingRequests {
				if(val == resp.URI) {
					if i != len(PendingRequests) - 1 {
						PendingRequests = append(PendingRequests[:i], PendingRequests[i+1:]...)
					} else {
						PendingRequests = PendingRequests[:i]
					}
					*ResponceChan <- resp
					*controlChannel <- CON_CHAN_ID
				}
			}
			Cache.AddPage(*resp)
			continue

		}
	}
}

// A coroutine summoned by the ConnectionTask
func GetPageTask(URI string, target TargetActionType, ResponceChan *chan *Request) {
	var resp = SendRequest(URI, DEFAULT_PORT)
	resp.target = target
	*ResponceChan <- resp
}

func CreateConnectionTask(RequestChan *chan RequestCommand, TerminationChan *chan bool, controlChannel *chan int) (*chan *Request, *chan *Request) {
	var outpChannel = make(chan *Request, globalstate.State.ChannelLengths.ConnectionBuffer)
	var downlaodChannel = make(chan *Request, globalstate.State.ChannelLengths.DownloadBuffer)
	go ConnectionTask(RequestChan, &outpChannel, &downlaodChannel, TerminationChan, controlChannel)
	return &outpChannel, &downlaodChannel
}
