package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dawidhermann/shortener-api/internal/db"
	"github.com/dawidhermann/shortener-api/internal/users"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type UsersController struct {
	service users.ServiceUsers
}

func newUsersController(connDb db.SqlConnection) UsersController {
	return UsersController{
		service: users.NewServiceUsers(connDb),
	}
}

func (controller UsersController) createUser(w http.ResponseWriter, r *http.Request) {
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
	userId, err := controller.service.CreateUser(userCreateModel)
	userLocation := fmt.Sprintf("/users/%d", userId)
	w.Header().Set("Location", userLocation)
	w.Header().Set("Content-Location", userLocation)
	w.WriteHeader(http.StatusCreated)
}

func (controller UsersController) updateUser(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		log.Println(err.Error())
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "Incorrect user id"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	var userPatchModel users.UserPatchModel
	err = json.NewDecoder(r.Body).Decode(&userPatchModel)
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
	err = controller.service.UpdateUser(userId, userPatchModel)
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

func (controller UsersController) getUser(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "Incorrect user id"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	userData, err := controller.service.GetUser(userId)
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

func (controller UsersController) deleteUser(w http.ResponseWriter, r *http.Request) {
	userIdParam := chi.URLParam(r, "userId")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "Incorrect user id"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	err = controller.service.DeleteUser(userId)
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
