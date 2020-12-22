package scrapers

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"stock_scraper/types"
	"strings"
	"time"
)

func RunWeb(item types.Item, checkContent func(string, []types.ParsedResults) (string, error)) (types.Result, string, error) {
	itemUrl := item.TrackedUrl
	if itemUrl == "" {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("No url provided")
		return types.Result{}, "", errors.New("No url provided")
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
			chromedp.ActionFunc(func(ctx context.Context) error {
				time.Sleep(time.Second * 1)
				return nil
			}),
			chromedp.OuterHTML("html", &body),
		},
	)

	if err != nil {
		logger.Error(err)
		return types.Result{}, "", err
	}

	return processReader(item, logger, ioutil.NopCloser(strings.NewReader(body)), checkContent)
}

func getBrowserCtx(headless bool) (context.Context, context.CancelFunc) {
	if headless {
		return chromedp.NewContext(context.Background())
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	actx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return chromedp.NewContext(actx)
}
