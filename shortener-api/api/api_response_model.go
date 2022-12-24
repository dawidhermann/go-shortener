package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type operationOkResult struct {
	Result operationOkDetails `json:"result"`
}
type operationOkDetails struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type operationErrorResult struct {
	Error operationErrorDetails `json:"error"`
}
type operationErrorDetails struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}

func JSONResponse(w http.ResponseWriter, data interface{}, code int) {
	encData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error occurred while encoding error response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write(encData)
}
