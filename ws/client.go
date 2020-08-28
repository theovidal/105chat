package ws

import (
	"time"

	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
)

type Client struct {
	User          *db.User
	Connection    *websocket.Conn
	Server        *Server
	Pipeline      chan *Event
	ClosePipeline chan bool
}

func NewClient(connection *websocket.Conn, server *Server) *Client {
	return &Client{
		User:          &db.User{},
		Connection:    connection,
		Server:        server,
		Pipeline:      make(chan *Event),
		ClosePipeline: make(chan bool),
	}
}

func (c *Client) Listen() {
	go c.ListenClientEvents()
	c.ListenServerEvents()
}

func (c *Client) ListenClientEvents() {
	for {
		select {
		case <-c.ClosePipeline:
			c.Server.DisconnectPipeline <- c
			c.ClosePipeline <- true
			return

		default:
			_ = websocket.JSON.Send(c.Connection, Event{
				Event: AUTHENTICATION_NEEDED,
			})

			var user db.User
			var trials int

			for {
				if trials >= 3 {
					_ = websocket.JSON.Send(c.Connection, Event{
						Event: CLOSE,
						Data:  "Too many failed authentications",
					})
					c.ClosePipeline <- true
					break
				}

				_ = c.Connection.SetReadDeadline(time.Now().Add(time.Minute))
				var event Event
				if err := websocket.JSON.Receive(c.Connection, &event); err != nil {
					c.ClosePipeline <- true
					break
				}

				if user.ID == 0 {
					if event.Event == CONNECT {
						token, valid := event.Data.(string)
						if !valid {
							_ = websocket.JSON.Send(c.Connection, Event{
								Event: ERROR,
								Data:  ERROR402,
							})
							continue
						}

						var err error
						user, err = db.FindUserByToken(token)
						if err != nil {
							_ = websocket.JSON.Send(c.Connection, Event{
								Event: AUTHENTICATION_FAIL,
							})
							trials += 1
						} else {
							if user.Disabled {
								_ = websocket.JSON.Send(c.Connection, Event{
									Event: ERROR,
									Data:  ERROR405,
								})
								c.ClosePipeline <- true
							}
							_ = websocket.JSON.Send(c.Connection, Event{
								Event: AUTHENTICATION_SUCCESS,
								Data:  &user,
							})
							c.User = &user
							c.Server.ConnectPipeline <- c
						}
					} else {
						_ = websocket.JSON.Send(c.Connection, Event{
							Event: ERROR,
							Data:  ERROR401,
						})
					}
					continue
				}

				handler, valid := clientEvents[event.Event]
				if !valid {
					_ = websocket.JSON.Send(c.Connection, Event{
						Event: ERROR,
						Data:  ERROR404,
					})
				} else {
					handler(c.Connection, &user, &event)
				}
			}
		}
	}
}

func (c *Client) ListenServerEvents() {
	for {
		select {
		case <-c.ClosePipeline:
			c.Server.DisconnectPipeline <- c
			c.ClosePipeline <- true
			return

		case event := <-c.Pipeline:
			err := websocket.JSON.Send(c.Connection, event)
			if err != nil {
				c.ClosePipeline <- true
			}
		}
	}
}
