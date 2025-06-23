package renderers

import (
	"encoding/json"
	"fmt"
	"os"

	logger "WillSmith/Logger"
)

const HomePageFile string = "../StaticPages/IndexPage"
const BookmarkCache string = "../StaticPages/Bookmarks.json"

type bookmarkList struct {
	Data []Bookmark `json:"data"`
}

type Bookmark struct {
	URL string `json:"url"`
	Description string `json:"description"`
}

var bookmarks []Bookmark = nil

func GetIndexPage() []byte {
	var file, fopenerr = os.ReadFile(HomePageFile)
	if fopenerr != nil {
		var errMsg = fmt.Sprintf("An error occured while trying to read the home page. \nOriginal error: %v\n", fopenerr.Error())
		return []byte(errMsg)
	}
	readBookmarks()
	for _, bm := range bookmarks {
		file = fmt.Appendf(file, "=> %v - %v\n", bm.URL, bm.Description)
	}
	// bookmarks = append(bookmarks, Bookmark{URL: "test", Description: "Test bookmark"})
	// updateBookmarks()
	return file
}

func readBookmarks() {
	if bookmarks != nil {
		return
	}
	var file, fopenerr = os.ReadFile(BookmarkCache)
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
	bookmarks = data.Data
}

func updateBookmarks() error {
	var file, err = os.Create(BookmarkCache)
	if err != nil {
		return err
	}
	defer file.Close()
	var inititalData = bookmarkList{Data: bookmarks}
	var body, jsonerr = json.MarshalIndent(inititalData, "", "	")
	if jsonerr != nil {
		return jsonerr
	}
	fmt.Fprint(file, string(body))
	return nil
}

func createBookmarkFile() (bool, []byte ) {
	var file, err = os.Create(BookmarkCache)
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
	bookmarks = append(bookmarks, toAdd)
	updateBookmarks()
}

func DeleteBookmark(URL string) {
	logger.SendInfo(fmt.Sprintf("Deleting %v from bookmarks", URL))
	for i, val := range bookmarks {
		if val.URL == URL {
			bookmarks = append(bookmarks[:i], bookmarks[i+1:]...)
			updateBookmarks()
			return
		}
	}
}
