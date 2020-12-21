package websocket

import (
	"encoding/json"
	"stock_scraper/types"
	"time"
)

func send(payload interface{}) {
	hub := GetHub()
	message, err := json.Marshal(&payload)
	if err == nil {
		hub.Broadcast(message)
	}
}

func SendUpdateMessage(scraperType string, item types.Item, status string, message string) {
	payload := types.WSUpdatePayload{
		Time:    time.Now().Unix(),
		Type:    "update",
		Scraper: scraperType,
		Item:    item,
		Status:  status,
		Message: message,
	}
	send(payload)
}

func SendActionMessage(action types.Action, item types.Item) {
	payload := types.WSActionPayload{
		Time:    time.Now().Unix(),
		Type:    "action",
		Item:    item,
		Action:  action.Type,
		Content: action.Content,
	}
	send(payload)
}
