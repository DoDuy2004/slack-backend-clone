package websocket

import (
	"sync"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *WSMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Map of workspace_id/channel_id to clients in that "room"
	rooms map[string]map[*Client]bool

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *WSMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Clean up client from all rooms
				for roomID := range h.rooms {
					delete(h.rooms[roomID], client)
					if len(h.rooms[roomID]) == 0 {
						delete(h.rooms, roomID)
					}
				}
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleBroadcast(message *WSMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var targetClients map[*Client]bool

	if message.ChannelID != nil {
		// Broadcast to a specific channel
		roomID := "channel:" + message.ChannelID.String()
		targetClients = h.rooms[roomID]
	} else if message.WorkspaceID != nil {
		// Broadcast to an entire workspace
		roomID := "workspace:" + message.WorkspaceID.String()
		targetClients = h.rooms[roomID]
	} else if message.UserID != nil {
		// Private message/notification to a specific user
		// We would need a user_id to clients mapping for this ideally
		for client := range h.clients {
			if client.userID == *message.UserID {
				select {
				case client.send <- message:
				default:
					// If the send channel is full, we don't want to block the hub
				}
			}
		}
		return
	} else {
		// Global broadcast (rarely used)
		targetClients = h.clients
	}

	for client := range targetClients {
		select {
		case client.send <- message:
		default:
			// If buffer is full, we'll unregister the client (standard practice)
			// For now, just skip to avoid blocking hub
		}
	}
}

// JoinRoom adds a client to a specific room
func (h *Hub) JoinRoom(roomType string, id uuid.UUID, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomID := roomType + ":" + id.String()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
}

// LeaveRoom removes a client from a specific room
func (h *Hub) LeaveRoom(roomType string, id uuid.UUID, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	roomID := roomType + ":" + id.String()
	if h.rooms[roomID] != nil {
		delete(h.rooms[roomID], client)
		if len(h.rooms[roomID]) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

// Broadcast sends a message to the hub for distribution
func (h *Hub) Broadcast(message *WSMessage) {
	h.broadcast <- message
}
