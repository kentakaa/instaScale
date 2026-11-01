package hub

import (
	"encoding/json"
	"log"
	"sync"
)

type Hub struct {
	// Registered clients indexed by UserID (a single user might have multiple active tabs)
	userClients map[string]map[*Client]bool

	// Broadcast channel for global announcements or feed events
	Broadcast chan []byte

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		userClients: make(map[string]map[*Client]bool),
		Broadcast:   make(chan []byte, 1024),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, exists := h.userClients[client.UserID]; !exists {
				h.userClients[client.UserID] = make(map[*Client]bool)
			}
			h.userClients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("[Hub] User connected: %s (Total online sessions: %d)", client.UserID, h.ActiveSessionsCount())

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, exists := h.userClients[client.UserID]; exists {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.userClients, client.UserID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("[Hub] User disconnected: %s", client.UserID)

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, clients := range h.userClients {
				for client := range clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// RouteMessage sends a message directly to a target recipient or broadcasts it
func (h *Hub) RouteMessage(payload MessagePayload) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return
	}

	if payload.RecipientID != "" {
		h.mu.RLock()
		defer h.mu.RUnlock()

		// Deliver to recipient's active devices
		if recipientClients, ok := h.userClients[payload.RecipientID]; ok {
			for client := range recipientClients {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(recipientClients, client)
				}
			}
		}

		// Also echo back to sender's other open tabs
		if senderClients, ok := h.userClients[payload.SenderID]; ok {
			for client := range senderClients {
				select {
				case client.Send <- data:
				default:
				}
			}
		}
	} else {
		// Broadcast to all
		h.Broadcast <- data
	}
}

// DirectNotify sends a push event directly to a user
func (h *Hub) DirectNotify(userID string, payload MessagePayload) bool {
	data, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.userClients[userID]; ok && len(clients) > 0 {
		for client := range clients {
			select {
			case client.Send <- data:
			default:
			}
		}
		return true
	}
	return false
}

func (h *Hub) ActiveSessionsCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, clients := range h.userClients {
		total += len(clients)
	}
	return total
}

func (h *Hub) ActiveUsersCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userClients)
}
