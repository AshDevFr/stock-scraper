package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"stock_scraper/pkg/config/parsers"
	"stock_scraper/types"
	"strings"
	"time"
)

var config types.Config
var userAgents []string

func GetConfig() types.Config {
	return config
}

func LoadConfig(filename string) types.Config {
	// load config.json file
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Fatal("Can't load config.json file with item numbers and email addresses.")
	}

	// unmarshal configs to Configs struct
	json.Unmarshal(file, &config)

	itemMap := make(map[string]*types.Item)
	var items []types.Item
	for _, item := range config.Items {
		parsedItem := parsers.ParseItem(config.DefaultConfig, item)
		items = append(items, parsedItem)
		itemMap[parsedItem.Uuid] = &parsedItem
	}
	config.Items = items
	config.ItemMap = itemMap

	log.Debug("Configs have been successfully loaded.")
	return config
}

func GetUserAgents() []string {
	return userAgents
}

func PickUserAgent() string {
	rand.Seed(time.Now().Unix())
	i := rand.Intn(len(userAgents))
	return userAgents[i]
}

func LoadUserAgents(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
		}
	}
	userAgents = strings.Split(string(content), "\n")
	return userAgents
}
