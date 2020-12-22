package parsers

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"stock_scraper/internal/scrapers"
	"stock_scraper/types"
	"strings"
)

type NeweggParser struct {
	label string
}

func (p *NeweggParser) ParseId(item types.Item) string {
	if item.Id != "" {
		return item.Id
	}
	if item.Url == "" {
		return ""
	}

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

	r := regexp.MustCompile(".*/p/(N\\d+)")
	match := r.FindStringSubmatch(u.Path)

	if len(match) > 1 {
		if match[1] != "" {
			return match[1]
		}
	}

	return ""
}

func (p *NeweggParser) Label() string {
	return p.label
}

func (p *NeweggParser) Parse(defaultConfig types.ItemConfig, item types.Item) types.Item {
	if item.Id != "" {
		item.TrackedUrl = "https://www.newegg.com/Product/Product.aspx?Item=" + item.Id
		item.AddToCartUrl = "https://secure.newegg.com/Shopping/AddtoCart.aspx?Submit=ADD&ItemList=" + item.Id
	}

	item.Config.Selectors = ParseSelectors(defaultConfig, item)
	if len(item.Config.Selectors) == 0 {
		item.Config.Selectors = []string{
			"#ProductBuy",
		}
	}

	if item.Config.PriceSelector == "" {
		item.Config.PriceSelector = "#app > div.page-content > div.page-section > div > div > div.row-side > div.product-buy-box > div.product-pane > div.product-price > ul > li.price-current"
	}

	defaultHeaders := make(map[string]string)
	item.Config.Headers = ParseHeaders(defaultConfig, defaultHeaders, item)

	return item
}

func (p *NeweggParser) Run(item types.Item) (types.Result, string, error) {
	return scrapers.Run(item, func(body string, results []types.ParsedResults) (string, error) {
		if strings.Contains(strings.ToLower(body), strings.ToLower("Are you a human?")) {
			return "", errors.New("Anti bot recaptcha")
		}
		return "", nil
	})
}
