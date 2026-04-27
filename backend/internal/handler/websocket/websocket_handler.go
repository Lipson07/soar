package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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

	h.userService.UpdateStatus(context.Background(), userID, "online")

	statusData, _ := json.Marshal(map[string]interface{}{
		"type":    "user-status",
		"user_id": userID.String(),
		"status":  "online",
	})
	h.hub.BroadcastToAll(statusData, userID)

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WebSocketHandler) readPump(client *Client) {
	defer func() {
		h.hub.Unregister <- client
		client.Conn.Close()

		h.userService.UpdateStatus(context.Background(), client.ID, "offline")

		statusData, _ := json.Marshal(map[string]interface{}{
			"type":    "user-status",
			"user_id": client.ID.String(),
			"status":  "offline",
		})
		h.hub.BroadcastToAll(statusData, client.ID)
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("Received from %s: %s", client.ID, string(message))

		var baseSignal struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(message, &baseSignal); err != nil {
			log.Printf("Failed to parse message type: %v", err)
			continue
		}

		switch baseSignal.Type {
		case "user-status":
			var statusSignal struct {
				Type   string `json:"type"`
				UserID string `json:"user_id"`
				Status string `json:"status"`
			}
			if err := json.Unmarshal(message, &statusSignal); err != nil {
				log.Printf("Failed to parse user-status: %v", err)
				continue
			}
			h.handleStatusSignal(client, &statusSignal)

		case "typing":
			var typingSignal struct {
				Type     string `json:"type"`
				UserID   string `json:"user_id"`
				ChatID   string `json:"chat_id"`
				IsTyping bool   `json:"is_typing"`
				Username string `json:"username"`
			}
			if err := json.Unmarshal(message, &typingSignal); err != nil {
				log.Printf("Failed to parse typing: %v", err)
				continue
			}
			log.Printf("Typing from %s: %+v", client.ID, typingSignal)
			h.handleTypingSignal(client, &typingSignal)

		case "call-start", "call-accept", "call-reject", "call-end", "offer", "answer", "ice-candidate":
			var signal domain.CallSignal
			if err := json.Unmarshal(message, &signal); err != nil {
				log.Printf("Failed to parse call signal: %v", err)
				continue
			}
			signal.FromID = client.ID
			h.handleCallSignal(client, &signal)

		default:
			log.Printf("Unknown message type: %s", baseSignal.Type)
		}
	}
}

func (h *WebSocketHandler) writePump(client *Client) {
	defer client.Conn.Close()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) handleStatusSignal(client *Client, signal *struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
	Status string `json:"status"`
}) {
	userID, _ := uuid.Parse(signal.UserID)
	h.userService.UpdateStatus(context.Background(), userID, signal.Status)

	data, _ := json.Marshal(map[string]interface{}{
		"type":    "user-status",
		"user_id": signal.UserID,
		"status":  signal.Status,
	})
	h.hub.BroadcastToAll(data, userID)
}

func (h *WebSocketHandler) handleTypingSignal(client *Client, signal *struct {
	Type     string `json:"type"`
	UserID   string `json:"user_id"`
	ChatID   string `json:"chat_id"`
	IsTyping bool   `json:"is_typing"`
	Username string `json:"username"`
}) {
	data, _ := json.Marshal(map[string]interface{}{
		"type":      "typing",
		"user_id":   signal.UserID,
		"chat_id":   signal.ChatID,
		"is_typing": signal.IsTyping,
		"username":  signal.Username,
	})

	log.Printf("Broadcasting typing: %s", string(data))
	h.hub.BroadcastToAll(data, client.ID)
}

func (h *WebSocketHandler) handleCallSignal(client *Client, signal *domain.CallSignal) {
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

	client.JoinRoom(call.RoomID)

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
	h.hub.SendToUser(call.CalleeID, data)
}

func (h *WebSocketHandler) handleCallAccept(client *Client, signal *domain.CallSignal) {
	call, err := h.callService.AcceptCall(context.Background(), signal.CallID)
	if err != nil {
		return
	}

	client.JoinRoom(call.RoomID)

	signal.RoomID = call.RoomID
	data, _ := json.Marshal(signal)
	h.hub.SendToUser(call.CallerID, data)
}

func (h *WebSocketHandler) handleCallReject(client *Client, signal *domain.CallSignal) {
	call, err := h.callService.RejectCall(context.Background(), signal.CallID)
	if err != nil {
		return
	}

	data, _ := json.Marshal(signal)
	h.hub.SendToUser(call.CallerID, data)
}

func (h *WebSocketHandler) handleCallEnd(client *Client, signal *domain.CallSignal) {
	call, err := h.callService.EndCall(context.Background(), signal.CallID)
	if err != nil {
		return
	}

	client.LeaveRoom(call.RoomID)

	data, _ := json.Marshal(signal)
	h.hub.Broadcast <- &Message{
		RoomID: call.RoomID,
		Data:   data,
		FromID: client.ID,
	}
}

func (h *WebSocketHandler) handleWebRTCSignal(client *Client, signal *domain.CallSignal) {
	client.JoinRoom(signal.RoomID)

	data, _ := json.Marshal(signal)
	h.hub.Broadcast <- &Message{
		RoomID: signal.RoomID,
		Data:   data,
		FromID: client.ID,
	}
}
