package dto

import "encoding/json"

type ClientMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ResponseMessage struct {
	ChatID uint        `json:"-"`
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
}

type SendMessage struct {
	RoomID  uint `json:"room_id"`
	Content uint `json:"content"`
}
