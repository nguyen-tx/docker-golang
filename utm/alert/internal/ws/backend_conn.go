package ws

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// BackendConn quản lý 1 kết nối WebSocket từ utm-backend.
type BackendConn struct {
	conn       *websocket.Conn
	send       chan []byte // nhận message từ bridge (control cmd từ FE) → gửi xuống backend
	toFE       chan []byte // message từ backend → đẩy lên broadcast của FeHub
}

func NewBackendConn(conn *websocket.Conn, toFE chan []byte) *BackendConn {
	return &BackendConn{
		conn: conn,
		send: make(chan []byte, 256),
		toFE: toFE,
	}
}

func (b *BackendConn) SendChan() chan<- []byte { return b.send }

// Serve bắt đầu write/read pump, block đến khi backend disconnect.
func (b *BackendConn) Serve() {
	go b.writePump()
	b.readPump() // block
}

func (b *BackendConn) writePump() {
	defer b.conn.Close()
	for msg := range b.send {
		if err := b.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Warn().Err(err).Msg("write to backend failed")
			return
		}
	}
}

// readPump nhận telemetry/alert từ backend → đẩy vào toFE để broadcast tới FE.
func (b *BackendConn) readPump() {
	defer b.conn.Close()
	for {
		_, raw, err := b.conn.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("backend disconnected")
			break
		}
		select {
		case b.toFE <- raw:
		default:
			log.Warn().Msg("toFE channel full, dropping message")
		}
	}
}
