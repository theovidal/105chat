package tests

import (
	"fmt"
	"net/http"
	"testing"
)

var roomID float64

func TestCreateRoom(t *testing.T) {
	data := H{
		"name":        "  Général  ",
		"avatar_url":  "https://exybore.becauseofprog.fr/img/avatar.png",
		"description": "Pour parler de tout",
		"color":       7764,
	}

	resp := MakeRequest("POST", "/rooms", data, http.StatusCreated)

	var room map[string]interface{}
	ParseBody(resp, &room)

	if room["name"].(string) != "Général" {
		t.Fatal("room doesn't have the expected name (trim isn't applied)")
	}

	roomID = room["id"].(float64)
}

func TestGetRoom(t *testing.T) {
	resp := MakeRequest("GET", fmt.Sprintf("/rooms/%19.f", roomID), nil, http.StatusOK)

	var room map[string]interface{}
	ParseBody(resp, &room)

	if room["name"].(string) != "Général" {
		t.Fatal("room doesn't have the expected name!")
	}
}

func TestDeleteRoom(t *testing.T) {
	resp := MakeRequest("DELETE", fmt.Sprintf("/rooms/%19.f", roomID), nil, http.StatusAccepted)

	var room map[string]interface{}
	ParseBody(resp, &room)

	if room["id"].(float64) != roomID {
		t.Fatal("returned ID doesn't correspond to the room")
	}
}
