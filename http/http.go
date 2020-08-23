package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseBody(r *http.Request, payload interface{}) (err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(body, &payload)
}

func Response(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
