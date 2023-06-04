package controllers

import (
	"database/sql"
	"github.com/dawidhermann/shortener-api/internal/auth"
	"net/http"
)

type AuthController struct {
	Auth *auth.AuthManager
}

func NewAuthController(db *sql.DB) AuthController {
	return AuthController{
		Auth: auth.NewAuthManager(db),
	}
}

func (controller AuthController) HandleAuthentication(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	user, err := controller.Auth.GetUserByUsername(username)
	if err != nil {
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "incorrect credentials"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	if user.Password != password {
		status := http.StatusBadRequest
		errorResult := operationErrorResult{
			Error: operationErrorDetails{Code: status,
				Message: "incorrect credentials"},
		}
		JSONResponse(w, errorResult, status)
		return
	}
	_, tokenString, _ := controller.authManager.EncodeJwtToken(map[string]interface{}{
		"userId": user.UserId,
		"email":  user.Email,
	})
	JSONResponse(w, operationOkResult{
		Result: operationOkDetails{
			Status: http.StatusOK,
			Data: AuthTokenOkResponse{
				Token: tokenString,
			},
		},
	}, http.StatusOK)
}
