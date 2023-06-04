package api

import (
	"fmt"
	"github.com/dawidhermann/shortener-api/config"
	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/dawidhermann/shortener-api/internal/db"
	"github.com/dawidhermann/shortener-api/internal/rpc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"log"
	"net/http"
)

func StartServer(serverConfig config.ServerConfig, connRpc rpc.ConnRpc, connDb db.SqlConnection, authManager auth.AuthenticationManager) {
	urlController := NewUrlController(connRpc, connDb)
	usersController := newUsersController(connDb)
	authController := NewAuthController(connDb, authManager)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/auth", authController.authHandler)
	r.Route("/url", func(r chi.Router) {
		r.Use(jwtauth.Verifier(authManager.TokenAuth))
		r.Group(func(r chi.Router) {
			r.Post("/", urlController.createShortenUrlHandler)
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Authenticator)
			r.Route("/{urlId}", func(r chi.Router) {
				r.Delete("/", urlController.deleteShortenUrlHandler)
				r.Get("/", urlController.getUrlHandler)
			})
		})
	})
	r.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/", usersController.createUser)
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(authManager.TokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Route("/{userId}", func(r chi.Router) {
				r.Delete("/", usersController.deleteUser)
				r.Get("/", usersController.getUser)
				r.Patch("/", usersController.updateUser)
			})
		})
	})
	//r.Get("/{urlId}", urlRedirectionHandler)
	httpPort := fmt.Sprintf(":%s", serverConfig.ServerPort)
	err := http.ListenAndServe(httpPort, r)
	if err != nil {
		log.Fatalf("Cannot start app on port %s because of %s", httpPort, err)
		return
	}
	log.Printf("Listening on port: %s", httpPort)
}
