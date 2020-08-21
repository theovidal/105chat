package main

import (
	"fmt"
	"strconv"

	"golang.org/x/net/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Event)

func WebSocketServer(ws *websocket.Conn) {
	websocket.Message.Send(ws, "Merci d'inscrire votre identifiant pour vous connecter.")

	var user User
	var err error
	var input string
	var trials int
	var closed bool
	for {
		websocket.Message.Receive(ws, &input)
		id, _ := strconv.Atoi(input)
		user, err = FindUser(id)
		if err != nil {
			websocket.Message.Send(ws, "Identifiant invalide!")
			trials += 1
			if trials == 3 {
				websocket.Message.Send(ws, "3 tentatives échouées - fermeture de la session")
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

	websocket.Message.Send(ws, fmt.Sprintf("Bienvenue dans le chat, %s !", user.Name))
	clients[ws] = true

	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			delete(clients, ws)
			break
		}

		websocket.Message.Send(ws, "Aucune commande disponible (pour le moment)")
	}
}

func HandleBroadcasts() {
	for {
		event := <-broadcast
		for client := range clients {
			err := websocket.JSON.Send(client, event)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
