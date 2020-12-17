package scrapers

import (
	"compress/gzip"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"stock_scraper/types"
	"time"
)

func Run(item types.Item, checkContent func(string) error) (string, string) {
	itemUrl := item.TrackedUrl
	if itemUrl == "" {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("No url provided")
		return "", "No url provided"
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: false,
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}

	logger := log.WithFields(log.Fields{
		"url":       itemUrl,
		"selectors": item.Config.Selectors,
	})

	logger.Debug("Fetching")

	req, err := http.NewRequest("GET", itemUrl, nil)
	if err != nil {
		logger.Error(err)
		return "", fmt.Sprintf("%s", err)
	}

	req.Header.Set("Accept", "text/html")
	for k, v := range item.Config.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", item.Config.UserAgent)

	res, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return "", fmt.Sprintf("%s", err)
	} else if res.StatusCode != 200 {
		msg := fmt.Sprintf("Request error (%d)", res.StatusCode)
		logger.Error(msg)
		return "", msg
	}

	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		defer reader.Close()
	default:
		reader = res.Body
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Error(err)
		return "", fmt.Sprintf("%s", err)
	}

	body := doc.Find("html").Text()
	err = checkContent(body)
	if err != nil {
		logger.Error(err)
	}

	content := ""
	for _, selector := range item.Config.Selectors {
		selection := doc.Find(selector).First()
		content = fmt.Sprintf("%s %s", content, selection.Text())
	}

	logger.WithFields(log.Fields{
		"content": content,
	}).Info("Success")

	return content, ""
}
