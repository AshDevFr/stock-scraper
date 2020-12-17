package types

type ItemConfig struct {
	Cron          string            `json:"cron"`
	Headers       map[string]string `json:"headers"`
	Rules         []Rule            `json:"rules"`
	Selectors     []string          `json:"selectors"`
	UserAgent     string            `json:"userAgent"`
	OpenLinks     *bool             `json:"openLinks"`
	MaxPrice      *float64          `json:"maxPrice"`
	PriceSelector string            `json:"priceSelector"`
}

type Rule struct {
	Condition string   `json:"condition"` // text, added, removed, changed
	Strategy  string   `json:"strategy"`  // has, match
	Text      string   `json:"text"`
	Actions   []string `json:"actions"` // notify, open
}

type Action struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Price struct {
	Symbol string
	Value  float64
}

type Parser interface {
	Label() string
	Parse(ItemConfig, Item) Item
	Run(Item) (string, string, error)
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
	Concurrency *int `json:"concurrency"`
}
