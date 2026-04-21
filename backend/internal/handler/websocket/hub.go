package websocket

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

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
			log.Printf("Client %s registered", client.ID)

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
			log.Printf("Client %s unregistered", client.ID)

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
