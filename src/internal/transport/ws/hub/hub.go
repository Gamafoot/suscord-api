package hub

import (
	"fmt"
	"net/http"
	"suscord/internal/config"
	domainError "suscord/internal/domain/errors"
	"suscord/internal/domain/eventbus"
	"suscord/internal/domain/storage"
	"suscord/internal/transport/ws/hub/dto"
	"suscord/internal/transport/ws/hub/model"
	"sync"
	"time"

	"github.com/samber/lo"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	pkgErrors "github.com/pkg/errors"
)

type Clients map[uint]*Client

type ChatRooms map[uint]map[uint]bool
type SFURooms map[uint]map[uint]bool

type Hub struct {
	cfg        *config.Config
	chatRooms  ChatRooms
	sfuRooms   SFURooms
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
		chatRooms:  make(ChatRooms),
		sfuRooms:   make(SFURooms),
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
			hub.mutex.Unlock()

		case client := <-hub.unregister:
			affectedSFURooms := make([]uint, 0)

			hub.mutex.Lock()
			if _, exists := hub.clients[client.ID]; exists {
				// Сохраняем комнаты для очистки
				for roomID := range client.Rooms {
					// Удаляем клиента из комнаты
					if room, exists := hub.chatRooms[roomID]; exists {
						delete(room, client.ID)
						if len(room) == 0 {
							delete(hub.chatRooms, roomID)
						}
					}
				}

				// Удаляем клиента из SFU-комнат (на случай, если соединение оборвалось без call-leave)
				for roomID, room := range hub.sfuRooms {
					if _, ok := room[client.ID]; ok {
						delete(room, client.ID)
						affectedSFURooms = append(affectedSFURooms, roomID)

						if len(room) == 0 {
							delete(hub.sfuRooms, roomID)
						}
					}
				}

				delete(hub.clients, client.ID)
				client.Conn.Close()
			}
			hub.mutex.Unlock()

			// Уведомляем оставшихся участников SFU-комнат, чтобы они могли убрать аудио пользователя
			for _, roomID := range affectedSFURooms {
				clients, err := hub.clientsSFURoom(roomID)
				if err != nil {
					continue
				}

				hub.broadcastToSFURoom(roomID, &dto.ResponseMessage{
					Type: "call-leave",
					Data: map[string]any{
						"clients": clients,
						"user_id": client.ID,
					},
				})
			}

		case message := <-hub.broadcast:
			hub.broadcastToChatRoom(message.ChatID, message)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		cfg := config.GetConfig()
		origin := r.Header.Get("Origin")
		return lo.Contains(cfg.CORS.Origins, origin)
	},
}

func (hub *Hub) WebsocketHandler(c echo.Context) error {
	sessionUUID := c.QueryParam("session")
	if len(sessionUUID) == 0 {
		return c.NoContent(http.StatusForbidden)
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		fmt.Printf("ws error: %+v\n", pkgErrors.WithStack(err))
		return nil
	}

	session, err := hub.storage.Database().Session().GetByUUID(c.Request().Context(), sessionUUID)
	if err != nil {
		fmt.Printf("ws error: %+v\n", err)
		_ = conn.Close()
		if pkgErrors.Is(err, domainError.ErrRecordNotFound) {
			return c.NoContent(http.StatusForbidden)
		}
		return err
	}

	user, err := hub.storage.Database().User().GetByID(c.Request().Context(), session.UserID)
	if err != nil {
		fmt.Printf("ws error: %+v\n", err)
		conn.Close()
		return nil
	}

	conn.SetReadDeadline(time.Now().Add(hub.cfg.WebSocket.PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(hub.cfg.WebSocket.PongWait))
		return nil
	})

	go func() {
		ticker := time.NewTicker(hub.cfg.WebSocket.PingPeriod)
		defer ticker.Stop()

		for range ticker.C {
			err := conn.WriteControl(
				websocket.PingMessage,
				[]byte{},
				time.Now().Add(hub.cfg.WebSocket.Timeout),
			)
			if err != nil {
				conn.Close()
				return
			}
		}
	}()

	client := &Client{
		Conn: conn,
		Client: model.Client{
			ID:         user.ID,
			Username:   user.Username,
			AvatarPath: user.AvatarPath,
		},
		Rooms: make(map[uint]bool),
	}

	chats, err := hub.storage.Database().Chat().GetUserChats(c.Request().Context(), user.ID)
	if err != nil {
		fmt.Printf("ws error: %+v\n", pkgErrors.WithStack(err))
		conn.Close()
		return nil
	}

	err = hub.joinToUserChatRooms(client, chats)
	if err != nil {
		if pkgErrors.Is(err, domainError.ErrUserIsNotMemberOfChat) {
			sendErr := client.SendMessage(&dto.ResponseMessage{
				Type: "join_room_error",
				Data: map[string]interface{}{
					"message": "You are not member this room",
				},
			})
			if sendErr != nil {
				fmt.Printf("ws error: %+v\n", pkgErrors.WithStack(sendErr))
			}
			conn.Close()
			return nil
		}
		fmt.Printf("ws error: %+v\n", pkgErrors.WithStack(err))
		conn.Close()
		return nil
	}

	hub.register <- client
	hub.receiveMessageHandler(conn, client)
	return nil
}

func (hub *Hub) receiveMessageHandler(conn *websocket.Conn, client *Client) {
	for {
		message := new(dto.ClientMessage)
		err := conn.ReadJSON(message)
		if err != nil {
			fmt.Printf("ws readJSON: %v\n", pkgErrors.WithStack(err))
			hub.unregister <- client
			return
		}

		err = hub.handleClientMessage(client, message)
		if err != nil {
			fmt.Printf("ws handleClientMessage: %+v\n", err)
		}
	}
}
