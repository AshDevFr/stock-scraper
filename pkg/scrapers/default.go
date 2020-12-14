package scrapers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
	"stock_scraper/types"
)

func RunDefault(item types.Item) (string, string) {
	return Run(item, []string{"body"}, make(map[string]string))
}

func Run(item types.Item, defaultSelectors []string, extraHeaders map[string]string) (string, string) {
	itemUrl := item.Url
	if itemUrl == "" {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("No url provided")
		return "", "No url provided"
	}

	selectors := defaultSelectors

	if len(item.Selectors) > 0 {
		selectors = item.Selectors
	}

	client := &http.Client{}

	logger := log.WithFields(log.Fields{
		"url":       itemUrl,
		"selectors": selectors,
	})

	logger.Info("Fetching")

	req, err := http.NewRequest("GET", itemUrl, nil)
	if err != nil {
		logger.Error(err)
		return "", fmt.Sprintf("%s", err)
	}

	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"
	if item.UserAgent != "" {
		userAgent = item.UserAgent
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

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
