package scrapers

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"stock_scraper/types"
	"strings"
)

func RunWeb(item types.Item) (string, string) {
	itemUrl := item.TrackedUrl
	if itemUrl == "" {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("No url provided")
		return "", "No url provided"
	}

	logger := log.WithFields(log.Fields{
		"url":       itemUrl,
		"selectors": item.Config.Selectors,
	})

	logger.Debug("Fetching")

	headers := make(map[string]interface{})
	headers["Accept"] = "text/html"
	for k, v := range item.Config.Headers {
		headers[k] = v
	}
	headers["User-Agent"] = item.Config.UserAgent

	ctx, cancel := getBrowserCtx(true)
	defer cancel()

	var body string
	err := chromedp.Run(ctx,
		chromedp.Tasks{
			network.Enable(),
			network.SetExtraHTTPHeaders(headers),
			chromedp.Navigate(itemUrl),
			chromedp.WaitVisible(item.Config.Selectors[0]),
			chromedp.OuterHTML("html", &body),
		},
	)

	if err != nil {
		logger.Error(err)
		return "", fmt.Sprintf("%s", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		log.Error(err)
		return "", fmt.Sprintf("%s", err)
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

func getBrowserCtx(headless bool) (context.Context, context.CancelFunc) {
	if headless {
		return chromedp.NewContext(context.Background())
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	actx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return chromedp.NewContext(actx)
}
