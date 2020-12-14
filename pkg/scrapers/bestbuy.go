package scrapers

import (
	"stock_scraper/types"
)

func RunBestBuy(item types.Item) (string, string) {
	defaultBestBuySelectors := []string{
		"div.row.v-m-bottom-g div.col-xs-12  div.fulfillment-add-to-cart-button  button.add-to-cart-button",
	}
	extraHeaders := make(map[string]string)
	//extraHeaders["accept-encoding"] = "gzip, deflate"
	//extraHeaders["cache-control"] = "no-cache"
	extraHeaders["accept-language"] = "en-US,en;q=0.9"
	//extraHeaders["pragma"] = "no-cache"
	//extraHeaders["sec-fetch-dest"] = "document"
	//extraHeaders["sec-fetch-mode"] = "navigate"
	//extraHeaders["sec-fetch-site"] = "none"
	//extraHeaders["upgrade-insecure-requests"] = "1"
	return Run(item, defaultBestBuySelectors, extraHeaders)
}
