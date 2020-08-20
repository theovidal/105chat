package main

import "fmt"

type Message struct {
	Author  *User
	Content string
}

func (m *Message) String() string {
	return fmt.Sprint(m.Author.Name, " : ", m.Content)
}
