package scrapers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"stock_scraper/internal/utils"
	"stock_scraper/types"
	"strings"
	"time"
)

func processDoc(item types.Item, reader io.ReadCloser) ([]types.ParsedResults, string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Error(err)
		return []types.ParsedResults{}, "", err
	}

	var results []types.ParsedResults
	itemSelector := item.Config.ItemSelector
	if itemSelector == "" {
		itemSelector = "html"
	}

	itemList := doc.Find(itemSelector)

	itemList.Each(func(i int, s *goquery.Selection) {
		result := types.ParsedResults{Result: make(map[string]string)}
		for _, selector := range item.Config.Selectors {
			selection := s.Find(selector)
			text := strings.TrimSpace(selection.Text())
			result.Result[selector] = text
			result.Content = fmt.Sprintf("%s %s", result.Content, text)
		}

		if item.Config.ItemLinkSelector != "" {
			selection := s.Find(item.Config.ItemLinkSelector)
			url, exists := selection.Attr("href")
			if exists && url != "" {
				trackedUrl, addToCartUrl := item.Parser.ParseUrls(types.Item{Url: url}, item.TrackedUrl)
				result.ItemLink = &trackedUrl
				result.ItemAddToCartLink = &addToCartUrl
			}
		}

		if item.Config.PriceSelector != "" {
			selection := s.Find(item.Config.PriceSelector)
			text := strings.TrimSpace(selection.Text())
			if text != "" {
				result.Price = utils.ParsePrice(text)
			}
		}

		results = append(results, result)
	})

	return results, doc.Find("html").Text(), nil
}

func processPrices(item types.Item, results []types.ParsedResults) string {
	prices := 0
	validPrices := 0
	for _, result := range results {
		if result.Price != nil {
			prices++
		}
		if result.Price != nil && item.Config.MaxPrice != nil {
			if result.Price.Value <= *item.Config.MaxPrice {
				validPrices++
			}
		}
	}

	if prices == 0 || validPrices > 0 {
		return ""
	}

	warn := "No valid price found"
	return warn
}

func getContent(results []types.ParsedResults) (string, string) {
	content := ""
	prices := ""
	for _, result := range results {
		if content == "" {
			content = result.Content
		} else {
			content = fmt.Sprintf("%s|%s", content, result.Content)
		}

		if result.Price != nil {
			if prices == "" {
				prices = fmt.Sprintf("%s%.2f", result.Price.Symbol, result.Price.Value)
			} else {
				prices = fmt.Sprintf("%s|%s%.2f", prices, result.Price.Symbol, result.Price.Value)
			}
		}

	}
	return content, prices
}

func processReader(item types.Item, logger *log.Entry, reader io.ReadCloser, checkContent types.CheckContentFunc) (types.Result, string, error) {
	results, body, err := processDoc(item, reader)
	warn := processPrices(item, results)
	if warn != "" {
		logger.Warn(warn)
		return types.Result{Results: results}, warn, nil
	}

	warn, err = checkContent(body, results)
	if err != nil {
		logger.Error(err)
		return types.Result{Results: results}, "", err
	}
	if warn != "" {
		logger.Warn(warn)
		return types.Result{Results: results}, warn, nil
	}

	content, prices := getContent(results)
	logger.WithFields(log.Fields{
		"content": content,
		"prices":  prices,
	}).Info("Success")

	return types.Result{
		Content: content,
		Results: results,
	}, "", nil
}

func Run(item types.Item, checkContent types.CheckContentFunc) (types.Result, string, error) {
	if item.Config.RunWeb != nil && *item.Config.RunWeb {
		return RunWeb(item, checkContent)
	}

	itemUrl := item.TrackedUrl
	if itemUrl == "" {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("No url provided")
		return types.Result{}, "", errors.New("No url provided")
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
		return types.Result{}, "", err
	}

	// Make sure we close the request to prevent EOF errors
	req.Close = true
	req.Header.Set("Accept", "text/html")
	for k, v := range item.Config.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", item.Config.UserAgent)

	res, err := client.Do(req)
	if err != nil {
		logger.Error(err)
		return types.Result{}, "", err
	} else if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Request error (%d)", res.StatusCode))
		logger.Error(err)
		return types.Result{}, "", err
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

	return processReader(item, logger, reader, checkContent)
}
