package api

import (
	"net/http"
)

type shortenUrlSave struct {
	Url string `json:"url"`
}

func createShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	//var saveUrlModel shortenUrlSave
	//err := json.NewDecoder(r.Body).Decode(&saveUrlModel)
	//if err != nil {
	//	log.Println(err.Error())
	//	status := http.StatusBadRequest
	//	errorResult := operationErrorResult{
	//		Error: operationErrorDetails{Status: status,
	//			Message: "Failed to decode request body"},
	//	}
	//	JSONResponse(w, errorResult, status)
	//	return
	//}
	//shortenedUrl := shortener.ShortenUrl()
	//err = db.SaveUrl(shortenedUrl, saveUrlModel.Url)
	//if err != nil {
	//	log.Println(err.Error())
	//	status := http.StatusInternalServerError
	//	errorResult := operationErrorResult{
	//		Error: operationErrorDetails{Status: status,
	//			Message: "Failed to save url to database"},
	//	}
	//	JSONResponse(w, errorResult, status)
	//	return
	//}
	//okStatus := http.StatusCreated
	//createdResponse := operationOkResult{Result: operationOkDetails{Status: okStatus, Data: shortenedUrl}}
	//JSONResponse(w, createdResponse, okStatus)
}

func deleteShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	//urlId := chi.URLParam(r, "urlId")
	//err := db.DeleteUrl(urlId)
	//if err != nil {
	//	log.Println(err.Error())
	//	errorStatus := http.StatusInternalServerError
	//	errorResponse := operationErrorResult{
	//		Error: operationErrorDetails{
	//			Status:    errorStatus,
	//			Message: "Failed to delete URL",
	//		},
	//	}
	//	JSONResponse(w, errorResponse, errorStatus)
	//	return
	//}
	//okStatus := http.StatusNoContent
	//w.WriteHeader(okStatus)
}

func urlRedirectionHandler(w http.ResponseWriter, r *http.Request) {
	//shortenedUrl := chi.URLParam(r, "urlId")
	//originalUrl, err := db.GetUrl(shortenedUrl)
	//if err != nil {
	//	log.Println(err.Error())
	//	status := http.StatusBadRequest
	//	errorResponse := operationErrorResult{Error: operationErrorDetails{Status: status, Message: "URL not found"}}
	//	JSONResponse(w, errorResponse, status)
	//	return
	//}
	//http.Redirect(w, r, originalUrl, http.StatusTemporaryRedirect)
}
