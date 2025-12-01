package dto

import "encoding/json"

type ClientMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ResponseMessage struct {
	RoomID uint        `json:"-"`
	Type   string      `json:"type"`
	Data   interface{} `json:"data"`
}

type WsGetRoomClients struct {
	RoomID uint `json:"room_id"`
}

type SendMessage struct {
	RoomID  uint `json:"room_id"`
	Content uint `json:"content"`
}
