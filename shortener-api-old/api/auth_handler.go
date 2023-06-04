package api

import (
	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/dawidhermann/shortener-api/internal/db"
	"github.com/dawidhermann/shortener-api/internal/users"
	"net/http"
)

type AuthTokenOkResponse struct {
	Token string `json:"token"`
}

type AuthController struct {
	service     users.ServiceUsers
	authManager auth.AuthenticationManager
}

func NewAuthController(connDb db.SqlConnection, authManager auth.AuthenticationManager) AuthController {
	return AuthController{
		service:     users.NewServiceUsers(connDb),
		authManager: authManager,
	}
}

func (controller AuthController) authHandler(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()
	user, err := controller.service.GetUserByUsername(username)
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
