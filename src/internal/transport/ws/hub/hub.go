package hub

import (
	"fmt"
	"net/http"
	"suscord/internal/config"
	domainError "suscord/internal/domain/errors"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/storage"
	"suscord/internal/transport/ws/hub/dto"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

type Clients map[uint]*Client

type Rooms map[uint]map[uint]bool

type Hub struct {
	cfg        *config.Config
	rooms      Rooms
	clients    Clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan *dto.ResponseMessage
	mutex      *sync.RWMutex
	storage    storage.Storage
}

func NewHub(cfg *config.Config, storage storage.Storage, eventbus eventbus.Bus) *Hub {
	hub := &Hub{
		cfg:        cfg,
		rooms:      make(Rooms),
		clients:    make(Clients),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		broadcast:  make(chan *dto.ResponseMessage, 10),
		mutex:      &sync.RWMutex{},
		storage:    storage,
	}
	hub.RegisterEventSubscribers(eventbus)
	return hub
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.mutex.Lock()
			hub.clients[client.ID] = client
			hub.clients[client.ID].Rooms = make(map[uint]bool)
			hub.mutex.Unlock()

		case client := <-hub.unregister:
			hub.mutex.Lock()
			if _, exists := hub.clients[client.ID]; exists {
				// Сохраняем комнаты для очистки
				roomsToLeave := make([]uint, 0, len(client.Rooms))
				for roomID := range client.Rooms {
					roomsToLeave = append(roomsToLeave, roomID)
					// Удаляем клиента из комнаты
					if room, exists := hub.rooms[roomID]; exists {
						delete(room, client.ID)
						if len(room) == 0 {
							delete(hub.rooms, roomID)
						}
					}
				}
				delete(hub.clients, client.ID)
				client.Conn.Close()
			}
			hub.mutex.Unlock()

		case message := <-hub.broadcast:
			hub.broadcastToRoom(message.ChatID, message)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (hub *Hub) WebsocketHandler(c echo.Context) error {
	sessionUUID := c.QueryParam("session")
	if len(sessionUUID) == 0 {
		return c.NoContent(http.StatusForbidden)
	}

	session, err := hub.storage.Database().Session().GetByUUID(c.Request().Context(), sessionUUID)
	if err != nil {
		if pkgErrors.Is(err, domainError.ErrRecordNotFound) {
			return c.NoContent(http.StatusForbidden)
		}
		return err
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	user, err := hub.storage.Database().User().GetByID(c.Request().Context(), session.UserID)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	client := &Client{
		Conn:       conn,
		ID:         user.ID,
		Username:   user.Username,
		AvatarPath: user.AvatarPath,
		Rooms:      make(map[uint]bool),
	}

	hub.register <- client

	chats, err := hub.storage.Database().Chat().GetUserChats(c.Request().Context(), user.ID)
	if err != nil {
		conn.Close()
		return err
	}

	err = hub.joinToUserRooms(client, chats)
	if err != nil {
		if pkgErrors.Is(err, domainError.ErrUserIsNotMemberOfChat) {
			err = client.SendMessage(&dto.ResponseMessage{
				Type: "join_room_error",
				Data: map[string]interface{}{
					"message": "You are not member this room",
				},
			})
			if err != nil {
				return err
			}
		}
		return err
	}

	hub.receiveMessageHandler(conn, client)
	return nil
}

func (hub *Hub) receiveMessageHandler(conn *websocket.Conn, client *Client) {
	for {
		message := new(dto.ClientMessage)
		err := conn.ReadJSON(message)
		if err != nil {
			hub.unregister <- client
			return
		}

		err = hub.handleClientMessage(client, message)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}
