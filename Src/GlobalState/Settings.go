package globalstate

import (
	"encoding/json"
	"fmt"
	"os"
)

const SettingsLocation = "../Settings.json"

type Settings struct {
	EnableLogging bool `json:"enablelogging"`
	CacheTTL int `json:"cachettl"`
}

var CurrentSettings Settings

var defaultSettings = Settings {
	EnableLogging: false,
	CacheTTL: 5,
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

