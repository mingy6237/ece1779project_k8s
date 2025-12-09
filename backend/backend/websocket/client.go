package websocket

import (
	"log"
	"net/http"
	"time"

	"inventory-manager-server/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development and demo purposes
		// In production, restrict to specific allowed origins
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024
)

// Client represents a WebSocket client connection
type Client struct {
	Hub    *HubType
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
	Role   string
}

// NewClient creates a new client
func NewClient(conn *websocket.Conn, userID string, role string) *Client {
	return &Client{
		Hub:    Hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
		Role:   role,
	}
}

// ReadPump reads messages from WebSocket connection
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Echo message back to all clients (broadcast)
		// In a real application, you might want to process the message first
		c.Hub.Broadcast(message)
	}
}

// WritePump writes messages to WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS handles WebSocket upgrade request
func ServeWS(w http.ResponseWriter, r *http.Request) {
	// Extract token from query parameter
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Unauthorized: token parameter required", http.StatusUnauthorized)
		return
	}

	// Validate token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		http.Error(w, "Unauthorized: invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Hub check
	if Hub == nil {
		log.Printf("WebSocket Hub not initialized")
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create client with user info from claims
	client := NewClient(conn, claims.UserID.String(), claims.Role)
	Hub.Register(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}
