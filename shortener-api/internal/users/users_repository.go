package users

import (
	"github.com/dawidhermann/shortener-api/internal/db"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type RepositoryUsers struct {
	connDb db.SqlConnection
}

func newRepositoryUsers(connDb db.SqlConnection) RepositoryUsers {
	return RepositoryUsers{
		connDb: connDb,
	}
}

func (repository RepositoryUsers) createUserEntity(username string, encPass string, email string) (int, error) {
	timestamp := time.Now()
	userId := 0
	err := repository.connDb.Db.QueryRow(
		"INSERT INTO users (username, password, email, created_at) VALUES($1, $2, $3, $4) RETURNING user_id",
		username, encPass, email, timestamp).Scan(&userId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return userId, nil
}

func (repository RepositoryUsers) getUserEntity(userId int) (user, error) {
	user := user{}
	err := repository.connDb.Db.QueryRow("SELECT user_id, username, password, email FROM users WHERE user_id=$1", userId).Scan(&user.UserId, &user.Username, &user.Password, &user.Email)
	return user, err
}

func (repository RepositoryUsers) getUserEntityByUsername(username string) (user, error) {
	user := user{}
	err := repository.connDb.Db.QueryRow("SELECT user_id, username, password, email FROM users WHERE username=$1", username).Scan(&user.UserId, &user.Username, &user.Password, &user.Email)
	return user, err
}

func (repository RepositoryUsers) updateUserEntity(user user) error {
	_, err := repository.connDb.Db.Exec("UPDATE users SET password=$1, email=$2 WHERE user_id=$3", user.Password, user.Email, user.UserId)
	return err
}

func (repository RepositoryUsers) deleteUserEntity(userId int) error {
	_, err := repository.connDb.Db.Exec("DELETE FROM users WHERE user_id=$1", userId)
	return err
}
