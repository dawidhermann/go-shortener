package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dawidhermann/go-shortener/internal/db"
	"github.com/dawidhermann/go-shortener/internal/shortener"
	"github.com/go-chi/chi/v5"
)

type shortenUrlSave struct {
	Url string `json:"url"`
}

type operationOkResult struct {
	Result operationOkDetails `json:"result"`
}
type operationOkDetails struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

type operationErrorResult struct {
	Error operationErrorDetails `json:"error"`
}
type operationErrorDetails struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}

func createShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	var saveUrlModel shortenUrlSave
	err := json.NewDecoder(r.Body).Decode(&saveUrlModel)
	if err != nil {
		log.Println(err.Error())
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "Failed to decode request body"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	shortenedUrl := shortener.ShortenUrl()
	err = db.SaveUrl(shortenedUrl, saveUrlModel.Url)
	if err != nil {
		log.Println(err.Error())
		status := http.StatusInternalServerError
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "Failed to save url to database"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	okStatus := http.StatusCreated
	createdResponse := operationOkResult{Result: operationOkDetails{Code: okStatus, Data: shortenedUrl}}
	JSONResponse(w, createdResponse, okStatus)
}

func deleteShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "urlId")
	err := db.DeleteUrl(urlId)
	if err != nil {
		log.Println(err.Error())
		errorStatus := http.StatusInternalServerError
		errorResponse := operationErrorResult{
			Error: operationErrorDetails{
				Code:    errorStatus,
				Message: "Failed to delete URL",
			},
		}
		JSONResponse(w, errorResponse, errorStatus)
		return
	}
	okStatus := http.StatusNoContent
	w.WriteHeader(okStatus)
}

func urlRedirectionHandler(w http.ResponseWriter, r *http.Request) {
	shortenedUrl := chi.URLParam(r, "urlId")
	originalUrl, err := db.GetUrl(shortenedUrl)
	if err != nil {
		log.Println(err.Error())
		status := http.StatusBadRequest
		errorResponse := operationErrorResult{Error: operationErrorDetails{Code: status, Message: "URL not found"}}
		JSONResponse(w, errorResponse, status)
		return
	}
	http.Redirect(w, r, originalUrl, http.StatusTemporaryRedirect)
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
