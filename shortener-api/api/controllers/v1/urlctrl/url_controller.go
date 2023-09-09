// Defines endpoints for url management
package urlctrl

import (
	"errors"
	"fmt"
	"net/http"

	v1 "github.com/dawidhermann/shortener-api/api/v1"
	"github.com/dawidhermann/shortener-api/internal/core/url"
	"github.com/dawidhermann/shortener-api/internal/web"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UrlsController struct {
	Core *url.Core
}

var (
	ErrBindingRequestData = errors.New("failed to bind request data")
	ErrInvalidId          = errors.New("invalid url id")
)

// Validate and create new url
func (ctrl UrlsController) CreateUrl(c echo.Context) error {
	var urlCreateModel url.UrlCreateViewModel
	if err := c.Bind(&urlCreateModel); err != nil {
		return v1.NewRequestError(ErrBindingRequestData, http.StatusBadRequest)
	}
	claims, err := web.GetUserClaims(c)
	if err != nil {
		return v1.NewRequestError(
			fmt.Errorf("invalid user claims: %w", err),
			http.StatusBadRequest,
		)
	}
	urlData, err := ctrl.Core.Create(c.Request().Context(), urlCreateModel, claims)
	if err != nil {
		return v1.NewRequestError(
			fmt.Errorf("failed to create url: %w", err),
			http.StatusInternalServerError,
		)
	}
	urlModel := url.NewUrlViewModel(urlData)
	return c.JSON(http.StatusCreated, urlModel)
}

// Delete url by ID
func (ctrl UrlsController) DeleteUrl(c echo.Context) error {
	urlId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return v1.NewRequestError(ErrInvalidId, http.StatusBadRequest)
	}
	if err := ctrl.Core.Delete(c.Request().Context(), urlId); err != nil {
		if errors.Is(err, url.ErrUrlNotFound) {
			return v1.NewRequestError(err, http.StatusNotFound)
		}
		return v1.NewRequestError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}
