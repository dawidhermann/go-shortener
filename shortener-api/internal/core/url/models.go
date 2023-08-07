package url

import (
	"time"

	"github.com/dawidhermann/shortener-api/internal/core/url/store"
	"github.com/google/uuid"
)

type UrlViewModel struct {
	UrlId  uuid.UUID `json:"id"`
	Key    string    `json:"key"`
	UserId uuid.UUID `json:"userId"`
}

type UrlCreateViewModel struct {
	TargetUrl string `json:"targetUrl" validate:"required,http_url"`
}

type Url struct {
	UrlId       uuid.UUID
	Key         string
	UserId      uuid.UUID
	DateCreated time.Time
	DateUpdated time.Time
}

func NewUrlViewModel(url Url) UrlViewModel {
	return UrlViewModel{
		UrlId:  url.UrlId,
		Key:    url.Key,
		UserId: url.UserId,
	}
}

func toUrl(dbUrl store.DbUrl) Url {
	return Url{
		UrlId:       dbUrl.ID,
		Key:         dbUrl.Key,
		UserId:      dbUrl.UserId,
		DateCreated: dbUrl.DateCreated,
		DateUpdated: dbUrl.DateUpdated,
	}
}

func toDbUrl(url Url) store.DbUrl {
	return store.DbUrl{
		ID:          url.UrlId,
		Key:         url.Key,
		UserId:      url.UserId,
		DateCreated: url.DateCreated,
		DateUpdated: url.DateUpdated,
	}
}
