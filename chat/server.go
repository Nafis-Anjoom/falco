package chat

import (
	"log"
)

type Server struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.Register:
			log.Println("client connected:", client.conn.RemoteAddr().String())
			s.Clients[client] = true
		case client := <-s.unregister:
			log.Println("client disconnected:", client.conn.RemoteAddr().String())
			if _, ok := s.Clients[client]; ok {
				delete(s.Clients, client)
			}
		// case message := <-s.Broadcast:
		//           log.Println("broadcasting message:", string(message))
		// 	for client := range s.Clients {
		// 		client.buff <- message
		// 	}
		}
	}
}
