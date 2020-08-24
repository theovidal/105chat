package controllers

import "errors"

var (
	// InvalidType is an error returned if a parameter type is wrong (e.g: string instead of int)
	InvalidType = errors.New("invalid parameter type")
	// UnknownRoom is an error returned if the user requested for a room that doesn't exists
	UnknownRoom = errors.New("unknown room")
	// UnknownUser is an error returned if the user requested for an user that doesn't exists
	UnknownUser = errors.New("unknown user")
	// UnknownMessage is an error returned if the user requested for a room message that doesn't exists
	UnknownMessage = errors.New("unknown message")
)
