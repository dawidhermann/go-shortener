// Main app mux
package api

import (
	"net/http"

	"github.com/dawidhermann/shortener-api/api/controllers/v1/authctrl"
	"github.com/dawidhermann/shortener-api/api/controllers/v1/urlctrl"
	"github.com/dawidhermann/shortener-api/api/controllers/v1/userctrl"
	v1 "github.com/dawidhermann/shortener-api/api/v1"
	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/dawidhermann/shortener-api/internal/core/url"
	"github.com/dawidhermann/shortener-api/internal/core/user"
	"github.com/dawidhermann/shortener-api/internal/rpc"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// type Handler func(c echo.Context) error

type App struct {
	RpcConn rpc.ConnRpc
	*echo.Echo
}

func newApp(rpcConn rpc.ConnRpc) *App {
	return &App{rpcConn, echo.New()}
}

type AppConfig struct {
	Auth    auth.Auth
	Db      *sqlx.DB
	RpcConn rpc.ConnRpc
}

// Create new app mux
func APIMux(cfg AppConfig) *App {
	app := newApp(cfg.RpcConn)
	usrctrl := userctrl.UsersController{
		Core: user.NewUserCore(cfg.Db),
	}
	urlctrl := urlctrl.UrlsController{
		Core: url.NewUrlCore(cfg.Db, cfg.RpcConn),
	}
	authctrl := authctrl.AuthController{
		Core: user.NewUserCore(cfg.Db),
		Auth: cfg.Auth,
	}
	config := echojwt.Config{
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.Auth.Secret), nil
		}}
	authMiddleware := echojwt.WithConfig(config)
	appApi := app.Group("/api")
	apiVersionGroup := appApi.Group("/v1")
	userGroup := apiVersionGroup.Group("/user")
	userGroup.POST("/", usrctrl.CreateUser)
	userGroup.GET("/:id", usrctrl.GetUserById, authMiddleware)
	userGroup.PATCH("/:id", usrctrl.UpdateUser, authMiddleware)
	userGroup.DELETE("/:id", usrctrl.DeleteUser, authMiddleware)
	urlGroup := apiVersionGroup.Group("/url", authMiddleware)
	urlGroup.POST("/", urlctrl.CreateUrl)
	urlGroup.DELETE("/:id", urlctrl.DeleteUrl)
	apiVersionGroup.POST("/auth", authctrl.LoginUser)
	app.HTTPErrorHandler = errorHandler
	return app
}

func errorHandler(err error, c echo.Context) {
	echoErr, isEchoErr := err.(*echo.HTTPError)
	switch {
	case v1.IsRequestError(err):
		reqErr := v1.GetRequestError(err)
		c.JSON(reqErr.Status, v1.ErrorResponse{Error: reqErr.Error()})
	case isEchoErr:
		var message string
		if msg, ok := echoErr.Message.(string); ok {
			message = msg
		} else {
			message = http.StatusText(echoErr.Code)
		}
		c.JSON(echoErr.Code, v1.ErrorResponse{Error: message})
	default:
		status := http.StatusInternalServerError
		c.JSON(status, v1.ErrorResponse{Error: http.StatusText(status)})
	}
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
