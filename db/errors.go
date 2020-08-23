package db

import "errors"

// InvalidType is an error returned if a parameter type is wrong (e.g: string instead of int)
var InvalidType = errors.New("invalid parameter type")

// UnknownRoom is an error returned if the user requested for a room that doesn't exists
var UnknownRoom = errors.New("unknown room")

// UnknownRoom is an error returned if the user requested for an user that doesn't exists
var UnknownUser = errors.New("unknown user")
