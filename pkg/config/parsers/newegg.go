package parsers

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"stock_scraper/pkg/scrapers"
	"stock_scraper/types"
	"strings"
)

type NeweggParser struct {
	label string
}

func getNeweggItemId(item types.Item) string {
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

	r := regexp.MustCompile(".*/p/([^\\/]+)")
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
	itemId := item.Id
	if itemId == "" && item.Url != "" {
		itemId = getNeweggItemId(item)
	}

	if itemId != "" {
		item.TrackedUrl = "https://www.newegg.com/Product/Product.aspx?Item=" + itemId
		item.AddToCartUrl = "https://secure.newegg.com/Shopping/AddtoCart.aspx?Submit=ADD&ItemList=" + itemId
	}

	item.Config.Selectors = ParseSelectors(defaultConfig, item)
	if len(item.Config.Selectors) == 0 {
		item.Config.Selectors = []string{
			"#ProductBuy",
		}
	}

	defaultHeaders := make(map[string]string)
	item.Config.Headers = ParseHeaders(defaultConfig, defaultHeaders, item)

	return item
}

func (p *NeweggParser) Run(item types.Item) (string, string) {
	return scrapers.Run(item, func(text string) error {
		if strings.Contains(strings.ToLower(text), strings.ToLower("Are you a human?")) {
			return errors.New("Anti bot recaptcha")
		}
		return nil
	})
}
