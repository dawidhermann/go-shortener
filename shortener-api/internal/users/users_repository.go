package users

import (
	"github.com/dawidhermann/shortener-api/internal/db"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"time"
)

func createUserEntity(username string, encPass string, email string) (int, error) {
	timestamp := time.Now()
	userId := 0
	err := db.Db.QueryRow(
		"INSERT INTO users (username, password, email, created_at) VALUES($1, $2, $3, $4) RETURNING user_id",
		username, encPass, email, timestamp).Scan(&userId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id := strconv.Itoa(userId)
	getUserEntity(id)
	return userId, nil
}

func getUserEntity(userId string) (user, error) {
	user := user{}
	err := db.Db.QueryRow("SELECT user_id, username, password, email FROM users WHERE user_id=$1", userId).Scan(&user.UserId, &user.Username, &user.Password, &user.Email)
	log.Println(user)
	log.Println(err)
	return user, err
}

func updateUserEntity(user user) error {
	log.Println(user)
	_, err := db.Db.Exec("UPDATE users SET password=$1, email=$2 WHERE user_id=$3", user.Password, user.Email, user.UserId)
	return err
}

func deleteUserEntity(userId string) error {
	_, err := db.Db.Exec("DELETE FROM users WHERE user_id=$1", userId)
	return err
}
