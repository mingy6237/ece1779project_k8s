package websocket

import (
	"log"
)

// Hub maintains all active client connections
var Hub *HubType

// HubType is the type for Hub
type HubType struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

// InitHub initializes the WebSocket Hub
func InitHub() error {
	Hub = &HubType{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go Hub.Run()
	log.Println("WebSocket Hub initialized and started")
	return nil
}

// GetHub returns the global Hub instance
func GetHub() *HubType {
	return Hub
}

// Run runs the Hub
func (h *HubType) Run() {
	log.Println("WebSocket Hub started")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client registered. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("Client unregistered. Total clients: %d", len(h.clients))
			}

		case message := <-h.broadcast:
			// Broadcast message to all registered clients
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// If client's send channel is full, close connection
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Register registers a client
func (h *HubType) Register(client *Client) {
	h.register <- client
}

// Unregister unregisters a client
func (h *HubType) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast broadcasts a message
func (h *HubType) Broadcast(message []byte) {
	h.broadcast <- message
}
