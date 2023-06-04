package urls

import (
	"errors"
	"fmt"
	"github.com/dawidhermann/shortener-api/internal/db"
	"github.com/dawidhermann/shortener-api/internal/rpc"
	"github.com/dawidhermann/shortener-api/internal/users"
	"github.com/dawidhermann/shortener-api/internal/validators"
	"log"
)

var (
	ErrIncorrectUserId  = errors.New("incorrect user id")
	ErrUserDoesNotExist = errors.New("user does not exist")
	ErrInvalidUrl       = errors.New("url is not valid")
	ErrCreateUrlFailed  = errors.New("failed to create url")
	ErrUrlDoesNotExist  = errors.New("url does not exist")
	ErrDeleteUrlFailed  = errors.New("failed to delete url")
)

type ServiceUrls struct {
	connRpc        rpc.ConnRpc
	repositoryUrls RepositoryUrls
	serviceUsers   users.ServiceUsers
}

func NewServiceUrls(connRpc rpc.ConnRpc, connDb db.SqlConnection) ServiceUrls {
	return ServiceUrls{
		connRpc:        connRpc,
		repositoryUrls: newRepositoryUrls(connDb),
		serviceUsers:   users.NewServiceUsers(connDb),
	}
}

func (service ServiceUrls) CreateUrl(createModel UrlCreateModel, userId int) (int, error) {
	if userId != 0 {
		userData, err := service.serviceUsers.GetUser(userId)
		if err != nil {
			return 0, err
		}
		if userData.UserId == 0 {
			return 0, ErrUserDoesNotExist
		}
	}
	fmt.Println(userId)
	err := validators.ValidateUrl(createModel.TargetUrl)
	if err != nil {
		return 0, ErrInvalidUrl
	}
	shortenedUrl, err := service.connRpc.CreateShortenUrl(createModel.TargetUrl)
	if err != nil {
		return 0, ErrCreateUrlFailed
	}
	urlId, err := service.repositoryUrls.createUrlEntity(shortenedUrl, userId)
	if err != nil {
		log.Print(err)
		return 0, ErrCreateUrlFailed
	}
	return urlId, err
}

func (service ServiceUrls) GetUrl(urlId int) (url, error) {
	urlData, err := service.repositoryUrls.getUrlEntity(urlId)
	if err != nil {
		log.Println(err)
		return url{}, ErrUrlDoesNotExist
	}
	return urlData, err
}

func (service ServiceUrls) DeleteUrl(urlId int) error {
	urlData, err := service.GetUrl(urlId)
	if err != nil {
		return ErrUrlDoesNotExist
	}
	err = service.connRpc.DeleteShortenedUrl(urlData.urlKey)
	if err != nil {
		return ErrDeleteUrlFailed
	}
	err = service.repositoryUrls.deleteUrlEntity(urlId)
	if err != nil {
		return ErrUrlDoesNotExist
	}
	return nil
}
