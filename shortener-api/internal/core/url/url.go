package url

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/dawidhermann/shortener-api/internal/auth/"
	"github.com/dawidhermann/shortener-api/internal/core/url/store"
	"github.com/dawidhermann/shortener-api/internal/rpc"
	"github.com/dawidhermann/shortener-api/internal/sys/validate"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotValid = errors.New("user is not valid")
)

type Core struct {
	store   *store.Store
	rpcConn rpc.ConnRpc
}

func NewUrlCore(db *sqlx.DB) *Core {
	return &Core{
		store: store.NewUrlStore(db),
	}
}

func (core *Core) Create(ctx context.Context, urlCreateModel UrlCreateViewModel, userClaims auth.UserClaims) (Url, error) {
	err := validate.ValidateStruct(urlCreateModel)
	if err != nil {
		return Url{}, fmt.Errorf("url data validation error: %w", err)
	}
	urlAddress, err := url.Parse(urlCreateModel.TargetUrl)
	if err != nil {
		return Url{}, fmt.Errorf("failed to parse url: %w", err)
	}
	urlKey, err := core.rpcConn.CreateShortenUrl(urlAddress.String())
	if err != nil {
		return Url{}, fmt.Errorf("failed to shorten url: %w", err)
	}
	url := Url{
		Key:    urlKey,
		UserId: userClaims.UserId,
	}
	createUrlRes, err := core.store.Create(ctx, toDbUrl(url))
	if err != nil {
		return Url{}, fmt.Errorf("failed to create url: %w", err)
	}
	url.UrlId = createUrlRes.UrlId
	url.DateCreated = createUrlRes.CreatedAt
	return url, nil
}

func (core *Core) GetById(ctx context.Context, id uuid.UUID) (Url, error) {
	dbUrl, err := core.store.GetById(ctx, id)
	if err != nil {
		return Url{}, fmt.Errorf("failed to fetch url: %w", err)
	}
	return toUrl(dbUrl), nil
}

func (core *Core) Delete(ctx context.Context, id uuid.UUID) error {
	err := core.store.Delete(ctx, id)
	if err != nil {
		fmt.Errorf("failed to delete url: %w", err)
	}
}