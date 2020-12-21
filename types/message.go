package types

type WSUpdatePayload struct {
	Time         int64
	Type         string `json:"type"`
	Scraper      string `json:"scraper"`
	Item         Item   `json:"item"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	AddToCartUrl string `json:"addToCartUrl"`
}

type WSActionPayload struct {
	Time    int64
	Type    string `json:"type"`
	Action  string `json:"action"`
	Item    Item   `json:"item"`
	Content string `json:"content"`
}
