package scrapers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"stock_scraper/pkg/utils"
	"stock_scraper/types"
	"strings"
	"time"
)

func Run(item types.Item, checkContent func(string, *types.Price, map[string]string) (string, error)) (string, string, error) {
	itemUrl := item.TrackedUrl
	if itemUrl == "" {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("No url provided")
		return "", "", errors.New("No url provided")
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
		return "", "", err
	}

	req.Header.Set("Accept", "text/html")
	for k, v := range item.Config.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", item.Config.UserAgent)

	res, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return "", "", err
	} else if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Request error (%d)", res.StatusCode))
		logger.Error(err)
		return "", "", err
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
		return "", "", err
	}

	content := ""
	results := make(map[string]string)
	for _, selector := range item.Config.Selectors {
		selection := doc.Find(selector).First()
		text := strings.TrimSpace(selection.Text())
		results[selector] = text
		content = fmt.Sprintf("%s %s", content, text)
	}

	var price *types.Price
	if item.Config.PriceSelector != "" {
		selection := doc.Find(item.Config.PriceSelector).First()
		text := strings.TrimSpace(selection.Text())
		if text != "" {
			price = utils.ParsePrice(text)
		}
	}

	body := doc.Find("html").Text()
	warn, err := checkContent(body, price, results)
	if err != nil {
		logger.Error(err)
	}
	if warn != "" {
		logger.Warn(warn)
		return "", warn, nil
	}

	logger.WithFields(log.Fields{
		"content": content,
	}).Info("Success")

	return content, "", nil
}
