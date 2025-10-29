package main

import (
	"git.jaezmien.com/Jaezmien/lemonade-stand/buffer"
	"github.com/gorilla/websocket"
)

type Server struct {
	Stand *LemonadeStand

	Clients map[*Client]bool
	Join chan *Client
	Leave chan *Client

	Quit chan struct{}
}

func NewServer(l *LemonadeStand) *Server {
	return &Server{
		Stand: l,
		Clients: make(map[*Client]bool),
		Join: make(chan *Client),
		Leave: make(chan *Client),
		Quit: make(chan struct{}),
	}
}

func (s *Server) Close() {
	if s.Quit == nil {
		return
	}

	for c := range s.Clients {
		c.Close()
	}

	close(s.Quit)

}

func (s *Server) Run() {
	for {
		select {
		case <-s.Quit:
			return
		case c := <-s.Join:
			s.JoinClient(c)
		case c := <-s.Leave:
			s.CloseClient(c)
		}
	}
}

func (s *Server) JoinClient(c *Client) {
	s.Clients[c] = true

	if s.Stand.HasNotITG() {
		c.Send <- []byte{0x01}
	}
}
func (s *Server) CloseClient(c *Client) {
	if _, ok := s.Clients[c]; ok {
		c.Close()
		delete(s.Clients, c)
	}
}

func (s *Server) Broadcast(data []byte) {
	for c := range s.Clients {
		select {
		case c.Send <- data:
		default:
			s.CloseClient(c)
		}
	}
}
func (s *Server) BroadcastToID(data []byte, appid int32) {
	for c := range s.Clients {
		if c.appid != appid { continue }

		select {
		case c.Send <- data:
		default:
			s.CloseClient(c)
		}
	}
}

func (s *Server) NewClient(con *websocket.Conn, appid int32) *Client {
	c := &Client{
		con: con,
		appid: appid,
		server: s,
	}	
	go c.Read()
	go c.Write()

	s.Join <- c
	return c
}

func (s *Server) GetClientsByID(appid int32) []*Client {
	c := make([]*Client, 0)
	for cl := range s.Clients {
		if cl.appid != appid { continue }
		c = append(c, cl)
	}
	return c
}

func (s *Server) ReadMessage(m []int32, appid int32) {
	split := buffer.SplitBuffer(m)
	
	for i, b := range split {
		buf := s.Stand.writeManager.NewBuffer(appid)
		buf.AppendBuffer(b)
		if i+1 == len(split) {
			buf.Set = buffer.BUFFER_END
		} else {
			buf.Set = buffer.BUFFER_PARTIAL
		}
	}
}
