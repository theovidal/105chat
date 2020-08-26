package utils

import "time"

// H is a shortcut to easily create a JSON map
type H map[string]interface{}

// Error is a shortcut to return errors over the HTTP API
type Error struct {
	// Key of the error, dot-separated strings
	// The last element oftens contains an identifier, e.g the ID of the room that is unknown
	Key string `json:"key"`
	// Message associated with the error
	Message string `json:"message"`
}

// Now returns the timestamp for now (logical.)
func Now() int64 {
	return time.Now().Unix()
}

// Contains checks if an identifiers list has one
func Contains(slice []uint, text uint) bool {
	for _, item := range slice {
		if item == text {
			return true
		}
	}

	return false
}
