package api

import (
	"encoding/json"
	"github.com/dawidhermann/go-shortener/internal/db"
	"github.com/dawidhermann/go-shortener/internal/shortener"
	"log"
	"net/http"
	"strings"
)

type shortenUrlSave struct {
	Url string `json:"url"`
}

func createShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var saveUrlModel shortenUrlSave
		err := json.NewDecoder(r.Body).Decode(&saveUrlModel)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortenedUrl := shortener.ShortenUrl()
		err = db.SaveUrl(shortenedUrl, saveUrlModel.Url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Method not implemented", http.StatusNotImplemented)
	}
}

func urlRedirectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		pathData := r.URL.Path
		if pathData[0] == '/' {
			pathData = pathData[1:]
		}
		pathParams := strings.Split(pathData, "/")
		if len(pathParams) > 1 {
			log.Print("Incorrect url")
			http.Error(w, "Incorrect URL", http.StatusBadRequest)
			return
		}
		originalUrl, err := db.GetUrl(pathParams[0])
		if err != nil {
			http.Error(w, "URL not found", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, originalUrl, http.StatusTemporaryRedirect)
	default:
		http.Error(w, "Method not implemented", http.StatusNotImplemented)
	}
}
