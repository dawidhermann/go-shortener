package controllers

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"log"
	"net/http"
	"strconv"
)

type UrlController struct {
	db *sql.DB
}

func NewUrlController(db *sql.DB) *UrlController {
	return &UrlController{db}
}

func (controller UrlController) createShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	fmt.Println(claims)
	userId := 0
	userIdVal, ok := claims["userId"]
	if ok {
		userIdStrVal := fmt.Sprintf("%v", userIdVal)
		id, err := strconv.Atoi(userIdStrVal)
		if err != nil {
			log.Println(err.Error())
			status := http.StatusBadRequest
			errorResult := operationErrorResult{
				Error: operationErrorDetails{Code: status,
					Message: ErrIncorrectUserId.Error()},
			}
			JSONResponse(w, errorResult, status)
			return
		}
		userId = id
	}
	var urlCreateModel urls.UrlCreateModel
	err := json.NewDecoder(r.Body).Decode(&urlCreateModel)
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
	urlId, err := controller.service.CreateUrl(urlCreateModel, userId)
	if err != nil {
		log.Println(err.Error())
		status := http.StatusInternalServerError
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: err.Error()},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	urlLocation := fmt.Sprintf("/urls/%d", urlId)
	okStatus := http.StatusCreated
	w.Header().Set("Location", urlLocation)
	w.Header().Set("Content-Location", urlLocation)
	w.WriteHeader(okStatus)
	okResp := operationOkResult{Result: operationOkDetails{
		Status: okStatus,
		Data:   nil,
	}}
	JSONResponse(w, okResp, okStatus)
}

func (controller UrlController) getUrlHandler(w http.ResponseWriter, r *http.Request) {
	urlId, err := strconv.Atoi(chi.URLParam(r, "urlId"))
	if err != nil {
		log.Println(err.Error())
		errorStatus := http.StatusBadRequest
		errorResponse := operationErrorResult{
			Error: operationErrorDetails{
				Code:    errorStatus,
				Message: "Incorrect url id",
			},
		}
		JSONResponse(w, errorResponse, errorStatus)
		return
	}
	urlData, err := controller.service.GetUrl(urlId)
	if err != nil {
		log.Println(err.Error())
		errorStatus := http.StatusBadRequest
		errorResponse := operationErrorResult{
			Error: operationErrorDetails{
				Code:    errorStatus,
				Message: err.Error(),
			},
		}
		JSONResponse(w, errorResponse, errorStatus)
		return
	}
	url := urls.NewUrlViewModel(urlData)
	okStatus := http.StatusOK
	createdResponse := operationOkResult{Result: operationOkDetails{Status: okStatus, Data: url}}
	JSONResponse(w, createdResponse, okStatus)
}

func (controller UrlController) deleteShortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	urlId, err := strconv.Atoi(chi.URLParam(r, "urlId"))
	if err != nil {
		log.Println(err.Error())
		errorStatus := http.StatusBadRequest
		errorResponse := operationErrorResult{
			Error: operationErrorDetails{
				Code:    errorStatus,
				Message: "Incorrect url id",
			},
		}
		JSONResponse(w, errorResponse, errorStatus)
		return
	}
	err = controller.service.DeleteUrl(urlId)
	if err != nil {
		log.Println(err.Error())
		errorStatus := http.StatusInternalServerError
		errorResponse := operationErrorResult{
			Error: operationErrorDetails{
				Code:    errorStatus,
				Message: err.Error(),
			},
		}
		JSONResponse(w, errorResponse, errorStatus)
		return
	}
	okStatus := http.StatusNoContent
	w.WriteHeader(okStatus)
}
