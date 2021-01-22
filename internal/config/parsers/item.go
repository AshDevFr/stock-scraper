package parsers

import (
	"fmt"
	"github.com/imdario/mergo"
	"regexp"
	"stock_scraper/internal/utils"
	"stock_scraper/types"
)

var (
	defaultUserAgent      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"
	defaultCron           = "@every 2 minutes"
	defaultOpenLinks      = false
	defaultOpenAddToCart  = false
	defaultRunWeb         = false
	defaultWebRetries     = 0
	defaultInitWaitingSec = 0
	neweggRegex           = regexp.MustCompile(`(?i)newegg\.com`)
	bestbuyRegex          = regexp.MustCompile(`(?i)bestbuy\.com`)
	amazonRegex           = regexp.MustCompile(`(?i)amazon\.com`)
	defaultParser         = &DefaultParser{label: "default"}
	amazonParser          = &AmazonParser{label: "amazon"}
	bestbuyParser         = &BestBuyParser{label: "bestbuy"}
	neweggParser          = &NeweggParser{label: "newegg"}
)

func ParseItem(defaultConfig types.ItemConfig, item types.Item) types.Item {
	item.Uuid = genUuid(item)
	item.TrackedUrl = item.Url
	item.Parser = getParser(item)
	item.Id = ParseId(item)
	item.Type = ParseType(item)
	item.Config.Cron = ParseCron(defaultConfig, item)
	item.Config.UserAgent = ParseUserAgent(defaultConfig, item)
	item.Config.Rules = ParseRules(defaultConfig, item)
	item.Config.OpenLinks = ParseOpenLinks(defaultConfig, item)
	item.Config.OpenAddToCart = ParseOpenAddToCart(defaultConfig, item)
	item.Config.ItemSelector = ParseItemSelector(defaultConfig, item)
	item.Config.ItemLinkSelector = ParseItemLinkSelector(defaultConfig, item)
	item.Config.PriceSelector = ParsePriceSelector(defaultConfig, item)
	item.Config.MaxPrice = ParseMaxPrice(defaultConfig, item)
	item.Config.RunWeb = ParseRunWeb(defaultConfig, item)
	item.Config.WebRetries = ParseWebRetries(defaultConfig, item)
	item.Config.InitWaitingSec = ParseInitWaitingSec(defaultConfig, item)

	return item.Parser.Parse(defaultConfig, item)
}

func ParseId(item types.Item) string {
	if item.Id == "" {
		return item.Parser.ParseId(item)
	}
	return item.Type
}

func ParseType(item types.Item) string {
	if item.Type == "" {
		return item.Parser.Label()
	}
	return item.Type
}

func ParseItemSelector(defaultConfig types.ItemConfig, item types.Item) string {
	if item.Config.ItemSelector == "" {
		if defaultConfig.ItemSelector != "" {
			return defaultConfig.ItemSelector
		}
	}
	return item.Config.ItemSelector
}

func ParseItemLinkSelector(defaultConfig types.ItemConfig, item types.Item) string {
	if item.Config.ItemLinkSelector == "" {
		if defaultConfig.ItemLinkSelector != "" {
			return defaultConfig.ItemLinkSelector
		}
	}
	return item.Config.ItemLinkSelector
}

func ParsePriceSelector(defaultConfig types.ItemConfig, item types.Item) string {
	if item.Config.PriceSelector == "" {
		if defaultConfig.PriceSelector != "" {
			return defaultConfig.PriceSelector
		}
	}
	return item.Config.PriceSelector
}

func ParseMaxPrice(defaultConfig types.ItemConfig, item types.Item) *float64 {
	if item.Config.MaxPrice == nil {
		if defaultConfig.MaxPrice != nil {
			return defaultConfig.MaxPrice
		}
	}
	return item.Config.MaxPrice
}

func ParseCron(defaultConfig types.ItemConfig, item types.Item) string {
	if item.Config.Cron == "" {
		if defaultConfig.Cron != "" {
			return defaultConfig.Cron
		}
		return defaultCron
	}
	return item.Config.Cron
}

func ParseUserAgent(defaultConfig types.ItemConfig, item types.Item) string {
	if item.Config.UserAgent == "" {
		if defaultConfig.UserAgent != "" {
			return defaultConfig.UserAgent
		}
		return defaultUserAgent
	}
	return item.Config.UserAgent
}

func ParseOpenLinks(defaultConfig types.ItemConfig, item types.Item) *bool {
	if item.Config.OpenLinks != nil {
		return item.Config.OpenLinks
	} else if defaultConfig.OpenLinks != nil {
		return defaultConfig.OpenLinks
	}
	return &defaultOpenLinks
}

func ParseRunWeb(defaultConfig types.ItemConfig, item types.Item) *bool {
	if item.Config.RunWeb != nil {
		return item.Config.RunWeb
	} else if defaultConfig.RunWeb != nil {
		return defaultConfig.RunWeb
	}
	return &defaultRunWeb
}

func ParseWebRetries(defaultConfig types.ItemConfig, item types.Item) *int {
	if item.Config.WebRetries != nil {
		return item.Config.WebRetries
	} else if defaultConfig.WebRetries != nil {
		return defaultConfig.WebRetries
	}
	return &defaultWebRetries
}

func ParseInitWaitingSec(defaultConfig types.ItemConfig, item types.Item) *int {
	if item.Config.InitWaitingSec != nil {
		return item.Config.InitWaitingSec
	} else if defaultConfig.InitWaitingSec != nil {
		return defaultConfig.InitWaitingSec
	}
	return &defaultInitWaitingSec
}

func ParseOpenAddToCart(defaultConfig types.ItemConfig, item types.Item) *bool {
	if item.Config.OpenAddToCart != nil {
		return item.Config.OpenAddToCart
	} else if defaultConfig.OpenAddToCart != nil {
		return defaultConfig.OpenAddToCart
	}
	return &defaultOpenAddToCart
}

func ParseRules(defaultConfig types.ItemConfig, item types.Item) []types.Rule {
	if len(item.Config.Rules) == 0 {
		if len(defaultConfig.Rules) > 0 {
			return defaultConfig.Rules
		}
	}
	return item.Config.Rules
}

func ParseSelectors(defaultConfig types.ItemConfig, item types.Item) []string {
	if len(item.Config.Selectors) == 0 {
		if len(defaultConfig.Selectors) > 0 {
			return defaultConfig.Selectors
		}
	}
	return item.Config.Selectors
}

func ParseHeaders(defaultConfig types.ItemConfig, parserHeaders map[string]string, item types.Item) map[string]string {
	headers := make(map[string]string)
	mergo.Merge(&headers, defaultConfig.Headers)
	mergo.Merge(&headers, parserHeaders)
	mergo.Merge(&headers, item.Config.Headers)

	return headers
}

func genUuid(item types.Item) string {
	identifier := fmt.Sprintf("%s%s%s%s%v%v",
		item.Id,
		item.Url,
		item.Type,
		item.Config.UserAgent,
		item.Config.Selectors,
		item.Config.Rules,
	)
	return utils.Hash(identifier)
}

func getParser(item types.Item) types.Parser {
	var parser types.Parser
	switch {
	case item.Type == "newegg" || neweggRegex.MatchString(item.Url):
		parser = neweggParser
	case item.Type == "bestbuy" || bestbuyRegex.MatchString(item.Url):
		parser = bestbuyParser
	case item.Type == "amazon" || amazonRegex.MatchString(item.Url):
		parser = amazonParser
	default:
		parser = defaultParser
	}

	return parser
}
