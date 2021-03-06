package parsers

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"stock_scraper/internal/scrapers"
	"stock_scraper/internal/utils"
	"stock_scraper/types"
)

type AmazonParser struct {
	label string
}

func (p *AmazonParser) ParseId(item types.Item) string {
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

	r := regexp.MustCompile("(?:dp|o|gp|-)\\/(B[0-9]{2}[0-9A-Z]{7}|[0-9]{9}(?:X|[0-9]))")
	match := r.FindStringSubmatch(u.Path)

	if len(match) > 1 {
		if match[1] != "" {
			return match[1]
		}
	}

	return ""
}

func (p *AmazonParser) ParseUrls(item types.Item, trackedUrl string) (string, string) {
	itemId := item.Id
	if itemId == "" {
		itemId = p.ParseId(item)
	}
	if itemId == "" {
		return utils.CompleteUrl(item.Url, trackedUrl), ""
	}

	addToCartUrl := "https://smile.amazon.com/gp/aws/cart/add-res.html?ASIN.1=" + itemId + "&Quantity.1=1"
	if item.Config.ForceUrl {
		return utils.CompleteUrl(item.Url, trackedUrl), addToCartUrl
	}
	return "https://smile.amazon.com/dp/" + itemId, addToCartUrl
}

func (p *AmazonParser) Label() string {
	return p.label
}

func (p *AmazonParser) Parse(defaultConfig types.ItemConfig, item types.Item) types.Item {
	item.TrackedUrl, item.AddToCartUrl = p.ParseUrls(item, item.Url)

	item.Config.Selectors = ParseSelectors(defaultConfig, item)
	if len(item.Config.Selectors) == 0 {
		item.Config.Selectors = []string{
			"span#price_inside_buybox",
			"input#add-to-cart-button",
		}
	}

	if item.Config.PriceSelector == "" {
		item.Config.PriceSelector = "span#price_inside_buybox"
	}

	defaultHeaders := make(map[string]string)
	item.Config.Headers = ParseHeaders(defaultConfig, defaultHeaders, item)

	return item
}

func checkAmazonContent(body string, results []types.ParsedResults) (string, error) {
	return "", nil
}

func (p *AmazonParser) Run(item types.Item) (types.Result, string, error) {
	return scrapers.Run(item, checkAmazonContent)
}
