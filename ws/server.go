package ws

import (
	"sync"

	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
)

type StationEntity struct {
	sync.RWMutex
	clients map[string]*Client
}

func (s *StationEntity) Initialize() {
	s.clients = make(map[string]*Client)
}

func (s *StationEntity) GetUser(token string) (user *db.User, found bool) {
	client, found := s.clients[token]
	if !found {
		return &db.User{}, false
	} else {
		return client.User, true
	}
}

func (s *StationEntity) EditUser(user *db.User) {
	s.clients[user.Token].User = user
}

// Station stores connections to the WebSocket server
// In the future, would be replaced by a concurrent-specific map : https://github.com/streamrail/concurrent-map
var Station StationEntity

type Server struct {
	ConnectPipeline    chan *Client
	DisconnectPipeline chan *Client
}

func NewServer() *Server {
	return &Server{
		ConnectPipeline:    make(chan *Client),
		DisconnectPipeline: make(chan *Client),
	}
}

func (s *Server) AddClient(client *Client) {
	Station.Lock()
	Station.clients[client.User.Token] = client
	Station.Unlock()
}

func (s *Server) DeleteClient(client *Client) {
	Station.Lock()
	delete(Station.clients, client.User.Token)
	Station.Unlock()
}

func (s *Server) Handle(connection *websocket.Conn) {
	client := NewClient(connection, s)
	client.Listen()
	connection.Close()
}

func (s *Server) Listen() {
	Station.Initialize()
	go HandlePipeline()

	for {
		select {
		case client := <-s.ConnectPipeline:
			s.AddClient(client)
			Pipeline <- Event{
				Event: USER_CONNECT,
				Data:  client.User.ID,
			}

		case client := <-s.DisconnectPipeline:
			s.DeleteClient(client)
			Pipeline <- Event{
				Event: USER_DISCONNECT,
				Data:  client.User.ID,
			}
		}
	}
}
