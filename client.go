package main

import (
	"log/slog"
	"time"

	"git.jaezmien.com/Jaezmien/lemonade-stand/bytebuffer"
	"git.jaezmien.com/Jaezmien/lemonade-stand/encoder"
	"github.com/gorilla/websocket"
)

var clientWriteWait = time.Second * 10
var clientPongWait = time.Second * 60
var clientPingPeriod = (clientPongWait * 9) / 10

type Client struct {
	server *Server
	con    *websocket.Conn
	appid  int32

	Send   chan []byte
	exited bool
}

func (c *Client) Close() {
	if c.exited {
		return
	}
	c.exited = true

	c.con.Close()
}

func (c *Client) Read() {
	c.con.SetReadDeadline(time.Now().Add(clientPongWait))
	c.con.SetPongHandler(func(appData string) error {
		c.con.SetReadDeadline(time.Now().Add(clientPongWait))
		return nil
	})

readLoop:
	for {
		slogAppID := slog.Int("appid", int(c.appid))

		t, message, err := c.con.ReadMessage()
		if c.exited {
			return
		}
		if err != nil {
			c.server.Stand.logger.Debug("error on read messsage", slog.Any("error", err), slogAppID)
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
				c.server.Stand.logger.Warn("received invalid client buffer", slogAppID)
				continue
			}
		case websocket.TextMessage:
			// Fine, we'll handle your TextMessage
			data, err = encoder.StringToBuffer(string(message))
			if err != nil {
				c.server.Stand.logger.Warn("received invalid client text", slogAppID)
				continue
			}
		default:
			c.server.Stand.logger.Debug("invalid message type", slogAppID)
			continue readLoop
		}

		c.server.Stand.logger.Debug("received client message", slogAppID)
		c.server.ReadMessage(data, c.appid)
	}

	c.server.CloseClient(c)
}
func (c *Client) Write() {
	ticker := time.NewTicker(clientPingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.con.SetWriteDeadline(time.Now().Add(clientWriteWait))
			if !ok {
				c.con.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

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
		case <-ticker.C:
			c.con.SetWriteDeadline(time.Now().Add(clientWriteWait))
			if err := c.con.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
