package api

// func NewAPIMux() *App {
// 	app := NewApp()

// 	authController := controllers.NewAuthController(nil)
// 	urlController := controllers.NewUrlController(nil)
// 	app.Post("/auth", authController.HandleAuthentication)
// 	app.Route("/url", func(r chi.Router) {
// 		//r.Use(jwtauth.Verifier(authManager.TokenAuth))
// 		app.Group(func(r chi.Router) {
// 			app.Post("/", urlController.createShortenUrlHandler)
// 		})
// 		app.Group(func(r chi.Router) {
// 			//r.Use(jwtauth.Authenticator)
// 			app.Route("/{urlId}", func(r chi.Router) {
// 				app.Delete("/", urlController.deleteShortenUrlHandler)
// 				app.Get("/", urlController.getUrlHandler)
// 			})
// 		})
// 	})
// 	return app
// }
