package hub

import (
	"suscord/internal/transport/ws/hub/dto"
	"suscord/internal/transport/ws/hub/mapper"
)

func (hub *Hub) joinSFURoom(chatID uint, client *Client) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if len(hub.sfuRooms[chatID]) == 0 {
		hub.sfuRooms[chatID] = make(map[uint]bool)
	}

	hub.sfuRooms[chatID][client.ID] = true

	return nil
}

func (hub *Hub) leaveSFURoom(roomID uint, client *Client) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if _, exists := hub.sfuRooms[roomID]; exists {
		delete(hub.sfuRooms[roomID], client.ID)
	}

	if len(hub.sfuRooms[roomID]) == 0 {
		delete(hub.sfuRooms, roomID)
	}

	return nil
}

func (hub *Hub) isMemberOfSFURoom(roomID uint, clientID uint) bool {
	if _, ok := hub.sfuRooms[roomID]; ok {
		return hub.sfuRooms[roomID][clientID]
	}
	return false
}

func (hub *Hub) clientsSFURoom(chatID uint) ([]*dto.Client, error) {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()

	roomClients := hub.sfuRooms[chatID]

	clients := make([]*dto.Client, 0, len(roomClients))

	for clientID := range roomClients {
		if hubClient, ok := hub.clients[clientID]; ok {
			clients = append(clients, mapper.NewClient(&hubClient.Client, hub.cfg.Media.Url))
		}
	}

	return clients, nil
}
