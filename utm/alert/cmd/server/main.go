package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/utm/alert/internal/bridge"
	"github.com/utm/alert/internal/config"
	wsinternal "github.com/utm/alert/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	cfg := config.Load()

	br := bridge.New()
	feHub := wsinternal.NewFeHub(br.ToBackend)
	go feHub.Run()

	// goroutine forward: message từ backend → broadcast tất cả FE
	go func() {
		for msg := range br.ToFE {
			feHub.BroadcastChan() <- msg
		}
	}()

	mux := http.NewServeMux()

	// utm-backend kết nối vào đây (1 connection)
	mux.HandleFunc("/ws/backend", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("backend ws upgrade failed")
			return
		}
		log.Info().Msg("utm-backend connected")
		bc := wsinternal.NewBackendConn(conn, br.ToFE)

		// forward control command từ FE → backend
		go func() {
			for msg := range br.ToBackend {
				bc.SendChan() <- msg
			}
		}()

		bc.Serve() // block đến khi backend disconnect
		log.Warn().Msg("utm-backend disconnected")
	})

	// FE kết nối vào đây (N connections)
	mux.HandleFunc("/ws/monitor", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error().Err(err).Msg("fe ws upgrade failed")
			return
		}
		feHub.ServeClient(conn) // block đến khi FE disconnect
	})

	addr := ":" + cfg.Port
	log.Info().Str("addr", addr).Msg("utm-alert starting")
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal().Err(err).Msg("utm-alert server error")
	}
}
