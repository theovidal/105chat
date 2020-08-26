package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	HTTP_URL = "http://localhost:1052/v1/http"
	WS_URL   = "ws://localhost:1051/v1/ws"

	TOKEN = "98"
)

type H map[string]interface{}

func MakeRequest(method, url string, data map[string]interface{}) (resp *http.Response, err error) {
	body, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, HTTP_URL+url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", TOKEN)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	return
}
