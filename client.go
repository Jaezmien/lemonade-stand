package main

import (
	"log/slog"

	"git.jaezmien.com/Jaezmien/lemonade-stand/bytebuffer"
	"git.jaezmien.com/Jaezmien/lemonade-stand/encoder"
	"github.com/gorilla/websocket"
)

type Client struct {
	server *Server
	con   *websocket.Conn
	appid int32

	Send chan []byte
	exited bool
}

func (c *Client) Close() {
	if c.exited { return }
	c.exited = true

	c.con.Close()
}

func (c *Client) Read() {
	readLoop:
	for {
		t, message, err := c.con.ReadMessage()
		if c.exited {
			return
		}
		if err != nil {
			c.server.Stand.logger.Debug("error on read messsage", slog.Any("error", err))
			break
		}

		if len(message) == 0 {
			continue
		}

		var data []int32
		switch t {
		case websocket.BinaryMessage:
			data, err = bytebuffer.BytesToBuffer(message)
			if err != nil {
				c.server.Stand.logger.Warn("received invalid client buffer", slog.Int("appid", int(c.appid)))
				continue
			}
		case websocket.TextMessage:
			// Fine, we'll handle your TextMessage
			data, err = encoder.StringToBuffer(string(message))
			if err != nil {
				c.server.Stand.logger.Warn("received invalid client text", slog.Int("appid", int(c.appid)))
				continue
			}
		default:
			c.server.Stand.logger.Debug("invalid message type")
			continue readLoop
		}

		c.server.ReadMessage(data, c.appid)
	}

	c.server.CloseClient(c)
}
func (c *Client) Write() {
	for message := range c.Send {
		w, err := c.con.NextWriter(websocket.BinaryMessage)	
		if err != nil {
			c.server.Stand.logger.Debug("error while creating client writer", slog.Any("error", err))
			continue
		}
		w.Write(message)
		if err := w.Close(); err != nil {
			c.server.Stand.logger.Debug("error while closing client writer", slog.Any("error", err))
			continue
		}
	}
}
