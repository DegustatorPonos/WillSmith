package globalstate

import (
	"encoding/json"
	"fmt"
	"os"
)

const SettingsLocation = "../Settings.json"

type Settings struct {
	ConnectionTimeout int `json:"connectiontimeout"`
	EnableLogging bool `json:"enablelogging"`
	CacheTTL int `json:"cachettl"`
	DownloadFolder string `json:"downloadfolder"`
	BookmarksFile string  `json:"bookmarksfile"`
}

var CurrentSettings Settings

var defaultSettings = Settings {
	ConnectionTimeout: 10,
	EnableLogging: true,
	CacheTTL: 5,
	DownloadFolder: "../Dowloads",
	BookmarksFile: "../StaticPages/Bookmarks.json",
}

func ReadSettings() {
	var file, fopenerr = os.ReadFile(SettingsLocation)
	if fopenerr != nil {
		var success, inintBody = createSettingsFile()
		if success {
			file = inintBody
		} else {
			panic(fmt.Sprintf("Unable to load or create settings file. \nOriginal error: %v\n", fopenerr.Error()))
		}
	}
	var data = Settings{}
	json.Unmarshal(file, &data)
	var validationErr = data.Validate()
	if validationErr != nil {
		panic(fmt.Sprintf("Failed to load settings file: %s", validationErr.Error()))
	}
	CurrentSettings = data
}

func createSettingsFile() (bool, []byte ) {
	var file, err = os.Create(SettingsLocation)
	if err != nil {
		return false, nil
	}
	defer file.Close()
	var inititalData = defaultSettings
	var body, jsonerr = json.MarshalIndent(inititalData, "", "	")
	if jsonerr != nil {
		return false, nil
	}
	fmt.Fprint(file, string(body))
	return true, body
}

// Returns an error with the issue description if the settings are invalid
func (base *Settings) Validate() error {
	if len(base.BookmarksFile) == 0 {
		return fmt.Errorf("The length of the bookmarks file path must be greather than 0")
	}
	if len(base.DownloadFolder) == 0 {
		return fmt.Errorf("The length of the downlaods folder path must be greather than 0")
	}
	return nil
}
