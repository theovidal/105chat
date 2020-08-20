package main

import "errors"

type User struct {
	ID    int
	Email string
	Name  string
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

func FindUser(id int) (User, error) {
	for _, account := range accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return User{}, errors.New("Impossible de trouver le compte")
}
