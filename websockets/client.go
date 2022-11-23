package websockets

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	id       string
	socket   *websocket.Conn
	outbound chan []byte
}

func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case message, ok := <-c.outbound:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				c.hub.unregister <- c
				return
			}
			if err := c.socket.WriteMessage(websocket.TextMessage, message); err != nil {
				c.socket.Close()
				c.hub.unregister <- c
			}
		}
	}
}
