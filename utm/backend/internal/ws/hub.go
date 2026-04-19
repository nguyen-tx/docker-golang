package ws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// ControlCommand là lệnh FE gửi lên BE qua WebSocket.
type ControlCommand struct {
	Type   string          `json:"type"`   // luôn là "control"
	Action string          `json:"action"` // "mute_alert" | "set_sensitivity" | "acknowledge"
	Params json.RawMessage `json:"params,omitempty"`
}

// CommandHandler xử lý lệnh nhận được từ FE.
type CommandHandler func(cmd ControlCommand)

// Client đại diện cho một frontend kết nối WebSocket.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub quản lý toàn bộ clients, dùng channel để tránh race condition.
type Hub struct {
	clients        map[*Client]bool
	broadcast      chan []byte
	register       chan *Client
	unregister     chan *Client
	commandHandler CommandHandler
}

func NewHub(cmdHandler CommandHandler) *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		broadcast:      make(chan []byte, 256),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		commandHandler: cmdHandler,
	}
}

// Run phải được chạy trong goroutine riêng.
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		}
	}
}

func (h *Hub) Broadcast(msg []byte) {
	h.broadcast <- msg
}

// ServeClient đăng ký client, bắt đầu write/read pump, block đến khi disconnect.
func (h *Hub) ServeClient(conn *websocket.Conn) {
	c := &Client{conn: conn, send: make(chan []byte, 256)}
	h.register <- c

	go c.writePump()
	c.readPump(h) // block cho đến khi FE disconnect
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// readPump đọc lệnh từ FE và chuyển sang commandHandler.
func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if h.commandHandler == nil {
			continue
		}
		var cmd ControlCommand
		if err := json.Unmarshal(raw, &cmd); err != nil {
			log.Warn().Err(err).Msg("invalid control command from FE")
			continue
		}
		if cmd.Type != "control" {
			continue
		}
		h.commandHandler(cmd)
	}
}