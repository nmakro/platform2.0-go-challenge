package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ValidateRequest(r *http.Request, i interface{}) error {
	rawRequestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawRequestBody, i)
	if err != nil {
		return err
	}
	return nil
}

func ResponseWithJSON(status int, i interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(i)
}
