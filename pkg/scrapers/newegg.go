package scrapers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"regexp"
	"stock_scraper/types"
)

var (
	defaultNeweggSelectors = []string{
		"#ProductBuy",
	}
)

func getNeweggItemId(item types.Item) string {
	u, err := url.Parse(item.Url)
	if err != nil {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("Invalid Url")
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err == nil {
		if itemId, ok := params["Item"]; ok {
			return itemId[0]
		}
	}

	r := regexp.MustCompile(".*/p/([^\\/]+)")
	match := r.FindStringSubmatch(u.Path)

	if len(match) > 1 {
		if match[1] != "" {
			return match[1]
		}
	}

	return ""
}

func RunNewegg(item types.Item) (string, string) {
	itemId := item.Id
	if itemId == "" {
		if item.Url == "" {
			log.WithFields(log.Fields{
				"item": item,
			}).Error("Invalid configuration")
			return "", "Invalid configuration"
		}

		itemId = getNeweggItemId(item)
	}

	itemUrl := "https://www.newegg.com/Product/Product.aspx?Item=" + itemId
	if itemId == "" {
		log.WithFields(log.Fields{
			"item": item,
		}).Warn("Could not find the item ID")
		itemUrl = item.Url
	}

	selectors := defaultNeweggSelectors

	if len(item.Selectors) > 0 {
		selectors = item.Selectors
	}

	client := &http.Client{}

	logger := log.WithFields(log.Fields{
		"website":   "NewEgg",
		"id":        itemId,
		"selectors": selectors,
	})

	logger.Info("Fetching")

	req, err := http.NewRequest("GET", itemUrl, nil)
	if err != nil {
		logger.Error(err)
		return "", fmt.Sprintf("%s", err)
	}

	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
	if item.UserAgent != "" {
		userAgent = item.UserAgent
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return "", fmt.Sprintf("%s", err)
	} else if res.StatusCode != 200 {
		msg := fmt.Sprintf("Request error (%d)", res.StatusCode)
		logger.Info(msg)
		return "", msg
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error(err)
		return "", fmt.Sprintf("%s", err)
	}

	content := ""
	for _, selector := range selectors {
		selection := doc.Find(selector).First()
		content = fmt.Sprintf("%s %s", content, selection.Text())
	}

	logger.WithFields(log.Fields{
		"content": content,
	}).Info("Found")
	return content, ""
}
