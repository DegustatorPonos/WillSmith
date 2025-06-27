package localresources

import (
	globalstate "WillSmith/GlobalState"
	logger "WillSmith/Logger"
	"encoding/json"
	"fmt"
	"os"
)

type bookmarkList struct {
	Data []Bookmark `json:"data"`
}

type Bookmark struct {
	URL string `json:"url"`
	Description string `json:"description"`
}

var Bookmarks []Bookmark = nil

func ReadBookmarks() {
	if Bookmarks != nil {
		return
	}
	var file, fopenerr = os.ReadFile(globalstate.CurrentSettings.BookmarksFile)
	if fopenerr != nil {
		var success, inintBody = createBookmarkFile()
		if success {
			file = inintBody
		} else {
			panic(fmt.Sprintf("Unable to load or create bookmark file. \nOriginal error: %v\n", fopenerr.Error()))
		}
	}
	var data = bookmarkList{}
	json.Unmarshal(file, &data)
	Bookmarks = data.Data
}

func updateBookmarks() error {
	var file, err = os.Create(globalstate.CurrentSettings.BookmarksFile)
	if err != nil {
		return err
	}
	defer file.Close()
	var inititalData = bookmarkList{Data: Bookmarks}
	var body, jsonerr = json.MarshalIndent(inititalData, "", "	")
	if jsonerr != nil {
		return jsonerr
	}
	fmt.Fprint(file, string(body))
	return nil
}

func createBookmarkFile() (bool, []byte ) {
	var file, err = os.Create(globalstate.CurrentSettings.BookmarksFile)
	if err != nil {
		return false, nil
	}
	defer file.Close()
	var inititalData = bookmarkList{Data: []Bookmark{}}
	var body, jsonerr = json.MarshalIndent(inititalData, "", "	")
	if jsonerr != nil {
		return false, nil
	}
	fmt.Fprint(file, string(body))
	return true, body
}

func AddBookmark(toAdd Bookmark) {
	logger.SendInfo(fmt.Sprintf("Adding %v to bookmarks", toAdd.URL))
	Bookmarks = append(Bookmarks, toAdd)
	updateBookmarks()
}

func DeleteBookmark(URL string) {
	logger.SendInfo(fmt.Sprintf("Deleting %v from bookmarks", URL))
	for i, val := range Bookmarks {
		if val.URL == URL {
			Bookmarks = append(Bookmarks[:i], Bookmarks[i+1:]...)
			updateBookmarks()
			return
		}
	}
}

