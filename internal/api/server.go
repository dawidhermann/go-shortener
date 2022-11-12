package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/url", func(r chi.Router) {
		r.Post("/", createShortenUrlHandler)
		r.Route("/{urlId}", func(r chi.Router) {
			r.Delete("/", deleteShortenUrlHandler)
		})
	})
	r.Get("/{urlId}", urlRedirectionHandler)

	err := http.ListenAndServe(":8090", r)
	if err != nil {
		return
	}
}
