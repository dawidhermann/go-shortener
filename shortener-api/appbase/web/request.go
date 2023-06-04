package web

import (
	"encoding/json"
	"net/http"
)

// Decode converts JSON data sent in request's body to struct
func Decode(req *http.Request, val any) error {
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}
	return nil
}