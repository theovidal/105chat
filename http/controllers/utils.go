package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// FindIDFromURL extracts a snowflake identifier from the URL, using mux variables
func FindIDFromURL(r *http.Request, value string) (int, error) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars[value])
	if err != nil {
		return 0, InvalidType
	}
	return id, nil
}
