package scrapers

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"stock_scraper/types"
	"strings"
	"time"
)

type BrowserSession struct {
	Context context.Context
	Cancel  context.CancelFunc
}

var retryWaitingSec = 1

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

	browserSession := getBrowserCtx(item)
	defer browserSession.Cancel()

	var body string
	err := chromedp.Run(browserSession.Context,
		chromedp.Tasks{
			network.Enable(),
			network.SetExtraHTTPHeaders(headers),
			chromedp.Navigate(itemUrl),
			chromedp.ActionFunc(func(ctx context.Context) error {
				retries := *item.Config.WebRetries
				for retries > 0 {
					chromedp.JavascriptAttribute("html", "outerHTML", &body).Do(ctx)
					_, err := checkContent(body, []types.ParsedResults{})
					if err != nil {
						logger.Warn(fmt.Sprintf("%s Waiting %ds...", err, retryWaitingSec))
						time.Sleep(time.Second * time.Duration(retryWaitingSec))
					} else {
						return nil
					}
					retries--
				}
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

func getBrowserCtx(item types.Item) BrowserSession {
	var ctx context.Context
	var cancel context.CancelFunc
	if !item.Config.NoHeadless {
		ctx, cancel = chromedp.NewContext(context.Background())
	} else {
		opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
		actx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel = chromedp.NewContext(actx)
	}

	return BrowserSession{Context: ctx, Cancel: cancel}
}
