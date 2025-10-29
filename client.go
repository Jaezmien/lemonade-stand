package main

import (
	"git.jaezmien.com/Jaezmien/lemonade-stand/bytebuffer"
	"github.com/gorilla/websocket"
)

type Client struct {
	server *Server
	con   *websocket.Conn
	appid int32

	Send chan []byte
}

func (c *Client) Close() {
	c.con.Close()
}

func (c *Client) Read() {
	for {
		t, message, err := c.con.ReadMessage()
		if err != nil {
			return
		}
		if t != websocket.BinaryMessage {
			return
		}
		if len(message) == 0 {
			return
		}

		data, err := bytebuffer.BytesToBuffer(message)
		if err != nil {
			c.server.Stand.logger.Warn("invalid client message")
			return
		}

		c.server.ReadMessage(data, c.appid)
	}
}
func (c *Client) Write() {
	for message := range c.Send {
		w, err := c.con.NextWriter(websocket.BinaryMessage)	
		if err != nil {
			return
		}
		w.Write(message)
		if err := w.Close(); err != nil {
			return
		}
	}
}
