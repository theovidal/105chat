package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

var roomID float64

func TestCreateRoom(t *testing.T) {
	data := H{
		"name":        "Général",
		"avatar_url":  "https://exybore.becauseofprog.fr/img/avatar.png",
		"description": "Pour parler de tout",
		"color":       7764,
	}

	resp, err := MakeRequest("POST", "/rooms", data)
	if err != nil {
		t.Error(err)
	}
	println(resp.StatusCode)

	var room map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &room)
	if err != nil {
		t.Error(err)
	}

	name := room["name"].(string)
	if name != "Général" {
		t.Error("not good name")
	}

	roomID = room["id"].(float64)
}

func TestGetRoom(t *testing.T) {
	resp, err := MakeRequest("GET", fmt.Sprintf("/rooms/%19.f", roomID), nil)
	if err != nil {
		t.Error(err)
	}
	println(resp.StatusCode)

	var room map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &room)
	if err != nil {
		t.Error(err)
	}

	name := room["name"].(string)
	if name != "Général" {
		t.Error("room doesn't have the expected name")
	}
}

func TestDeleteRoom(t *testing.T) {
	resp, err := MakeRequest("DELETE", fmt.Sprintf("/rooms/%19.f", roomID), nil)
	if err != nil {
		t.Error(err)
	}
	println(resp.StatusCode)

	var room map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &room)
	if err != nil {
		t.Error(err)
	}

	if room["id"].(float64) != roomID {
		t.Error("returned ID doesn't correspond to the room")
	}
}
