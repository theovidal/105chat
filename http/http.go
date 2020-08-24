package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ParseBody parses request's body (JSON data) into an interface
func ParseBody(r *http.Request, payload interface{}) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, &payload)
}

// Response encodes a JSON response to the user
func Response(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Event", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
