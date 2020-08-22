package ws

import (
	"golang.org/x/net/websocket"

	"github.com/theovidal/105chat/db"
)

func Server(ws *websocket.Conn) {
	websocket.JSON.Send(ws, H{
		"code":    2,
		"message": "Please identify using your token",
	})

	var user db.User
	var err error
	var token string
	var trials int
	var closed bool
	for {
		websocket.Message.Receive(ws, &token)
		user, err = db.FindUserByToken(token)
		if err != nil {
			websocket.JSON.Send(ws, H{
				"code":    0,
				"message": "Invalid token",
			})
			trials += 1
			if trials == 3 {
				websocket.JSON.Send(ws, H{
					"code":    0,
					"message": "Too many failed requests, closing connection",
				})
				ws.Close()
				closed = true
			}
		} else {
			break
		}
	}

	if closed {
		return
	}

	websocket.JSON.Send(ws, H{
		"code": 1,
		"data": &user,
	})
	clients[ws] = true

	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			delete(clients, ws)
			break
		}

		websocket.JSON.Send(ws, H{
			"code":    0,
			"message": "Unknown command",
		})
	}
}
