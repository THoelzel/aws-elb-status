package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	//maxMessageSize = 512
)

var upgrader = websocket.Upgrader{}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	returnStatus(w)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Received websocket request",
		zap.String("remote_address", r.RemoteAddr),
	)

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		fmt.Println("403: " + r.Header.Get("Origin"))
		return
	}
	// websocket.Upgrader(w, r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Unable to upgrade connection",
			zap.Error(err),
		)
	}

	msgCh := b.Subscribe()

	go send(msgCh, conn)
	go receive(conn)
}

func send(msgCh chan interface{}, conn *websocket.Conn) {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case msg := <-msgCh:

			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteJSON(msg)

			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Warn("Unable to send through websocket, terminating connection",
						zap.String("remote_address", conn.RemoteAddr().String()),
						zap.Error(err),
					)
				} else {
					logger.Debug("Websocket connection terminated",
						zap.String("remote_address", conn.RemoteAddr().String()),
					)
				}
				b.Unsubscribe(msgCh)
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Debug("Unable to send websocket ping",
					zap.Error(err),
					zap.String("remote_address", conn.RemoteAddr().String()),
				)
				b.Unsubscribe(msgCh)
				return
			}
		}
	}
}

func receive(conn *websocket.Conn) {
	defer func() {
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Warn("Unexpected read error",
					zap.Error(err),
				)
			}
			break
		}
	}
}
