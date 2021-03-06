package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"stock_scraper/internal/browser"
	"stock_scraper/internal/process"
	"stock_scraper/internal/state"
	"stock_scraper/internal/websocket"
	"stock_scraper/types"
	"sync"
)

func runScraper(item types.Item) {
	parallelChan <- true
	scraperType := item.Parser.Label()
	result, warn, err := item.Parser.Run(item)
	defer func() { <-parallelChan }()

	if err != nil {
		websocket.SendUpdateMessage(scraperType, item, "error", fmt.Sprintf("%s", err))
	} else if warn != "" {
		websocket.SendUpdateMessage(scraperType, item, "warn", warn)
	} else {
		websocket.SendUpdateMessage(scraperType, item, "ok", result.Content)
		rules := config.DefaultConfig.Rules
		if item.Config.Rules != nil {
			rules = item.Config.Rules
		}
		previousContent := state.GetContent(item.Uuid)
		actions := process.ApplyRules(item, rules, previousContent, result)
		state.SetContent(item.Uuid, result)

		if len(actions) > 0 {
			if *item.Config.OpenLinks && item.Url != "" {
				state.ShouldRunAlert(item.Uuid, func() {
					for _, action := range actions {
						if *item.Config.OpenAddToCart && action.AddToCartLink != "" {
							browser.Open(action.AddToCartLink)
						} else if *item.Config.OpenAddToCart && action.AddToCartLink != "" {
							browser.Open(action.AddToCartLink)
						} else if action.Link != "" {
							browser.Open(action.Link)
						} else {
							browser.Open(item.Url)
						}
					}
				})
				log.Info("Opening link")
			}
			for _, action := range actions {
				log.WithFields(log.Fields{
					"url":         item.TrackedUrl,
					"selectors":   item.Config.Selectors,
					"addedText":   action.Diff.AddedText,
					"removedText": action.Diff.RemovedText,
				}).Debug("Action triggered")
				websocket.SendActionMessage(action, item)
			}
		}
	}
}

func runAllScrapers() {
	log.Debug("Running the scrapers: start")
	// loop for items in config to build and execute http requests

	var waitGroup sync.WaitGroup

	for _, item := range config.Items {
		waitGroup.Add(1)
		go func(item types.Item) {
			defer waitGroup.Done()
			runScraper(item)
		}(item)
	}
	waitGroup.Wait()
	log.Debug("Running the scrapers: done")
}
