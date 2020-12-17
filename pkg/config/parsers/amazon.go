package parsers

import (
	"stock_scraper/pkg/scrapers"
	"stock_scraper/types"
)

type AmazonParser struct {
	label string
}

func (p *AmazonParser) Label() string {
	return p.label
}

func (p *AmazonParser) Parse(defaultConfig types.ItemConfig, item types.Item) types.Item {
	item.Config.Selectors = ParseSelectors(defaultConfig, item)
	if len(item.Config.Selectors) == 0 {
		item.Config.Selectors = []string{
			"span#price_inside_buybox",
		}
	}

	defaultHeaders := make(map[string]string)
	item.Config.Headers = ParseHeaders(defaultConfig, defaultHeaders, item)

	return item
}

func (p *AmazonParser) Run(item types.Item) (string, string) {
	return scrapers.Run(item, func(text string) error {
		return nil
	})
}
