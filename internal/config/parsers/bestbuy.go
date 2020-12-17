package parsers

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"stock_scraper/internal/scrapers"
	"stock_scraper/types"
)

type BestBuyParser struct {
	label string
}

func getBestBuyItemId(item types.Item) string {
	u, err := url.Parse(item.Url)
	if err != nil {
		log.WithFields(log.Fields{
			"item": item,
		}).Error("Invalid Url")
	}

	params, err := url.ParseQuery(u.RawQuery)
	if err == nil {
		if itemId, ok := params["skuId"]; ok {
			return itemId[0]
		}
	}

	r := regexp.MustCompile(".*/(\\d+)\\.p\\??")
	match := r.FindStringSubmatch(u.Path)

	if len(match) > 1 {
		if match[1] != "" {
			return match[1]
		}
	}

	return ""
}

func (p *BestBuyParser) Label() string {
	return p.label
}

func (p *BestBuyParser) Parse(defaultConfig types.ItemConfig, item types.Item) types.Item {
	itemId := item.Id
	if itemId == "" && item.Url != "" {
		itemId = getBestBuyItemId(item)
	}

	if itemId != "" {
		item.TrackedUrl = "https://api.bestbuy.com/click/-/" + itemId + "/pdp"
		item.AddToCartUrl = "https://api.bestbuy.com/click/-/" + itemId + "/cart"
	}

	item.Config.Selectors = ParseSelectors(defaultConfig, item)
	if len(item.Config.Selectors) == 0 {
		item.Config.Selectors = []string{
			"div.row.v-m-bottom-g div.col-xs-12  div.fulfillment-add-to-cart-button  button.add-to-cart-button",
		}
	}

	if item.Config.PriceSelector == "" {
		item.Config.PriceSelector = "div > div > div > div > div.price-box > div:nth-child(1) > div > div > span:nth-child(1)"
	}

	defaultHeaders := make(map[string]string)
	defaultHeaders["accept-encoding"] = "gzip"
	//defaultHeaders["cache-control"] = "no-cache"
	defaultHeaders["accept-language"] = "en-US,en;q=0.9"
	//defaultHeaders["pragma"] = "no-cache"
	//defaultHeaders["sec-fetch-dest"] = "document"
	//defaultHeaders["sec-fetch-mode"] = "navigate"
	//defaultHeaders["sec-fetch-site"] = "none"
	//defaultHeaders["upgrade-insecure-requests"] = "1"
	item.Config.Headers = ParseHeaders(defaultConfig, defaultHeaders, item)

	return item
}

func (p *BestBuyParser) Run(item types.Item) (string, string, error) {
	return scrapers.Run(item, func(body string, price *types.Price, selectionTexts map[string]string) (string, error) {
		return "", nil
	})
}
