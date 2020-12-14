package scrapers

import (
	"stock_scraper/types"
)

func RunAmazon(item types.Item) (string, string) {
	defaultAmazonSelectors := []string{
		"span#price_inside_buybox",
	}
	return Run(item, defaultAmazonSelectors, make(map[string]string))
}
