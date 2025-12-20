package dto

import "encoding/json"

type ClientMessage struct {
	Type   string          `json:"type"`
	ChatID uint            `json:"chat_id,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

type ResponseMessage struct {
	ChatID uint   `json:"-"`
	Type   string `json:"type"`
	Data   any    `json:"data"`
}

type SendMessage struct {
	RoomID  uint `json:"room_id"`
	Content uint `json:"content"`
}
