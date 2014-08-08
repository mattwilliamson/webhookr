// Copyright 2013 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wh

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	wsWriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	wsPongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than wsPongWait.
	wsPingPeriod = (wsPongWait * 9) / 10

	// Maximum message size allowed from peer.
	wsMaxMessageSize = 512
)

// WSClient is a middleman between the websocket connection and the hub.
type WSClient struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Webhookr ID
	whid string
}

// readPump pumps messages from the websocket connection to the hub.
func (c *WSClient) readPump(h *ClientRegistry) {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(wsMaxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(wsPongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(wsPongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		m := ClientMessage{c.whid, message}
		h.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *WSClient) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(wsWriteWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *WSClient) writePump() {
	ticker := time.NewTicker(wsPingPeriod)
	defer ticker.Stop()
	defer c.ws.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// websocketHandler handles webocket requests from the peer.
func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	// if r.Header.Get("Origin") != "http://"+r.Host {
	// 	http.Error(w, "Origin not allowed", 403)
	// 	return
	// }
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	c := &WSClient{send: make(chan []byte, 256), ws: ws}
	server.Registry.register <- c
	go c.writePump()
	c.readPump(s.Registry)
}