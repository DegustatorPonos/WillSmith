package renderers

import (
	"fmt"
	"os"

	localresources "WillSmith/LocalResources"
)

const HomePageFile string = "../StaticPages/IndexPage"

func GetIndexPage() []byte {
	var file, fopenerr = os.ReadFile(HomePageFile)
	if fopenerr != nil {
		var errMsg = fmt.Sprintf("An error occured while trying to read the home page. \nOriginal error: %v\n", fopenerr.Error())
		return []byte(errMsg)
	}
	localresources.ReadBookmarks()
	for _, bm := range localresources.Bookmarks {
		file = fmt.Appendf(file, "=> %v - %v\n", bm.URL, bm.Description)
	}
	return file
}
