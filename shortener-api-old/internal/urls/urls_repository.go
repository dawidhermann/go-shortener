package urls

import (
	"database/sql"
	"github.com/dawidhermann/shortener-api/internal/db"
	"log"
	"strconv"
)

type RepositoryUrls struct {
	sqlConn db.SqlConnection
}

func newRepositoryUrls(connDb db.SqlConnection) RepositoryUrls {
	return RepositoryUrls{
		sqlConn: connDb,
	}
}

func (repository RepositoryUrls) createUrlEntity(shortenedUrl string, userId int) (int, error) {
	var userIdVal sql.NullString
	if userId != 0 {
		userIdVal = sql.NullString{String: strconv.Itoa(userId), Valid: true}
	} else {
		userIdVal = sql.NullString{}
	}
	urlId := 0
	err := repository.sqlConn.Db.QueryRow(
		"INSERT INTO urls (url_key, user_id) VALUES($1, $2) RETURNING url_id",
		shortenedUrl, userIdVal).Scan(&urlId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return urlId, err
}

func (repository RepositoryUrls) getUrlEntity(urlId int) (url, error) {
	var urlData url
	err := repository.sqlConn.Db.QueryRow("SELECT url_id, url_key FROM urls WHERE url_id=$1", urlId).Scan(&urlData.urlId, &urlData.urlKey)
	return urlData, err
}

func (repository RepositoryUrls) deleteUrlEntity(urlId int) error {
	_, err := repository.sqlConn.Db.Exec("DELETE FROM urls WHERE url_id=$1", urlId)
	return err
}
