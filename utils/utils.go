package utils

// H is a shortcut to easily create a JSON map
type H map[string]interface{}

type Error struct {
	Key     string `json:"key"`
	Message string `json:"message"`
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
