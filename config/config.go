package config

import (
	"encoding/json"
	"os"
)

/*
Configuration File (JSON)
{
	"credential.file": "path to google credendial file JSON",
	"token.file": "path to created token file",
	"delete.pattern": "search file string -> https://developers.google.com/drive/api/v3/reference/query-ref",
	"upload.max_tries": max number of attempts to upload the file, if fail
}
*/

// JSONConfig type
type JSONConfig struct {
	CredentialFile string `json:"credential.file"`
	TokenFile      string `json:"token.file"`
	DeletePattern  string `json:"delete.pattern"`
	UploadMaxTries int    `json:"upload.max_tries"`
}

// ReadConfigFile - Reads the JSON config file passed from command line
func ReadConfigFile(configFile string) (*JSONConfig, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	conf := &JSONConfig{}
	conf.UploadMaxTries = 3
	err = json.NewDecoder(f).Decode(conf)
	return conf, err
}
