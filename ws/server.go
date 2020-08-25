package ws

import (
	"time"

	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
)

// Server starts the WebSocket server
func Server(ws *websocket.Conn) {
	websocket.JSON.Send(ws, Event{
		Event: AUTHENTICATION_NEEDED,
	})

	var user db.User
	var trials int

	for {
		if trials >= 3 {
			websocket.JSON.Send(ws, Event{
				Event: CLOSE,
				Data:  "Too many failed authentications",
			})
			ws.Close()
			break
		}

		ws.SetReadDeadline(time.Now().Add(time.Minute))
		var event Event
		if err := websocket.JSON.Receive(ws, &event); err != nil {
			Pipeline <- Event{
				Event: USER_DISCONNECT,
				Data:  user.ID,
			}
			station.Lock()
			delete(station.clients, user.ID)
			station.Unlock()
			break
		}

		if user.ID == 0 {
			if event.Event == CONNECT {
				token, valid := event.Data.(string)
				if !valid {
					websocket.JSON.Send(ws, Event{
						Event: ERROR,
						Data:  ERROR402,
					})
					continue
				}

				var err error
				user, err = db.FindUserByToken(token)
				if err != nil {
					websocket.JSON.Send(ws, Event{
						Event: AUTHENTICATION_FAIL,
					})
					trials += 1
				} else {
					websocket.JSON.Send(ws, Event{
						Event: AUTHENTICATION_SUCCESS,
						Data:  &user,
					})
					station.Lock()
					station.clients[user.ID] = ws
					station.Unlock()
					Pipeline <- Event{
						Event: USER_CONNECT,
						Data:  user.ID,
					}
				}
			} else {
				websocket.JSON.Send(ws, Event{
					Event: ERROR,
					Data:  ERROR401,
				})
			}
			continue
		}

		handler, valid := clientEvents[event.Event]
		if !valid {
			websocket.JSON.Send(ws, Event{
				Event: ERROR,
				Data:  ERROR404,
			})
		} else {
			handler(ws, &user, &event)
		}
	}
}
