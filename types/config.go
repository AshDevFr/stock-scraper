package types

type Rule struct {
	Condition string `json:"condition"` // text, added, removed
	Strategy  string `json:"strategy"`  // has, match
	Text      string `json:"text"`
}

type Action struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Item struct {
	Id        string   `json:"id"`
	Url       string   `json:"url"`
	Type      string   `json:"type"`
	Selectors []string `json:"selectors"`
	UserAgent string   `json:"userAgent"`
	Rules     []Rule   `json:"rules"`
}

type Config struct {
	Items []Item `json:"items"`
	Rules []Rule `json:"rules"`
}
