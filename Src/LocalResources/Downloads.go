package localresources

import (
	globalstate "WillSmith/GlobalState"
	logger "WillSmith/Logger"
	"errors"
	"fmt"
	"os"
	"strings"
)

func Download(ResourceURL string, body []byte) {
	logger.SendWarning("flag")
	var fileName = getRequestedFileName(ResourceURL)
	var Path = getNewFilePath(fileName)
	var file, err = os.Create(Path)
	if err != nil {
		logger.SendError(fmt.Sprintf("Unable to create a file in downloads folder. Original error: %v", err.Error()))
		return
	}
	defer file.Close()
	var _, werr = file.Write(body)
	if werr != nil {
		logger.SendError(fmt.Sprintf("Unable to write in a file in downloads folder. Original error: %v", werr.Error()))
	}
}

func getRequestedFileName(FullURL string ) string {
	var parts = strings.Split(FullURL, "/")
	if len(parts) > 1 && parts[len(parts)-1] == "" {
		return parts[len(parts) - 2]
	}
	return parts[len(parts)-1]
}

func getNewFilePath(FileName string) string {
	var fullPath = fmt.Sprintf("%v/%v", globalstate.CurrentSettings.DownloadFolder, FileName)
	if !exists(fullPath) {
		return fullPath
	}
	var retry = 1
	var basePath, format = separateFormat(FileName)
	var newFileName string
	logger.SendInfo(fmt.Sprintf("%v dot %v", basePath, format))
	for {
		var newPath = strings.TrimRight(fmt.Sprintf("%v/%v(%d).%v", globalstate.CurrentSettings.DownloadFolder, basePath, retry, format), ".")
		logger.SendInfo(fmt.Sprintf("Checking %v", newPath))
		if !exists(newPath) {
			newFileName = newPath
			break
		}
		retry++
	}
	return newFileName
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func separateFormat(fullName string) (string, string) {
	logger.SendWarning(fullName)
	var parts = strings.Split(fullName, ".")
	if len(parts) < 2 {
		return fullName, ""
	}
	logger.SendInfo(fmt.Sprintf("Parts length: %d", len(parts)))
	for _, v := range parts {
		logger.SendInfo(v)
	}
	var base = strings.Join(parts[:len(parts)-2], "")
	return base, parts[len(parts)-1]
}
