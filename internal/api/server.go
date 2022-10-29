package api

import "net/http"

func StartServer() {
	http.HandleFunc("/url", createShortenUrlHandler)
	http.HandleFunc("/", urlRedirectionHandler)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		return
	}
}
