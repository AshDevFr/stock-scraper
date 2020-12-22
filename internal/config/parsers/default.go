package parsers

import (
	"stock_scraper/internal/scrapers"
	"stock_scraper/types"
)

type DefaultParser struct {
	label string
}

func (p *DefaultParser) ParseId(item types.Item) string {
	if item.Id != "" {
		return item.Id
	}

	return ""
}

func (p *DefaultParser) Label() string {
	return p.label
}

func (p *DefaultParser) Parse(defaultConfig types.ItemConfig, item types.Item) types.Item {
	item.Config.Selectors = ParseSelectors(defaultConfig, item)
	if len(item.Config.Selectors) == 0 {
		item.Config.Selectors = []string{
			"body",
		}
	}

	defaultHeaders := make(map[string]string)
	item.Config.Headers = ParseHeaders(defaultConfig, defaultHeaders, item)

	return item
}

func (p *DefaultParser) Run(item types.Item) (types.Result, string, error) {
	return scrapers.Run(item, func(body string, results []types.ParsedResults) (string, error) {
		return "", nil
	})
}
