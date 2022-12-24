package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dawidhermann/shortener-api/internal/users"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	var userCreateModel users.UserCreateViewModel
	err := json.NewDecoder(r.Body).Decode(&userCreateModel)
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
	userId, err := users.CreateUser(userCreateModel)
	userLocation := fmt.Sprintf("/users/%d", userId)
	w.Header().Set("Location", userLocation)
	w.Header().Set("Content-Location", userLocation)
	w.WriteHeader(http.StatusCreated)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	var userPatchModel users.UserPatchModel
	err := json.NewDecoder(r.Body).Decode(&userPatchModel)
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
	err = users.UpdateUser(userId, userPatchModel)
	if err != nil {
		log.Println(err.Error())
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: err.Error()},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	userLocation := fmt.Sprintf("/users/%s", userId)
	w.Header().Set("Location", userLocation)
	w.Header().Set("Content-Location", userLocation)
	w.WriteHeader(http.StatusOK)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	userData, err := users.GetUser(userId)
	if errors.Is(err, users.ErrUserNotFound) {
		status := http.StatusNotFound
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: err.Error()},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	if err != nil {
		status := http.StatusInternalServerError
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "Failed to fetch user"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	user := users.NewUserViewModel(userData)
	okStatus := http.StatusOK
	createdResponse := operationOkResult{Result: operationOkDetails{Status: okStatus, Data: user}}
	JSONResponse(w, createdResponse, okStatus)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	err := users.DeleteUser(userId)
	log.Println(err)
	if errors.Is(err, users.ErrIncorrectUserId) {
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: err.Error()},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	if err != nil {
		status := http.StatusInternalServerError
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "Failed to delete user"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
