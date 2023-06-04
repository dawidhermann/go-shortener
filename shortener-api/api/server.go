package api

import (
	"net/http"

	"github.com/dawidhermann/shortener-api/api/controllers/v1/userctrl"
	v1 "github.com/dawidhermann/shortener-api/business/web/v1"
	"github.com/dawidhermann/shortener-api/internal/core/user"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

// type Handler func(c echo.Context) error

type App struct {
	*echo.Echo
}

func NewApp() *App {
	return &App{echo.New()}
}

// func (app App) Handle(method string, path string, handler Handler) {
// 	h := func(c echo.Context) {

// 	}
// 	app.Router().Add(method, path, h)
// }

func APIMux(db *sqlx.DB) *App {
	app := NewApp()
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
		Core: user.NewUserCore(db),
	}
	group := app.Group("/user")
	group.POST("/", usrctrl.CreateUser)
	group.GET("/:id", usrctrl.GetUserById)
	group.PATCH("/:id", usrctrl.UpdateUser)
	group.DELETE("/:id", usrctrl.DeleteUser)
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
	switch {
	case v1.IsRequestError(err):
		reqErr := v1.GetRequestError(err)
		c.JSON(reqErr.Status, v1.ErrorResponse{Error: reqErr.Error()})
	default:
		status := http.StatusInternalServerError
		c.JSON(status, v1.ErrorResponse{Error: http.StatusText(status)})
	}
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
