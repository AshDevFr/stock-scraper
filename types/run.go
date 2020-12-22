package types

type Parser interface {
	Label() string
	ParseId(Item) string
	Parse(ItemConfig, Item) Item
	Run(Item) (Result, string, error)
}

type Action struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Price struct {
	Symbol string
	Value  float64
}

type Result struct {
	Results []ParsedResults
	Content string
}

type ParsedResults struct {
	Price   *Price
	Result  map[string]string
	Content string
}
