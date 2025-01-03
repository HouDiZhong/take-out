package websocket

import (
	"log/slog"
	"net/http"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHub struct {
	Hub *service.Hub
}

func NewWebSocketHub(hub *service.Hub) *WebSocketHub {
	return &WebSocketHub{
		Hub: hub,
	}
}

func (h *WebSocketHub) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("WebSocket upgrade error:", "error", err)
		return
	}
	defer conn.Close()

	userID := c.Param("id")
	if userID == "" {
		slog.Error("WebSocket user ID is empty")
		return
	}
	client := &service.Client{
		Send: make(chan []byte, 256),
		ID:   userID,
	}
	// 将客户端与用户 ID 关联
	h.Hub.Mutex.Lock()
	h.Hub.Clients[client] = true
	h.Hub.ClientsByID[userID] = client
	h.Hub.Mutex.Unlock()

	go h.writePump(conn, client)
	h.readPump(conn, client)
}

func (h *WebSocketHub) readPump(conn *websocket.Conn, client *service.Client) {
	defer func() {
		h.Hub.Unregister <- client
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("WebSocket error:", "error", err)
			}
			break
		}
		h.Hub.Broadcast <- message
	}
}

func (h *WebSocketHub) writePump(conn *websocket.Conn, client *service.Client) {
	defer conn.Close()

	for message := range client.Send {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			slog.Error("WebSocket write error:", "error", err)
			return
		}
	}
	conn.WriteMessage(websocket.CloseMessage, []byte{})
}
