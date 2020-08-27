package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	HTTP_URL = "http://localhost:1052/v1/http"
	WS_URL   = "ws://localhost:1051/v1/ws"

	TOKEN = "98"
)

type H map[string]interface{}

func MakeRequest(method, url string, data map[string]interface{}, expectedCode int) (resp *http.Response) {
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
	if resp.StatusCode != expectedCode {
		panic(fmt.Sprintf("wrong status code (expected %d, found %d)", expectedCode, resp.StatusCode))
	}

	return
}

func ParseBody(r *http.Response, payload interface{}) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		panic(err)
	}
}
