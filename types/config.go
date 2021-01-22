package types

type ItemConfig struct {
	Cron             string            `json:"cron"`
	Headers          map[string]string `json:"headers"`
	Rules            []Rule            `json:"rules"`
	ItemSelector     string            `json:"itemSelector"`
	ItemLinkSelector string            `json:"itemLinkSelector"`
	Selectors        []string          `json:"selectors"`
	UserAgent        string            `json:"userAgent"`
	OpenLinks        *bool             `json:"openLinks"`
	OpenAddToCart    *bool             `json:"openAddToCart"`
	MaxPrice         *float64          `json:"maxPrice"`
	PriceSelector    string            `json:"priceSelector"`
	RunWeb           *bool             `json:"runWeb"`
	InitWaitingSec   *int              `json:"initWaitingSec"`
	NoHeadless       bool              `json:"noHeadless"`
	WebRetries       *int              `json:"webRetries"`
	ForceUrl         bool              `json:"forceUrl"`
}

type Rule struct {
	Condition string   `json:"condition"` // text, added, removed, changed
	Strategy  string   `json:"strategy"`  // has, match
	Text      string   `json:"text"`
	Actions   []string `json:"actions"` // notify, open
}

type Item struct {
	Id           string     `json:"id"`
	Uuid         string     `json:"uuid"`
	Title        string     `json:"title"`
	Url          string     `json:"url"`
	TrackedUrl   string     `json:"trackedUrl"`
	AddToCartUrl string     `json:"addToCartUrl"`
	Type         string     `json:"type"`
	Config       ItemConfig `json:"config"`
	Parser       Parser
}

type Config struct {
	Items         []Item     `json:"items"`
	DefaultConfig ItemConfig `json:"defaultConfig"`
	ItemMap       map[string]*Item
	Concurrency   *int `json:"concurrency"`
}
