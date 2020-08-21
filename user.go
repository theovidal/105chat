package main

import (
	"errors"
	"net/http"
	"strconv"
)

type User struct {
	ID    int    `json:"-"`
	Email string `json:"-"`
	Name  string `json:"name"`
}

var accounts = []User{
	User{
		ID:    98,
		Email: "exybore@becauseofprog.fr",
		Name:  "Théo",
	},
	User{
		ID:    90,
		Email: "didier@becauseofprog.fr",
		Name:  "Didier",
	},
}

func FindUserFromRequest(r *http.Request) (user User, err error) {
	header := r.Header.Get("Authentication")
	token, _ := strconv.Atoi(header)

	user, err = FindUser(token)
	return
}

func FindUser(id int) (User, error) {
	for _, account := range accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return User{}, errors.New("Impossible de trouver le compte")
}
