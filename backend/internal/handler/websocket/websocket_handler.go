package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"backend/internal/domain"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	ID    uuid.UUID
	Conn  *websocket.Conn
	Send  chan []byte
	Hub   *Hub
	Rooms map[string]bool
	mu    sync.RWMutex
}

type WebSocketHandler struct {
	hub         *Hub
	callService service.CallService
	userService service.UserService
}

func NewWebSocketHandler(hub *Hub, callService service.CallService, userService service.UserService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		callService: callService,
		userService: userService,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	userIDStr := c.Query("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:    userID,
		Conn:  conn,
		Send:  make(chan []byte, 256),
		Hub:   h.hub,
		Rooms: make(map[string]bool),
	}

	h.hub.Register <- client

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WebSocketHandler) readPump(client *Client) {
	defer func() {
		h.hub.Unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var signal domain.CallSignal
		if err := json.Unmarshal(message, &signal); err != nil {
			continue
		}

		signal.FromID = client.ID
		h.handleSignal(client, &signal)
	}
}

func (h *WebSocketHandler) writePump(client *Client) {
	defer client.Conn.Close()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) handleSignal(client *Client, signal *domain.CallSignal) {
	switch signal.Type {
	case "call-start":
		h.handleCallStart(client, signal)
	case "call-accept":
		h.handleCallAccept(client, signal)
	case "call-reject":
		h.handleCallReject(client, signal)
	case "call-end":
		h.handleCallEnd(client, signal)
	case "offer", "answer", "ice-candidate":
		h.handleWebRTCSignal(client, signal)
	}
}

func (h *WebSocketHandler) handleCallStart(client *Client, signal *domain.CallSignal) {
	chatID, _ := uuid.Parse(signal.ChatID)
	calleeID, _ := uuid.Parse(signal.CalleeID)
	callType := domain.CallType(signal.CallType)

	call, err := h.callService.StartCall(context.Background(), chatID, client.ID, calleeID, callType)
	if err != nil {
		log.Printf("Failed to start call: %v", err)
		return
	}

	h.joinRoom(client, call.RoomID)

	signal.RoomID = call.RoomID
	signal.Call = &domain.CallResponse{
		ID:       call.ID,
		ChatID:   call.ChatID,
		CallerID: call.CallerID,
		CalleeID: call.CalleeID,
		Type:     call.Type,
		Status:   call.Status,
		RoomID:   call.RoomID,
	}

	caller, _ := h.userService.GetByID(context.Background(), call.CallerID)
	if caller != nil {
		signal.Call.Caller = &domain.UserInfo{
			ID:        caller.ID,
			Username:  caller.Username,
			AvatarURL: caller.AvatarURL,
		}
	}

	data, _ := json.Marshal(signal)
	h.hub.mu.RLock()
	if callee, ok := h.hub.Clients[call.CalleeID]; ok {
		callee.Send <- data
	}
	h.hub.mu.RUnlock()
}

func (h *WebSocketHandler) handleCallAccept(client *Client, signal *domain.CallSignal) {
	call, err := h.callService.AcceptCall(context.Background(), signal.CallID)
	if err != nil {
		return
	}

	h.joinRoom(client, call.RoomID)

	signal.RoomID = call.RoomID
	data, _ := json.Marshal(signal)

	h.hub.mu.RLock()
	if caller, ok := h.hub.Clients[call.CallerID]; ok {
		caller.Send <- data
	}
	h.hub.mu.RUnlock()
}

func (h *WebSocketHandler) handleCallReject(client *Client, signal *domain.CallSignal) {
	call, err := h.callService.RejectCall(context.Background(), signal.CallID)
	if err != nil {
		return
	}

	signal.RoomID = call.RoomID
	data, _ := json.Marshal(signal)

	h.hub.mu.RLock()
	if caller, ok := h.hub.Clients[call.CallerID]; ok {
		caller.Send <- data
	}
	h.hub.mu.RUnlock()
}

func (h *WebSocketHandler) handleCallEnd(client *Client, signal *domain.CallSignal) {
	call, err := h.callService.EndCall(context.Background(), signal.CallID)
	if err != nil {
		return
	}

	h.leaveRoom(client, call.RoomID)

	signal.RoomID = call.RoomID
	data, _ := json.Marshal(signal)
	h.hub.Broadcast <- &Message{
		RoomID: call.RoomID,
		Data:   data,
		FromID: client.ID,
	}
}

func (h *WebSocketHandler) handleWebRTCSignal(client *Client, signal *domain.CallSignal) {
	h.joinRoom(client, signal.RoomID)

	data, _ := json.Marshal(signal)
	h.hub.Broadcast <- &Message{
		RoomID: signal.RoomID,
		Data:   data,
		FromID: client.ID,
	}
}

func (h *WebSocketHandler) joinRoom(client *Client, roomID string) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.Rooms == nil {
		client.Rooms = make(map[string]bool)
	}
	client.Rooms[roomID] = true

	h.hub.mu.Lock()
	defer h.hub.mu.Unlock()

	if h.hub.Rooms[roomID] == nil {
		h.hub.Rooms[roomID] = make(map[uuid.UUID]*Client)
	}
	h.hub.Rooms[roomID][client.ID] = client
}

func (h *WebSocketHandler) leaveRoom(client *Client, roomID string) {
	client.mu.Lock()
	defer client.mu.Unlock()

	delete(client.Rooms, roomID)

	h.hub.mu.Lock()
	defer h.hub.mu.Unlock()

	if room, ok := h.hub.Rooms[roomID]; ok {
		delete(room, client.ID)
		if len(room) == 0 {
			delete(h.hub.Rooms, roomID)
		}
	}
}
