package api

import (
	"net/http"

	"github.com/dawidhermann/shortener-api/api/controllers/v1/authctrl"
	"github.com/dawidhermann/shortener-api/api/controllers/v1/userctrl"
	v1 "github.com/dawidhermann/shortener-api/api/v1"
	"github.com/dawidhermann/shortener-api/internal/auth"
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

func NewApp(rpcConn rpc.ConnRpc) *App {
	return &App{rpcConn, echo.New()}
}

type AppConfig struct {
	Auth    auth.Auth
	Db      *sqlx.DB
	RpcConn rpc.ConnRpc
}

// func (app App) Handle(method string, path string, handler Handler) {
// 	h := func(c echo.Context) {

// 	}
// 	app.Router().Add(method, path, h)
// }

func APIMux(cfg AppConfig) *App {
	app := NewApp(cfg.RpcConn)
	//r.Post("/auth", authController.authHandler)
	// group := app.Group("/url")
	// group.POST("/", urlController.createShortenUrlHandler)
	// app.Route("/url", func(r chi.Router) {
	// 	//r.Use(jwtauth.Verifier(authManager.TokenAuth))
	// 	app.Group(func(r chi.Router) {
	// 		app.Post("/", urlController.createShortenUrlHandler)
	// 	})
	// 	app.Group(func(r chi.Router) {
	// 		//r.Use(jwtauth.Authenticator)
	// 		app.Route("/{urlId}", func(r chi.Router) {
	// 			app.Delete("/", urlController.deleteShortenUrlHandler)
	// 			app.Get("/", urlController.getUrlHandler)
	// 		})
	// 	})
	// })
	usrctrl := userctrl.UsersController{
		Core: user.NewUserCore(cfg.Db),
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
	group := app.Group("/user")
	group.POST("/", usrctrl.CreateUser)
	group.GET("/:id", usrctrl.GetUserById, authMiddleware)
	group.PATCH("/:id", usrctrl.UpdateUser, authMiddleware)
	group.DELETE("/:id", usrctrl.DeleteUser, authMiddleware)
	app.POST("/auth", authctrl.LoginUser)
	app.HTTPErrorHandler = errorHandler
	//r.Route("/user", func(r chi.Router) {
	//	r.Group(func(r chi.Router) {
	//		r.Post("/", usersController.createUser)
	//	})
	//	r.Group(func(r chi.Router) {
	//		r.Use(jwtauth.Verifier(authManager.TokenAuth))
	//		r.Use(jwtauth.Authenticator)
	//		r.Route("/{userId}", func(r chi.Router) {
	//			r.Delete("/", usersController.deleteUser)
	//			r.Get("/", usersController.getUser)
	//			r.Patch("/", usersController.updateUser)
	//		})
	//	})
	//})
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
