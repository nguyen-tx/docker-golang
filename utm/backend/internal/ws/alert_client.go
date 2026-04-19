package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// AlertClient là WS client kết nối từ utm-backend đến utm-alert.
// Gửi telemetry/alert lên, nhận control command về.
type AlertClient struct {
	alertURL       string
	send           chan []byte
	commandHandler CommandHandler
	conn           *websocket.Conn
}

type CommandHandler func(raw []byte)

func NewAlertClient(alertURL string, cmdHandler CommandHandler) *AlertClient {
	return &AlertClient{
		alertURL:       alertURL,
		send:           make(chan []byte, 256),
		commandHandler: cmdHandler,
	}
}

// Run kết nối đến utm-alert và giữ kết nối sống, tự reconnect khi mất.
// Phải chạy trong goroutine riêng.
func (c *AlertClient) Run() {
	for {
		if err := c.connect(); err != nil {
			log.Warn().Err(err).Str("url", c.alertURL).Msg("utm-alert connect failed, retrying in 5s")
			time.Sleep(5 * time.Second)
			continue
		}
		// connect() block đến khi mất kết nối, sau đó retry
		log.Warn().Msg("utm-alert disconnected, reconnecting...")
		time.Sleep(2 * time.Second)
	}
}

func (c *AlertClient) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.alertURL, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	log.Info().Str("url", c.alertURL).Msg("connected to utm-alert")

	done := make(chan struct{})
	go c.writePump(conn, done)
	c.readPump(conn, done) // block
	return nil
}

// Send đẩy message vào queue để writePump gửi lên utm-alert.
func (c *AlertClient) Send(msg []byte) {
	select {
	case c.send <- msg:
	default:
		log.Warn().Msg("alert client send buffer full, dropping message")
	}
}

func (c *AlertClient) writePump(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for msg := range c.send {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// readPump nhận control command từ utm-alert (do FE gửi lên).
func (c *AlertClient) readPump(conn *websocket.Conn, done chan struct{}) {
	defer conn.Close()
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if c.commandHandler != nil {
			c.commandHandler(raw)
		}
	}
}
