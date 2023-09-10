package redirectionctrl

import (
	"net/http"

	"github.com/dawidhermann/shortener-redirect/internal/db"
)

type RedirectionController struct {
	store *db.KVStore
}

func New(store *db.KVStore) RedirectionController {
	return RedirectionController{store: store}
}

func (controller RedirectionController) Redirect(w http.ResponseWriter, r *http.Request) {
	targetUrl, err := controller.store.GetUrl(r.URL.Path[1:])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Location", targetUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}
