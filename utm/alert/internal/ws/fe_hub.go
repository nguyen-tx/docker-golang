package ws

import "github.com/gorilla/websocket"

// FeHub quản lý N kết nối WebSocket từ FE clients.
type FeHub struct {
	clients    map[*feClient]bool
	broadcast  chan []byte  // nhận message từ bridge → gửi tới tất cả FE
	toBackend  chan []byte  // nhận control command từ FE → gửi lên backend
	register   chan *feClient
	unregister chan *feClient
}

type feClient struct {
	conn *websocket.Conn
	send chan []byte
}

func NewFeHub(toBackend chan []byte) *FeHub {
	return &FeHub{
		clients:    make(map[*feClient]bool),
		broadcast:  make(chan []byte, 256),
		toBackend:  toBackend,
		register:   make(chan *feClient),
		unregister: make(chan *feClient),
	}
}

func (h *FeHub) BroadcastChan() chan<- []byte { return h.broadcast }

// Run phải chạy trong goroutine riêng.
func (h *FeHub) Run() {
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

// ServeClient đăng ký FE client, block đến khi disconnect.
func (h *FeHub) ServeClient(conn *websocket.Conn) {
	c := &feClient{conn: conn, send: make(chan []byte, 256)}
	h.register <- c
	go c.writePump()
	c.readPump(h) // block
}

func (c *feClient) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// readPump đọc control command từ FE, đẩy vào toBackend channel.
func (c *feClient) readPump(h *FeHub) {
	defer func() {
		h.unregister <- c
		c.conn.Close()
	}()
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// forward nguyên xi lên backend, không parse
		select {
		case h.toBackend <- raw:
		default:
		}
	}
}
