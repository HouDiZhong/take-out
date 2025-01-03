package service

import (
	"log/slog"
	"sync"
)

type Client struct {
	Send chan []byte // 发送消息的通道
	ID   string      // 用户 ID
}

type Hub struct {
	Clients     map[*Client]bool   // 所有客户端
	ClientsByID map[string]*Client // 根据用户 ID 管理客户端
	Register    chan *Client       // 注册客户端
	Unregister  chan *Client       // 注销客户端
	Broadcast   chan []byte        // 广播消息
	Mutex       sync.Mutex         // 互斥锁
}

func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		ClientsByID: make(map[string]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan []byte),
		Mutex:       sync.Mutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client] = true
			h.Mutex.Unlock()
		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			h.Mutex.Lock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.Mutex.Unlock()
		}
	}
}

// 向特定客户端发送消息
func (h *Hub) SendToClient(userID string, message []byte) {
	h.Mutex.Lock()
	client, ok := h.ClientsByID[userID]
	h.Mutex.Unlock()

	if !ok {
		slog.Error("Client not found:", "error", userID)
		return
	}

	client.Send <- message
}

// 向所有客户端广播消息
func (h *Hub) BroadcastMessage(message []byte) {
	h.Mutex.Lock()
	for client := range h.Clients {
		client.Send <- message
	}
	h.Mutex.Unlock()
}

type NoticeService struct {
	Hub *Hub
}

func NewNoticeService(hub *Hub) *NoticeService {
	return &NoticeService{
		Hub: hub,
	}
}

// 向特定用户发送消息
func (s *NoticeService) SendMessageToUser(userID string, message string) {
	s.Hub.SendToClient(userID, []byte(message))
}

// 向所有用户广播消息
func (s *NoticeService) BroadcastMessage(message string) {
	s.Hub.BroadcastMessage([]byte(message))
}
