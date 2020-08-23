package db

import "errors"

var InvalidType = errors.New("invalid parameter type")
var UnknownRoom = errors.New("unknown room")
var UnknownUser = errors.New("unknown user")
