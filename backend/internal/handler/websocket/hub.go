package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID    uuid.UUID
	Conn  *websocket.Conn
	Send  chan []byte
	Hub   *Hub
	Rooms map[string]bool
	mu    sync.RWMutex
}

type Hub struct {
	Clients    map[uuid.UUID]*Client
	Rooms      map[string]map[uuid.UUID]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	mu         sync.RWMutex
}

type Message struct {
	RoomID string
	Data   []byte
	FromID uuid.UUID
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[uuid.UUID]*Client),
		Rooms:      make(map[string]map[uuid.UUID]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			log.Printf("Client %s registered (total: %d)", client.ID, len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)

				client.mu.RLock()
				for roomID := range client.Rooms {
					if room, ok := h.Rooms[roomID]; ok {
						delete(room, client.ID)
						if len(room) == 0 {
							delete(h.Rooms, roomID)
						}
					}
				}
				client.mu.RUnlock()
			}
			h.mu.Unlock()
			log.Printf("Client %s unregistered (total: %d)", client.ID, len(h.Clients))

		case msg := <-h.Broadcast:
			h.mu.RLock()
			if room, ok := h.Rooms[msg.RoomID]; ok {
				for id, client := range room {
					if id != msg.FromID {
						select {
						case client.Send <- msg.Data:
						default:
							close(client.Send)
							delete(room, id)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToAll отправляет сообщение всем подключенным клиентам кроме excludeID
func (h *Hub) BroadcastToAll(data []byte, excludeID uuid.UUID) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for id, client := range h.Clients {
		if id != excludeID {
			select {
			case client.Send <- data:
				log.Printf("Sent to client: %s", id)
			default:
				log.Printf("Failed to send to client: %s", id)
			}
		}
	}
}

// SendToUser отправляет сообщение конкретному пользователю
func (h *Hub) SendToUser(userID uuid.UUID, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.Clients[userID]; ok {
		select {
		case client.Send <- data:
		default:
		}
	}
}

func (c *Client) JoinRoom(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Rooms == nil {
		c.Rooms = make(map[string]bool)
	}
	c.Rooms[roomID] = true

	c.Hub.mu.Lock()
	defer c.Hub.mu.Unlock()

	if c.Hub.Rooms[roomID] == nil {
		c.Hub.Rooms[roomID] = make(map[uuid.UUID]*Client)
	}
	c.Hub.Rooms[roomID][c.ID] = c
}

func (c *Client) LeaveRoom(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Rooms, roomID)

	c.Hub.mu.Lock()
	defer c.Hub.mu.Unlock()

	if room, ok := c.Hub.Rooms[roomID]; ok {
		delete(room, c.ID)
		if len(room) == 0 {
			delete(c.Hub.Rooms, roomID)
		}
	}
}

func (c *Client) SendSignal(signal interface{}) error {
	data, err := json.Marshal(signal)
	if err != nil {
		return err
	}

	select {
	case c.Send <- data:
	default:
	}
	return nil
}
