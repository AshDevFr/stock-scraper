package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"stock_scraper/types"
)

func LoadConfig(filename string) types.Config {
	var config types.Config
	// load config.json file
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Fatal("Can't load config.json file with item numbers and email addresses.")
	}

	// unmarshal configs to Configs struct
	json.Unmarshal(file, &config)

	log.Println("Configs have been successfully loaded.")
	return config
}
